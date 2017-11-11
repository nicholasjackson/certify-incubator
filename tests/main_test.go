package tests

import (
	"math/rand"
	"testing"
	"time"
)

func Test_Certify(t *testing.T) {
	randomTestRun(t, []func(t *testing.T){
		deploy,
		list,
		invoke,
		delete,
		alerts,
	})
}

// run the tests in random order to esnure there are no pre-dependencies between
// test runs
func randomTestRun(t *testing.T, tests []func(t *testing.T)) {
	shuffle(tests)

	for _, te := range tests {
		te(t)
	}
}

func deploy(t *testing.T) {
	t.Run("deploy", func(t *testing.T) {
		t.Run("Deploy function passing custom environment variables", passingCustomEnvVars)
	})

	cleanupDeployedFunctions()
}

func invoke(t *testing.T) {
	t.Run("invoke", func(t *testing.T) {
		t.Run("Invoke a function and check response", basicInvoke)
	})
	cleanupDeployedFunctions()
}

func delete(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		t.Run("Check it is possible to delete a function", basicDelete)
	})
}

func list(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		t.Run("When two functions are deployed list should return two functions", when2FunctionsDeployedListReturns2Functions)
	})

	cleanupDeployedFunctions()
}

func alerts(t *testing.T) {
	t.Run("alerts", func(t *testing.T) {
		t.Run("Scale alert scales function", scaleUpAlertScalesFunction)
	})

	cleanupDeployedFunctions()
}

// suffle the test order in random order
func shuffle(vals []func(t *testing.T)) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
