package external

import (
	"testing"
)

func TestGetTankAverages(t *testing.T) {
	averages, err := GetTankAverages()
	if err != nil {
		t.Fatal(err)
	}

	if len(averages) == 0 {
		t.Fatal("no averages returned")
	}
}
