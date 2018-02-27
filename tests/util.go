package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/openfaas/faas/gateway/requests"
	"github.com/pkg/errors"
)

var deployedFunctions []string

func init() {
	deployedFunctions = make([]string, 0)
}

func makeReader(input interface{}) *bytes.Buffer {
	res, _ := json.Marshal(input)
	return bytes.NewBuffer(res)
}

// httpReqWithRetry makes a http request n number of times until success code is matched or retry count exceeded
func httpReqWithRetry(path, method string, payload []byte, retries, timeoutMS, okStatus int) ([]byte, *http.Response, error) {
	r := retrier.New(
		retrier.ConstantBackoff(retries, time.Duration(timeoutMS)*time.Millisecond),
		nil,
	)

	var resp *http.Response
	var body []byte
	errs := errors.New("")

	err := r.Run(func() error {
		var err error
		body, resp, err = httpReq(
			os.Getenv("gateway_url")+path,
			method,
			bytes.NewBuffer(payload),
		)

		if err != nil {
			errs = errors.Wrap(errs, fmt.Sprintf("error executing request: %s\n", err))
			return errs
		}

		if resp.StatusCode != okStatus {
			errs = errors.Wrap(errs, fmt.Sprintf("expected status %d, got %d\n", okStatus, resp.StatusCode))
			return errs
		}

		return nil
	})

	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("call function %s failed after %d attempts\n", os.Getenv("gateway_url")+path, retries))
	}

	return body, resp, nil
}

func httpReq(url1, method string, reader io.Reader) ([]byte, *http.Response, error) {
	c := http.Client{}

	req, makeReqErr := http.NewRequest(method, url1, reader)
	if makeReqErr != nil {
		return nil, nil, fmt.Errorf("error with request %s ", makeReqErr)
	}

	res, callErr := c.Do(req)
	if callErr != nil {
		return nil, nil, fmt.Errorf("call error %s ", callErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
		bytesOut, err := ioutil.ReadAll(res.Body)

		return bytesOut, res, err
	}

	return nil, res, nil
}

func deployFunction(createRequest requests.CreateFunctionRequest) (*http.Response, error) {
	_, res, err := httpReq(
		os.Getenv("gateway_url")+"system/functions",
		"POST",
		makeReader(createRequest),
	)

	if err == nil && res.StatusCode == http.StatusOK {
		// check not already in deployed list
		found := false
		for _, f := range deployedFunctions {
			if f == createRequest.Service {
				found = true
			}
		}

		if !found {
			deployedFunctions = append(deployedFunctions, createRequest.Service)
		}
	}

	return res, err
}

func cleanupDeployedFunctions() {
	var cleanupFailed []string // list of functions unable to remove

	for _, f := range deployedFunctions {
		_, res, err := httpReq(
			os.Getenv("gateway_url")+"system/functions",
			"DELETE",
			makeReader(requests.DeleteFunctionRequest{
				FunctionName: f,
			}),
		)

		// if unsuccessful add to the failed list to try again later
		if err != nil || res.StatusCode != http.StatusOK {
			log.Println("Cleanup failed:", f, res.StatusCode)
			cleanupFailed = append(cleanupFailed, f)
		}
	}

	deployedFunctions = cleanupFailed
}

func listFunctions(t *testing.T) []requests.Function {
	bytesOut, res, err := httpReq(os.Getenv("gateway_url")+"system/functions", "GET", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
		return nil
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("error getting functions got status %d, wanted %d", res.StatusCode, http.StatusOK)
		t.Fail()
		return nil
	}

	fs := []requests.Function{}
	err = json.Unmarshal(bytesOut, &fs)
	if err != nil {
		t.Log(err)
		t.Fail()
		return nil
	}

	return fs
}
