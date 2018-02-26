package tests

import (
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

func list(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		t.Run("When two functions are deployed list should return two functions", when2FunctionsDeployedListReturns2Functions)
	})

	cleanupDeployedFunctions()
}

func when2FunctionsDeployedListReturns2Functions(t *testing.T) {
	envVars := map[string]string{}
	envVars["custom_env"] = "custom_env_value"

	deploy1 := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "list1",
		Network:    "func_functions",
		EnvProcess: "sha512sum",
		EnvVars:    envVars,
	}
	_, _ = deployFunction(deploy1)

	deploy2 := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "list2",
		Network:    "func_functions",
		EnvProcess: "sha512sum",
		EnvVars:    envVars,
	}

	_, _ = deployFunction(deploy2)

	assertList(t, 2)
}

func assertList(t *testing.T, count int) {
	fs := listFunctions(t)

	if len(fs) != count {
		t.Logf("List functions got: %d, want: %d", len(fs), count)
		t.Fail()
	}
}
