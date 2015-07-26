package main

import (
	"math/rand"
	"testing"
	"time"
)

/*
func TestDrawAndDiscard(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := new(State)
	s.Init()

	vic := s.TotalVictory()
	if vic != 3 {
		t.Errorf("Victory points: %d. Want: 3", vic)
	} else {
		t.Logf("Victory points initially: %d", vic)
	}

	for i := 0; i < 100; i++ {
		s.DrawHand()
		s.Discard()
	}

	vic = s.TotalVictory()
	if vic != 3 {
		t.Errorf("Victory points: %d. Want: 3", vic)
	}
}
*/

func TestAddCard(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	s := new(State)
	s.Init()

	vic := s.TotalVictory()
	if vic != 3 {
		t.Errorf("Victory points: %d. Want: 3", vic)
	}

	for i := 0; i < 100; i++ {
		s.DrawHand()
		s.AddCardAndDiscardHand("copper")

		c := s.TotalCardsByName("copper")
		if c != 8+i {
			t.Errorf("[%d] Total coppers: %d. Want: %d\n", i, c, 8+i)
		}
	}

	cards := s.TotalCards()
	if cards != 110 {
		t.Errorf("Total cards: %d. Want: 110", cards)
	}

	coppers := s.TotalCardsByName("copper")
	if coppers != 107 {
		t.Errorf("Total coppers: %d. Want: 107", coppers)
	}
	estates := s.TotalCardsByName("estate")
	if estates != 3 {
		t.Errorf("Total estates: %d. Want: 3", estates)
	}
}
