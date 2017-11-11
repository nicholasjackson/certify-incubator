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

// randomTestRun runs the tests in random order to esnure that a previous test
// suite is not inadvertinelyt setting up any conditions for the next suite
func randomTestRun(t *testing.T, tests []func(t *testing.T)) {
	shuffle(tests)

	for _, test := range tests {
		test(t)
	}
}

// suffle the order of the test slice
func shuffle(vals []func(t *testing.T)) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
