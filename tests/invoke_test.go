package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

func invoke(t *testing.T) {
	t.Run("invoke", func(t *testing.T) {
		t.Run("Invoke a function and check response", basicInvoke)
	})
	cleanupDeployedFunctions()
}

func basicInvoke(t *testing.T) {
	envVars := map[string]string{}
	envVars["custom_env"] = "custom_env_value"
	deploy := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "invoke-test",
		Network:    "func_functions",
		EnvProcess: "env",
		EnvVars:    envVars,
	}

	_, _ = deployFunction(deploy)

	assertInvoke(t, "invoke-test", "custom_env_value")
}

func assertInvoke(t *testing.T, name string, expected string) {
	body, _, err := httpReqWithRetry(
		"function/"+name,
		"POST",
		[]byte{},
		100,
		1000,
		http.StatusOK,
	)

	if err != nil {
		t.Fatal(err)
		return
	}

	out := string(body)
	if strings.Contains(out, expected) == false {
		t.Fatalf("want: %s, got: %s", expected, out)
	}
}
