package router

import "testing"

func TestSetupRouter(t *testing.T) {

	r := SetupRouter()

	if r == nil {
		t.Errorf("router is nil")
	}

}
