package presenter

import (
	"Goose47/goman/internal/types"
	"testing"
)

func TestGetModuleInfo(t *testing.T) {
	tests := []struct {
		in   types.Module
		want string
	}{
		{types.Module{Name: "pkg", Uri: "/uri", Description: "desc"}, "Package pkg\nDocumentation: https://pkg.go.dev/uri\ndesc\n---"},
	}

	for _, tc := range tests {
		res := GetModuleInfo(tc.in)
		if res != tc.want {
			t.Errorf("Expected %s, got %s", tc.want, res)
		}
	}
}

func TestGetModuleNotFound(t *testing.T) {
	name := "name"
	res := GetModuleNotFound(name)
	want := "Module " + name + " is not found in standard library. :("

	if res != want {
		t.Errorf("Expected %s, got %s", want, res)
	}
}

func TestGetItemInfo(t *testing.T) {
	tests := []struct {
		in   types.Item
		want string
	}{
		{types.Item{Type: "type", Name: "pkg", Description: "desc", Signature: "sig", Example: "exa"}, "type pkg \nsig\n---\ndesc\nexa\n---"},
		{types.Item{Type: "type", Name: "pkg", Signature: "sig"}, "type pkg \nsig"},
	}

	for _, tc := range tests {
		res := GetItemInfo(tc.in)
		if res != tc.want {
			t.Errorf("Expected %s, got %s", tc.want, res)
		}
	}
}
