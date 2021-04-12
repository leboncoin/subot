package globals_test

import (
	"github.com/leboncoin/subot/pkg/globals"
	"testing"
)

func TestCompareEmptyTimeStamps(t *testing.T) {
	first := ""
	second := ""
	_, err := globals.CompareTimeStamps(first, second)
	if err == nil {
		t.Errorf("Empty timestamp comparision impossible: %d is nil", err)
	}
}

func TestCompareNegativeTimeStamps(t *testing.T) {
	first := "20945"
	second := "1"
	diff, err := globals.CompareTimeStamps(first, second)
	if err != nil {
		t.Errorf("Error while comparing timestamps: %d is not nil", err)
	}
	if diff > 0 {
		t.Errorf("Negative difference: %f < %d", diff, 0)
	}
}

func TestCompareTimeStamps(t *testing.T) {
	first := "20945"
	second := "120945"
	diff, err := globals.CompareTimeStamps(first, second)
	if err != nil {
		t.Errorf("Error while comparing timestamps: %d is not nil", err)
	}
	if diff != 100000. {
		t.Errorf("Incorrect difference: %f is not equal to %f", diff, 100000.)
	}
}