package manifest

import (
	"testing"
)

func TestGetBytesToEstimate(t *testing.T) {
	if GetBytesToEstimate(900) != "0.90 KB" {
		t.Errorf("Under 1 Clearly Fraction KB value %s doesn't match expected: 0.90 KB", GetBytesToEstimate(900))
	}

	if GetBytesToEstimate(1000) != "1.00 KB" {
		t.Errorf("Exact KB value %s doesn't match expected: 1.00 KB", GetBytesToEstimate(1000))
	}

	if GetBytesToEstimate(1001) != "1.00 KB" {
		t.Errorf("Near exact KB value %s doesn't match expected: 1.00 KB", GetBytesToEstimate(1001))
	}

	if GetBytesToEstimate(1251) != "1.25 KB" {
		t.Errorf("Fractional KB value %s doesn't match expected: 1.25 KB", GetBytesToEstimate(1250))
	}

	if GetBytesToEstimate(2000000) != "2.00 MB" {
		t.Errorf("Exact MB value %s doesn't match expected: 2.00 MB", GetBytesToEstimate(2000000))
	}

	if GetBytesToEstimate(2000011) != "2.00 MB" {
		t.Errorf("Near exact MB value %s doesn't match expected: 2.00 MB", GetBytesToEstimate(2000011))
	}

	if GetBytesToEstimate(2340011) != "2.34 MB" {
		t.Errorf("Clearly, fractional mb value %s doesn't match expected: 2.34 MB", GetBytesToEstimate(2340011))
	}

	if GetBytesToEstimate(3000000000) != "3.00 GB" {
		t.Errorf("Exact GB value %s doesn't match expected: 3.00 GB", GetBytesToEstimate(3000000000))
	}

	if GetBytesToEstimate(3000000011) != "3.00 GB" {
		t.Errorf("Near exact GB value %s doesn't match expected: 3.00 GB", GetBytesToEstimate(3000000011))
	}

	if GetBytesToEstimate(3470000011) != "3.47 GB" {
		t.Errorf("Clearly fractional GB value %s doesn't match expected: 3.47 GB", GetBytesToEstimate(3470000011))
	}

	if GetBytesToEstimate(5000000000000) != "5.00 TB" {
		t.Errorf("Exact TB value %s doesn't match expected: 3.00 GB", GetBytesToEstimate(5000000000000))
	}

	if GetBytesToEstimate(5000000000011) != "5.00 TB" {
		t.Errorf("Near exact TB value %s doesn't match expected: 5.00 TB", GetBytesToEstimate(5000000000011))
	}

	if GetBytesToEstimate(5630000000011) != "5.63 TB" {
		t.Errorf("Clearly fractional TB value %s doesn't match expected: 5.63 TB", GetBytesToEstimate(5630000000011))
	}
}
