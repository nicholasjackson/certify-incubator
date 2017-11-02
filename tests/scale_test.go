package tests

import "testing"

func TestScale(t *testing.T) {
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Test_ScaleRequestScalesFunction", scaleRequestScalesFunction)
	})

	// <tear-down code>

}

func scaleRequestScalesFunction(t *testing.T) {

}
