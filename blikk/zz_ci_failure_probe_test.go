package blikk

import "testing"

func TestCIFailureCommentProbe(t *testing.T) {
	t.Fatal("deliberate failure to verify the PR failure comment; this test is removed in the next commit")
}
