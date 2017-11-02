package tests

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alexellis/faas/gateway/requests"
)

func Test_Invoke(t *testing.T) {
	t.Run("group", func(t *testing.T) {
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
	attempts := 30 // i.e. 30x2s = 1m
	delay := time.Millisecond * 2000
	uri := os.Getenv("gateway_url") + "function/" + name
	success := false

	for i := 0; i < attempts; i++ {
		bytesOut, res, err := httpReq(uri, "POST", nil)

		if err != nil {
			t.Log(err.Error())
			continue
		}
		if res.StatusCode != http.StatusOK {
			t.Logf("[%d/%d] Bad response want: %d, got: %d", i+1, attempts, http.StatusOK, res.StatusCode)
			t.Logf(uri)
			if i == attempts-1 {
				t.Logf("Failing after: %d attempts", attempts)
			}
			time.Sleep(delay)
			continue
		} else {
			t.Logf("[%d/%d] Correct response: %d", i+1, attempts, res.StatusCode)
		}

		out := string(bytesOut)
		if strings.Contains(out, expected) == false {
			t.Logf("want: %s, got: %s", expected, out)
		} else {
			success = true
		}

		break
	}

	if !success {
		t.Fail()
	}
}
