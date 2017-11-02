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

	"github.com/alexellis/faas/gateway/requests"
)

var deployedFunctions []string

func init() {
	deployedFunctions = make([]string, 0)
}

func makeReader(input interface{}) *bytes.Buffer {
	res, _ := json.Marshal(input)
	return bytes.NewBuffer(res)
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
