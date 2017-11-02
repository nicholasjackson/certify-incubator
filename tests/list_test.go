package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/alexellis/faas/gateway/requests"
)

func Test_List(t *testing.T) {
	t.Run("group", func(t *testing.T) {
		t.Run("Test_When2FunctionsDeployedListReturns2Functions", when2FunctionsDeployedListReturns2Functions)
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
	bytesOut, res, err := httpReq(os.Getenv("gateway_url")+"system/functions", "GET", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("got %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
		return
	}

	functions := []requests.Function{}
	err = json.Unmarshal(bytesOut, &functions)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if len(functions) != count {
		t.Logf("List functions got: %s, want: %s", len(functions), count)
		t.Fail()
	}
}
