package remotecfg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"

	"connectrpc.com/connect"
	collectorv1 "github.com/grafana/alloy-remote-config/api/gen/proto/go/collector/v1"
	"github.com/grafana/alloy-remote-config/api/gen/proto/go/collector/v1/collectorv1connect"
	"github.com/grafana/alloy/internal/component"
	_ "github.com/grafana/alloy/internal/component/loki/process"
	"github.com/grafana/alloy/internal/featuregate"
	alloy_runtime "github.com/grafana/alloy/internal/runtime"
	"github.com/grafana/alloy/internal/runtime/logging"
	"github.com/grafana/alloy/internal/service"
	"github.com/grafana/alloy/internal/service/livedebugging"
	"github.com/grafana/alloy/internal/util"
	"github.com/grafana/alloy/syntax"
	"github.com/grafana/alloy/syntax/ast"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnDiskCache(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	client := &collectorClient{}

	var registerCalled atomic.Bool
	client.registerCollectorFunc = buildRegisterCollectorFunc(&registerCalled)
	url := "https://example.com/"

	// The contents of the on-disk cache.
	cacheContents := `loki.process "default" { forward_to = [] }`

	// Create a new service.
	env := newTestEnvironment(t, client)
	require.NoError(t, env.ApplyConfig(fmt.Sprintf(`
		url = "%s"
	`, url)))

	// Mock client to return an unparseable response.
	client.getConfigFunc = buildGetConfigHandler("unparseable config", "", false)

	// Write the cache contents, and run the service.
	err := os.WriteFile(env.svc.dataPath, []byte(cacheContents), 0644)
	require.NoError(t, err)

	// Run the service.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, env.Run(ctx))
	}()

	// As the API response was unparseable, verify that the service has loaded
	// the on-disk cache contents.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		b, err := env.svc.getCachedConfig()
		assert.NoError(c, err)
		assert.Equal(c, cacheContents, string(b))
	}, time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestGoodBadGood(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	url := "https://example.com/"
	cfgGood := `loki.process "default" { forward_to = [] }`
	cfgBad := `unparseable config`

	client := &collectorClient{}

	// Mock client to return a valid response.
	var registerCalled atomic.Bool
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfgGood, "", false)
	client.registerCollectorFunc = buildRegisterCollectorFunc(&registerCalled)
	client.mut.Unlock()

	// Create a new service.
	env := newTestEnvironment(t, client)
	require.NoError(t, env.ApplyConfig(fmt.Sprintf(`
		url            = "%s"
		poll_frequency = "10s"
	`, url)))

	// Run the service.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, env.Run(ctx))
	}()

	require.Eventually(t, func() bool { return registerCalled.Load() }, 1*time.Second, 10*time.Millisecond)

	// As the API response was successful, verify that the service has loaded
	// the valid response.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfgGood)), env.svc.getLastLoadedCfgHash())
	}, time.Second, 10*time.Millisecond)

	// Update the response returned by the API to an invalid configuration.
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfgBad, "", false)
	client.mut.Unlock()

	// Verify that the service has still the same "good" configuration has
	// loaded and flushed on disk, but also recorded the "bad" hash saved for
	// comparison.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		b, err := env.svc.getCachedConfig()
		assert.NoError(c, err)
		assert.Equal(c, cfgGood, string(b))
	}, 1*time.Second, 10*time.Millisecond)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfgBad)), env.svc.getLastLoadedCfgHash())
	}, 1*time.Second, 10*time.Millisecond)

	// Update the response returned by the API to the previous "good"
	// configuration.
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfgGood, "", false)
	client.mut.Unlock()

	// Verify that the service has updated the hash.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfgGood)), env.svc.getLastLoadedCfgHash())
	}, 1*time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestAPIResponse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	url := "https://example.com/"
	cfg1 := `loki.process "default" { forward_to = [] }`
	cfg2 := `loki.process "updated" { forward_to = [] }`

	client := &collectorClient{}

	// Mock client to return a valid response.
	var registerCalled atomic.Bool
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfg1, "", false)
	client.registerCollectorFunc = buildRegisterCollectorFunc(&registerCalled)
	client.mut.Unlock()

	// Create a new service.
	env := newTestEnvironment(t, client)
	require.NoError(t, env.ApplyConfig(fmt.Sprintf(`
		url            = "%s"
		poll_frequency = "10s"
	`, url)))

	// Run the service.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, env.Run(ctx))
	}()

	require.Eventually(t, func() bool { return registerCalled.Load() }, 1*time.Second, 10*time.Millisecond)

	// As the API response was successful, verify that the service has loaded
	// the valid response.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfg1)), env.svc.getLastLoadedCfgHash())
	}, time.Second, 10*time.Millisecond)

	// Update the response returned by the API.
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfg2, "", false)
	client.mut.Unlock()

	// Verify that the service has loaded the updated response.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfg2)), env.svc.getLastLoadedCfgHash())
	}, 1*time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestAPIResponseNotModified(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	url := "https://example.com/"
	cfg1 := `loki.process "default" { forward_to = [] }`

	client := &collectorClient{}

	// Mock client to return a valid response.
	var registerCalled atomic.Bool
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler(cfg1, "12345", false)
	client.registerCollectorFunc = buildRegisterCollectorFunc(&registerCalled)
	client.mut.Unlock()

	// Create a new service.
	env := newTestEnvironment(t, client)
	require.NoError(t, env.ApplyConfig(fmt.Sprintf(`
		url            = "%s"
		poll_frequency = "10s"
	`, url)))

	// Run the service.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, env.Run(ctx))
	}()

	require.Eventually(t, func() bool { return registerCalled.Load() }, 1*time.Second, 10*time.Millisecond)

	// As the API response was successful, verify that the service has loaded
	// the valid response.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Equal(c, getHash([]byte(cfg1)), env.svc.getLastLoadedCfgHash())
	}, time.Second, 10*time.Millisecond)

	// Update the response returned by the API.
	client.mut.Lock()
	client.getConfigFunc = buildGetConfigHandler("", "12345", true)
	client.mut.Unlock()

	calls := client.getConfigCalls.Load()

	// Verify that the service has loaded the updated response.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		// Ensure that getConfig has been called again since changing the response.
		assert.Greater(c, client.getConfigCalls.Load(), calls)
		assert.Equal(c, getHash([]byte(cfg1)), env.svc.getLastLoadedCfgHash())
	}, 1*time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestUserAgentHeader(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	cfg := `loki.process "default" { forward_to = [] }`

	// Track captured User-Agent headers
	var capturedUserAgent atomic.Value
	var registerCalled atomic.Bool

	// Create a test server that captures the User-Agent header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the User-Agent header
		capturedUserAgent.Store(r.Header.Get("User-Agent"))

		// Mock a successful register collector response
		if r.URL.Path == "/collector.v1.CollectorService/RegisterCollector" {
			registerCalled.Store(true)
			w.Header().Set("Content-Type", "application/proto")
			w.WriteHeader(http.StatusOK)
			// Create empty protobuf response for RegisterCollectorResponse
			response := &collectorv1.RegisterCollectorResponse{}
			data, err := proto.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		}

		// Mock a successful get config response
		if r.URL.Path == "/collector.v1.CollectorService/GetConfig" {
			w.Header().Set("Content-Type", "application/proto")
			w.WriteHeader(http.StatusOK)
			// Create a minimal protobuf response for GetConfigResponse
			// This is a simple hardcoded protobuf message with the content field
			response := &collectorv1.GetConfigResponse{
				Content: cfg,
			}
			data, err := proto.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create a new service with default factory (uses real interceptor)
	svc, err := New(Options{
		Logger:      util.TestLogger(t),
		StoragePath: t.TempDir(),
	})
	require.NoError(t, err)

	env := &testEnvironment{
		t:   t,
		svc: svc,
	}

	// Configure the service to use our test server
	require.NoError(t, env.ApplyConfig(fmt.Sprintf(`
		url            = "%s"
		poll_frequency = "10s"
	`, server.URL)))

	// Run the service
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, env.Run(ctx))
	}()

	// Wait for the register call to complete
	require.Eventually(t, func() bool {
		return registerCalled.Load()
	}, 2*time.Second, 10*time.Millisecond)

	// Verify that the User-Agent header was captured and contains "Alloy"
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		userAgent := capturedUserAgent.Load()
		assert.NotNil(c, userAgent, "User-Agent header should be captured")
		if userAgent != nil {
			userAgentStr := userAgent.(string)
			assert.NotEmpty(c, userAgentStr, "User-Agent header should be present")
			assert.Contains(c, userAgentStr, "Alloy", "User-Agent should contain 'Alloy'")
		}
	}, 1*time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func buildGetConfigHandler(in string, hash string, notModified bool) func(context.Context, *connect.Request[collectorv1.GetConfigRequest]) (*connect.Response[collectorv1.GetConfigResponse], error) {
	return func(context.Context, *connect.Request[collectorv1.GetConfigRequest]) (*connect.Response[collectorv1.GetConfigResponse], error) {
		rsp := &connect.Response[collectorv1.GetConfigResponse]{
			Msg: &collectorv1.GetConfigResponse{
				Content:     in,
				NotModified: notModified,
				Hash:        hash,
			},
		}
		return rsp, nil
	}
}

