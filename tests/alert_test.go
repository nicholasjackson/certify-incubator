package tests

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/alexellis/faas/gateway/requests"
)

func alerts(t *testing.T) {
	t.Run("alerts", func(t *testing.T) {
		t.Run("Scale alert scales function", scaleUpAlertScalesFunction)
	})

	cleanupDeployedFunctions()
}

func scaleUpAlertScalesFunction(t *testing.T) {
	envVars := map[string]string{}
	envVars["custom_env"] = "custom_env_value"

	scale := requests.CreateFunctionRequest{
		Image:      "functions/alpine:latest",
		Service:    "scaletest",
		Network:    "func_functions",
		EnvProcess: "sha512sum",
		EnvVars:    envVars,
	}
	_, _ = deployFunction(scale)

	sendScaleRequest(t)

	assertScale(t)
}

func sendScaleRequest(t *testing.T) {
	_, r, err := httpReq(
		os.Getenv("gateway_url")+"system/alert",
		"POST",
		bytes.NewBuffer([]byte(scalePayload)),
	)

	if err != nil && r.StatusCode != http.StatusOK {
		t.Fatal(err)
	}
}

func assertScale(t *testing.T) {
	fs := listFunctions(t)

	if fs == nil {
		t.Fatal("no functions running")
	}

	if fs[0].Replicas != 5 {
		t.Logf("Expected function to have 2 instances got %d", fs[0].Replicas)
		t.Fail()
	}
}

var scalePayload = `
{"receiver": "scale-up",
  "status": "firing",
  "alerts": [{
      "status": "firing",
      "labels": {
          "alertname": "APIHighInvocationRate",
          "code": "200",
          "function_name": "scaletest",
          "instance": "gateway:8080",
          "job": "gateway",
          "monitor": "faas-monitor",
          "service": "gateway",
          "severity": "major",
          "value": "8.998200359928017"
      },
      "annotations": {
          "description": "High invocation total on gateway:8080",
          "summary": "High invocation total on gateway:8080"
      },
      "startsAt": "2017-03-15T15:52:57.805Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://4156cb797423:9090/graph?g0.expr=rate%28gateway_function_invocation_total%5B10s%5D%29+%3E+5\u0026g0.tab=0"
  }],
  "groupLabels": {
      "alertname": "APIHighInvocationRate",
      "service": "gateway"
  },
  "commonLabels": {
      "alertname": "APIHighInvocationRate",
      "code": "200",
      "function_name": "scaletest",
      "instance": "gateway:8080",
      "job": "gateway",
      "monitor": "faas-monitor",
      "service": "gateway",
      "severity": "major",
      "value": "8.998200359928017"
  },
  "commonAnnotations": {
      "description": "High invocation total on gateway:8080",
      "summary": "High invocation total on gateway:8080"
  },
  "externalURL": "http://f054879d97db:9093",
  "version": "3",
  "groupKey": 18195285354214864953
}
`
