package abtest

import (
	"errors"
	"os"

	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/api"
	cfgexp "github.com/aliyun/aliyun-pairec-config-go-sdk/v2/experiments"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/model"
)

// tokenClient implements ExperimentClient using the legacy host/token SDK.
// tokenClient is a minimal stub for the legacy host/token sdk. It currently
// does not implement real functionality but keeps the old interface compatible.
type tokenClient struct{}

func newTokenClientFromEnv() (ExperimentClient, error) {
	region := os.Getenv("REGION")
	if region == "" {
		return nil, errors.New("env REGION empty")
	}
	// legacy sdk not available; return a stub client
	return &tokenClient{}, nil
}

func (t *tokenClient) MatchExperiment(sceneName string, ctx *model.ExperimentContext) *model.ExperimentResult {
	// not implemented, return empty result
	return model.NewExperimentResult(sceneName, ctx)
}

func (t *tokenClient) GetSceneParams(sceneName string) model.SceneParams {
	// not supported in old sdk
	return model.NewSceneParams()
}

func (t *tokenClient) BackflowFeatureConsistencyCheckJobData(data *model.FeatureConsistencyBackflowData) (api.FeatureConsistencyBackflowResponse, error) {
	// not supported
	return api.FeatureConsistencyBackflowResponse{}, nil
}

func (t *tokenClient) SyncFeatureConsistencyCheckJobReplayLog(data *model.FeatureConsistencyReplyData) (api.FeatureConsistencyReplyResponse, error) {
	// not supported
	return api.FeatureConsistencyReplyResponse{}, nil
}

func (t *tokenClient) GetTrafficControlTaskMetaData(env string, currentTimestamp int64) []model.TrafficControlTask {
	return nil
}

func (t *tokenClient) GetTrafficControlTargetData(env, sceneName string, currentTimestamp int64) map[string]model.TrafficControlTarget {
	return nil
}

func (t *tokenClient) GetTrafficControlTargetTraffic(env, sceneName string, idList ...string) []cfgexp.TrafficControlTargetTraffic {
	return nil
}
