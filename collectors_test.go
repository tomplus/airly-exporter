package main

import (
	"testing"
)

func TestRegisterCollectors(t *testing.T) {

	pc := PromCollectors{}
	pc.RegisterCollectors()

	if pc.countTotal == nil {
		t.Errorf("countTotal not inited")
	}
	if pc.errorTotal == nil {
		t.Errorf("errorTotal not inited")
	}
}
