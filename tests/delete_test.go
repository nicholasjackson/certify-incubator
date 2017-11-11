package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/alexellis/faas/gateway/requests"
)

func basicDelete(t *testing.T) {
	envVars := map[string]string{}
	envVars["custom_env"] = "custom_env_value"
	deploy := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "delete-test",
		Network:    "func_functions",
		EnvProcess: "env",
		EnvVars:    envVars,
	}
	_, _, err := httpReq(
		os.Getenv("gateway_url")+"system/functions",
		"POST",
		makeReader(deploy),
	)

	if err != nil {
		t.Fatal(err)
	}
	_, res, err := httpReq(
		os.Getenv("gateway_url")+"system/functions",
		"DELETE",
		makeReader(requests.DeleteFunctionRequest{
			FunctionName: "delete-test",
		}),
	)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("Expected status code 200, got", res.StatusCode)
	}
}
