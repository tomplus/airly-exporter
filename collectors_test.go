package main

import (
	"testing"
)

func TestCollectors(t *testing.T) {

	pc := PromCollectors{}
	pc.RegisterCollectors()

	if pc.countTotal == nil {
		t.Errorf("countTotal not inited")
	}
	if pc.errorTotal == nil {
		t.Errorf("errorTotal not inited")
	}

	measurements := AllMeasurements{}
	pc.SetMeasurements("1234", measurements)
}
