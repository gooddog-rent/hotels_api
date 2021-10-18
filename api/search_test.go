package api

import "testing"

// TODO:
func TestValidBuildSuffixTree(t *testing.T) {
	l := Locations{}

	l.BuildSuffixTree()

	if l.Tree == nil {
		t.Errorf("BuildSuffixTree must be valid tree! got: %v", l.Tree)
	}
}
