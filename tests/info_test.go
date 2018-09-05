package tests

import (
	"net/http"
	"strings"
	"testing"
)

func info(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		t.Run("Invoke the Info endpoint and check response", assertInfo)
	})
}

func assertInfo(t *testing.T) {
	body, _, err := httpReqWithRetry(
		"system/info",
		"GET",
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
	if strings.Contains(out, "\"orchestration\":\"nomad\"") == false {
		t.Fatalf("want: %s, got: %s", "nomad", out)
	}
}
