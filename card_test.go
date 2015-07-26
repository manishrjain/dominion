package main

import "testing"

func TestCostAndVictory(t *testing.T) {
	cost := GetCost("estate")
	if cost != 2 {
		t.Errorf("Wrong answer")
	}

	vic := GetVictory("estate")
	if vic != 1 {
		t.Errorf("Wrong answer")
	}

	vic = GetVictory("copper")
	if vic != 0 {
		t.Errorf("Wrong answer")
	}

	val := GetValue("gold")
	if val != 3 {
		t.Errorf("Wrong answer")
	}

	val = GetValue("province")
	if val != 0 {
		t.Errorf("Wrong answer")
	}
}
