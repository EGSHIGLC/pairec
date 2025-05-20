package abtest

import (
	"errors"
	"os"

	"github.com/alibaba/pairec/v2/log"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/api"
	cfgexp "github.com/aliyun/aliyun-pairec-config-go-sdk/v2/experiments"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/model"
)

// pairecClient implements ExperimentClient using the pairec config SDK.
type pairecClient struct{ *cfgexp.ExperimentClient }

func newPairecClientFromEnv() (ExperimentClient, error) {
	env := os.Getenv("PAIREC_ENVIRONMENT")
	if env == "" {
		return nil, errors.New("env PAIREC_ENVIRONMENT empty")
	}
	region := os.Getenv("REGION")
	instanceID := os.Getenv("INSTANCE_ID")
	accessID := os.Getenv("AccessKey")
	accessSecret := os.Getenv("AccessSecret")
	if region == "" {
		return nil, errors.New("env REGION empty")
	}
	if instanceID == "" {
		return nil, errors.New("env INSTANCE_ID empty")
	}

	l := log.ABTestLogger{}
	opts := []cfgexp.ClientOption{cfgexp.WithLogger(cfgexp.LoggerFunc(l.Infof)), cfgexp.WithErrorLogger(cfgexp.LoggerFunc(l.Errorf))}
	if ep := os.Getenv("PAIREC_CONFIG_ENDPOINT"); ep != "" {
		opts = append(opts, cfgexp.WithDomain(ep))
	}
	client, err := cfgexp.NewExperimentClient(instanceID, region, accessID, accessSecret, env, opts...)
	if err != nil {
		return nil, err
	}
	return &pairecClient{client}, nil
}

// Wrapper methods
func (p *pairecClient) MatchExperiment(sceneName string, ctx *model.ExperimentContext) *model.ExperimentResult {
	return p.ExperimentClient.MatchExperiment(sceneName, ctx)
}

func (p *pairecClient) GetSceneParams(sceneName string) model.SceneParams {
	return p.ExperimentClient.GetSceneParams(sceneName)
}

func (p *pairecClient) BackflowFeatureConsistencyCheckJobData(data *model.FeatureConsistencyBackflowData) (api.FeatureConsistencyBackflowResponse, error) {
	return p.ExperimentClient.BackflowFeatureConsistencyCheckJobData(data)
}

func (p *pairecClient) SyncFeatureConsistencyCheckJobReplayLog(data *model.FeatureConsistencyReplyData) (api.FeatureConsistencyReplyResponse, error) {
	return p.ExperimentClient.SyncFeatureConsistencyCheckJobReplayLog(data)
}

func (p *pairecClient) GetTrafficControlTaskMetaData(env string, currentTimestamp int64) []model.TrafficControlTask {
	return p.ExperimentClient.GetTrafficControlTaskMetaData(env, currentTimestamp)
}

func (p *pairecClient) GetTrafficControlTargetData(env, sceneName string, currentTimestamp int64) map[string]model.TrafficControlTarget {
	return p.ExperimentClient.GetTrafficControlTargetData(env, sceneName, currentTimestamp)
}

func (p *pairecClient) GetTrafficControlTargetTraffic(env, sceneName string, idList ...string) []cfgexp.TrafficControlTargetTraffic {
	return p.ExperimentClient.GetTrafficControlTargetTraffic(env, sceneName, idList...)
}