func buildRegisterCollectorFunc(called *atomic.Bool) func(ctx context.Context, req *connect.Request[collectorv1.RegisterCollectorRequest]) (*connect.Response[collectorv1.RegisterCollectorResponse], error) {
	return func(ctx context.Context, req *connect.Request[collectorv1.RegisterCollectorRequest]) (*connect.Response[collectorv1.RegisterCollectorResponse], error) {
		called.Store(true)
		return &connect.Response[collectorv1.RegisterCollectorResponse]{
			Msg: &collectorv1.RegisterCollectorResponse{},
		}, nil
	}
}

type testEnvironment struct {
	t   *testing.T
	svc *Service
}

func newTestEnvironment(t *testing.T, client *collectorClient) *testEnvironment {
	svc, err := New(Options{
		Logger:      util.TestLogger(t),
		StoragePath: t.TempDir(),
	})
	svc.clientFactory = func(_ Arguments) (collectorv1connect.CollectorServiceClient, error) {
		return client, nil
	}
	require.NoError(t, err)

	return &testEnvironment{
		t:   t,
		svc: svc,
	}
}

func (env *testEnvironment) ApplyConfig(config string) error {
	var args Arguments
	if err := syntax.Unmarshal([]byte(config), &args); err != nil {
		return err
	}
	// The lower limit of the poll_frequency argument would slow our tests
	// considerably; let's artificially lower it after the initial validation
	// has taken place.
	args.PollFrequency /= 100
	return env.svc.Update(args)
}

