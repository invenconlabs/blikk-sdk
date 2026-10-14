package blikk

import "testing"

func TestCIFailureCommentProbe(t *testing.T) {
	t.Fatal("deliberate failure to verify the PR report; removed in the next commit")
}
