package util

import (
	"testing"
)

func TestGetFullName(t *testing.T) {
	parent := "astra foo"
	child := "bar"
	expected := parent + " " + child
	actual := GetFullName(parent, child)
	if expected != actual {
		t.Errorf("test failed, expected %s, got %s", expected, actual)
	}
}