func (env *testEnvironment) Run(ctx context.Context) error {
	return env.svc.Run(ctx, fakeHost{})
}

type fakeHost struct{}

var _ service.Host = (fakeHost{})

func (fakeHost) GetComponent(id component.ID, opts component.InfoOptions) (*component.Info, error) {
	return nil, fmt.Errorf("no such component %s", id)
}

func (fakeHost) ListComponents(moduleID string, opts component.InfoOptions) ([]*component.Info, error) {
	if moduleID == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("no such module %q", moduleID)
}

func (fakeHost) GetServiceConsumers(_ string) []service.Consumer { return nil }
func (fakeHost) GetService(_ string) (service.Service, bool)     { return nil, false }

func (f fakeHost) NewController(id string) service.Controller {
	logger, _ := logging.New(io.Discard, logging.DefaultOptions)
	ctrl := alloy_runtime.New(alloy_runtime.Options{
		ControllerID:    ServiceName,
		Logger:          logger,
		Tracer:          nil,
		DataPath:        "",
		MinStability:    featuregate.StabilityGenerallyAvailable,
		Reg:             prometheus.NewRegistry(),
		OnExportsChange: func(map[string]interface{}) {},
		Services:        []service.Service{livedebugging.New()},
	})

	return serviceController{ctrl}
}

type collectorClient struct {
	getConfigCalls        atomic.Int32
	mut                   sync.RWMutex
	getConfigFunc         func(context.Context, *connect.Request[collectorv1.GetConfigRequest]) (*connect.Response[collectorv1.GetConfigResponse], error)
	registerCollectorFunc func(ctx context.Context, req *connect.Request[collectorv1.RegisterCollectorRequest]) (*connect.Response[collectorv1.RegisterCollectorResponse], error)
}

func (c *collectorClient) GetConfig(ctx context.Context, req *connect.Request[collectorv1.GetConfigRequest]) (*connect.Response[collectorv1.GetConfigResponse], error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if c.getConfigFunc != nil {
		c.getConfigCalls.Inc()
		return c.getConfigFunc(ctx, req)
	}

	panic("getConfigFunc not set")
}

func (c *collectorClient) RegisterCollector(ctx context.Context, req *connect.Request[collectorv1.RegisterCollectorRequest]) (*connect.Response[collectorv1.RegisterCollectorResponse], error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if c.registerCollectorFunc != nil {
		return c.registerCollectorFunc(ctx, req)
	}

	panic("registerCollectorFunc not set")
}

func (c *collectorClient) UnregisterCollector(ctx context.Context, req *connect.Request[collectorv1.UnregisterCollectorRequest]) (*connect.Response[collectorv1.UnregisterCollectorResponse], error) {
	panic("unregisterCollector isn't wired yet")
}

type serviceController struct {
	f *alloy_runtime.Runtime
}

func (sc serviceController) Run(ctx context.Context) { sc.f.Run(ctx) }
func (sc serviceController) LoadSource(b []byte, args map[string]any, configPath string) (*ast.File, error) {
	source, err := alloy_runtime.ParseSource("", b)
	if err != nil {
		return nil, err
	}
	return source.SourceFiles()[""], sc.f.LoadSource(source, args, configPath)
}
func (sc serviceController) Ready() bool { return sc.f.Ready() }
