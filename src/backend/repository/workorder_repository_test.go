package repository

import "testing"

func TestIsWorkOrderInProgress(t *testing.T) {
	if !isWorkOrderInProgress(" Progress ") {
		t.Fatal("progress work order must be protected from deletion")
	}
	if isWorkOrderInProgress("completed") {
		t.Fatal("completed work order may be deleted")
	}
}
