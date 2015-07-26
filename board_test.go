package main

import "testing"

func TestBoardGet(t *testing.T) {
	var b Board
	b.Init(2)

	for i := 0; i < 8; i++ {
		card, ok := b.Get("estate")
		if !ok || card != "estate" {
			t.Errorf("Can't get estate")
		}
	}
	card, ok := b.Get("estate")
	if ok {
		t.Errorf("Extra estate found")
	}

	card, ok = b.Get("province")
	if !ok || card != "province" {
		t.Error("Can't get province")
	}
}
