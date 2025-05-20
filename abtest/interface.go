package abtest

import (
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/api"
	cfgexp "github.com/aliyun/aliyun-pairec-config-go-sdk/v2/experiments"
	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/model"
)

// ExperimentClient defines the subset of methods used in pairec
// so different implementations can be plugged in.
type ExperimentClient interface {
	MatchExperiment(sceneName string, ctx *model.ExperimentContext) *model.ExperimentResult
	GetSceneParams(sceneName string) model.SceneParams
	BackflowFeatureConsistencyCheckJobData(data *model.FeatureConsistencyBackflowData) (api.FeatureConsistencyBackflowResponse, error)
	SyncFeatureConsistencyCheckJobReplayLog(data *model.FeatureConsistencyReplyData) (api.FeatureConsistencyReplyResponse, error)
	GetTrafficControlTaskMetaData(env string, currentTimestamp int64) []model.TrafficControlTask
	GetTrafficControlTargetData(env, sceneName string, currentTimestamp int64) map[string]model.TrafficControlTarget
	GetTrafficControlTargetTraffic(env, sceneName string, idList ...string) []cfgexp.TrafficControlTargetTraffic
}
