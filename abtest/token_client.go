package abtest

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alibaba/pairec/v2/log"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pairecservice"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/api"
	cfgexp "github.com/aliyun/aliyun-pairec-config-go-sdk/v2/experiments"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/model"
)

// tokenClient implements ExperimentClient using a bearer token to authenticate
// with the pairec service API.
type tokenClient struct{ *cfgexp.ExperimentClient }

// defaultTransport replicates the settings from the SDK's internal client so
// tokenClient behaves the same as pairecClient.
var defaultTransport = &http.Transport{
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		return d.DialContext(ctx, "tcp4", addr)
	},
	MaxIdleConns:          200,
	MaxIdleConnsPerHost:   200,
	MaxConnsPerHost:       200,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func newAPIClientWithToken(instanceID, region, token string) (*api.APIClient, error) {
	// Create a regular API client to initialize internal fields
	c, err := api.NewAPIClient(instanceID, region, "", "")
	if err != nil {
		return nil, err
	}

	sdkCfg := sdk.NewConfig()
	sdkCfg.Scheme = "https"
	cred := credentials.NewBearerTokenCredential(token)

	client, err := pairecservice.NewClientWithOptions(region, sdkCfg, cred)
	if err != nil {
		return nil, err
	}
	client.SetTransport(defaultTransport)

	// replace the underlying client with the bearer token version
	c.Client = client
	return c, nil
}

func newExperimentClientWithToken(instanceID, region, token, environment string, opts ...cfgexp.ClientOption) (*cfgexp.ExperimentClient, error) {
	client, err := cfgexp.NewExperimentClient(instanceID, region, "", "", environment, opts...)
	if err != nil {
		return nil, err
	}

	apiClient, err := newAPIClientWithToken(instanceID, region, token)
	if err != nil {
		return nil, err
	}
	client.APIClient = apiClient
	return client, nil
}

func newTokenClientFromEnv() (ExperimentClient, error) {
	env := os.Getenv("PAIREC_ENVIRONMENT")
	if env == "" {
		return nil, errors.New("env PAIREC_ENVIRONMENT empty")
	}
	region := os.Getenv("REGION")
	instanceID := os.Getenv("INSTANCE_ID")
	token := os.Getenv("ABTEST_TOKEN")
	if region == "" {
		return nil, errors.New("env REGION empty")
	}
	if instanceID == "" {
		return nil, errors.New("env INSTANCE_ID empty")
	}
	if token == "" {
		return nil, errors.New("env ABTEST_TOKEN empty")
	}

	l := log.ABTestLogger{}
	opts := []cfgexp.ClientOption{cfgexp.WithLogger(cfgexp.LoggerFunc(l.Infof)), cfgexp.WithErrorLogger(cfgexp.LoggerFunc(l.Errorf))}
	if ep := os.Getenv("PAIREC_CONFIG_ENDPOINT"); ep != "" {
		opts = append(opts, cfgexp.WithDomain(ep))
	}

	client, err := newExperimentClientWithToken(instanceID, region, token, env, opts...)
	if err != nil {
		return nil, err
	}

	return &tokenClient{client}, nil
}

// Wrapper methods
func (t *tokenClient) MatchExperiment(sceneName string, ctx *model.ExperimentContext) *model.ExperimentResult {
	return t.ExperimentClient.MatchExperiment(sceneName, ctx)
}

func (t *tokenClient) GetSceneParams(sceneName string) model.SceneParams {
	return t.ExperimentClient.GetSceneParams(sceneName)
}

func (t *tokenClient) BackflowFeatureConsistencyCheckJobData(data *model.FeatureConsistencyBackflowData) (api.FeatureConsistencyBackflowResponse, error) {
	return t.ExperimentClient.BackflowFeatureConsistencyCheckJobData(data)
}

func (t *tokenClient) SyncFeatureConsistencyCheckJobReplayLog(data *model.FeatureConsistencyReplyData) (api.FeatureConsistencyReplyResponse, error) {
	return t.ExperimentClient.SyncFeatureConsistencyCheckJobReplayLog(data)
}

func (t *tokenClient) GetTrafficControlTaskMetaData(env string, currentTimestamp int64) []model.TrafficControlTask {
	return t.ExperimentClient.GetTrafficControlTaskMetaData(env, currentTimestamp)
}

func (t *tokenClient) GetTrafficControlTargetData(env, sceneName string, currentTimestamp int64) map[string]model.TrafficControlTarget {
	return t.ExperimentClient.GetTrafficControlTargetData(env, sceneName, currentTimestamp)
}

func (t *tokenClient) GetTrafficControlTargetTraffic(env, sceneName string, idList ...string) []cfgexp.TrafficControlTargetTraffic {
	return t.ExperimentClient.GetTrafficControlTargetTraffic(env, sceneName, idList...)
}
