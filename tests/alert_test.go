package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/openfaas/faas/gateway/requests"
	"github.com/pkg/errors"
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
	_, _, err := httpReqWithRetry(
		"system/alert",
		"POST",
		[]byte(scalePayload),
		4,
		100,
		http.StatusOK,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func assertScale(t *testing.T) {
	// retry until the replicas have been created, this can take variable time based on the machine performance
	r := retrier.New(
		retrier.ConstantBackoff(30, 1000*time.Millisecond),
		nil,
	)

	errs := errors.New("")

	err := r.Run(func() error {

		fs := listFunctionDetail(t, "scaletest")

		if fs.Name == "" {
			errs = errors.Wrap(errs, "no functions running\n")
			return errs
		}

		if fs.Replicas != 4 {
			errs = errors.Wrap(errs, fmt.Sprintf("Expected function to have 4 instances got %d\n", fs.Replicas))
			return errs
		}

		if fs.AvailableReplicas != 4 {
			errs = errors.Wrap(errs, fmt.Sprintf("Expected function to have 4 available replicas got %d\n", fs.AvailableReplicas))
			return errs
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
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
