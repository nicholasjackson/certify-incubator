package tests

import (
	"net/http"
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

/*
func Test_Pipeline(t *testing.T) {
	envVars := map[string]string{}
	deploy := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "stronghash",
		Network:    "func_functions",
		EnvProcess: "sha512sum",
		EnvVars:    envVars,
	}

	DeployTest(t, deploy)

	TestList(t)
}
*/

func deploy(t *testing.T) {
	t.Run("deploy", func(t *testing.T) {
		t.Run("Deploy function passing custom environment variables", passingCustomEnvVars)
	})

	cleanupDeployedFunctions()
}

func passingCustomEnvVars(t *testing.T) {
	envVars := map[string]string{}
	envVars["custom_env"] = "custom_env_value"

	deploy := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "env-test",
		Network:    "func_functions",
		EnvProcess: "env",
		EnvVars:    envVars,
	}

	assertDeploy(t, deploy)
}

func assertDeploy(t *testing.T, createRequest requests.CreateFunctionRequest) {
	res, err := deployFunction(createRequest)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
	}
}
