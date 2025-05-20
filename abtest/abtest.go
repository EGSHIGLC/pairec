package abtest

import (
	"os"

	"github.com/aliyun/aliyun-pairec-config-go-sdk/v2/model"
)

var experimentClient ExperimentClient

// LoadFromEnvironment create abtest instance use env, env list:
//
// ENV params list:
//
//		PAIREC_ENVIRONMENT is the environment type, valid values are: daily, prepub, product
//		REGION region of pairec console instance, like cn-beijing,cn-hangzhou
//		INSTANCE_ID id of pairec console instance
//	    AccessKey  aliyun accessKeyId
//	    AccessSecret  aliyun accessKeySecret
func LoadFromEnvironment() {
	abType := os.Getenv("ABTEST_TYPE")
	if abType == "" {
		abType = "pairec"
	}

	var client ExperimentClient
	var err error
	switch abType {
	case "token":
		client, err = newTokenClientFromEnv()
	default:
		client, err = newPairecClientFromEnv()
	}
	if err != nil {
		panic(err)
	}

	experimentClient = client
}
func GetExperimentClient() ExperimentClient {
	return experimentClient
}

func GetParams(sceneName string) model.SceneParams {
	if experimentClient == nil {
		return model.NewSceneParams()
	}
	return experimentClient.GetSceneParams(sceneName)
}
