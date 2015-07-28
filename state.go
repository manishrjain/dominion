package main

import (
	"fmt"
	"sort"
	"strings"
)

type State struct {
	hand    []string
	discard []string
	deck    []string
	Picks   []string

	victory int
	numProv int
}

func (s *State) NewCopy() State {
	ds := State{}
	ds.hand = make([]string, len(s.hand))
	ds.discard = make([]string, len(s.discard))
	ds.deck = make([]string, len(s.deck))
	ds.Picks = make([]string, len(s.Picks))

	copy(ds.hand, s.hand)
	copy(ds.discard, s.discard)
	copy(ds.deck, s.deck)
	copy(ds.Picks, s.Picks)

	ds.victory = -1
	ds.numProv = -1
	return ds
}

func (s *State) PickState() string {
	return strings.Join(s.Picks, " ")
}

func (s *State) Print() {
	var c []string
	c = append(c, "HAND")
	c = append(c, s.hand...)
	c = append(c, "DISCARD")
	c = append(c, s.discard...)
	c = append(c, "DECK")
	c = append(c, s.deck...)
	c = append(c, "PICKS")
	c = append(c, s.Picks...)

	fmt.Println(strings.Join(c, " "))
	fmt.Printf("Victory points: %d\n", s.TotalVictory())
}

func (s *State) StringHand() string {
	return strings.Join(s.hand, ",")
}

func (s *State) Init() {
	for i := 0; i < 3; i++ {
		s.deck = append(s.deck, "estate")
	}
	for i := 0; i < 7; i++ {
		s.deck = append(s.deck, "copper")
	}
	shuffle(s.deck)
	if s.TotalCards() != 10 {
		panic("Invalid initialization of state")
	}
	s.victory = -1
	s.numProv = -1
}

func (s *State) drawCards(num int) []string {
	numcards := len(s.deck) + len(s.discard)

	if len(s.deck) < num {
		s.discard = append(s.discard, s.deck...)
		s.deck = s.deck[:0]

		shuffle(s.discard)

		s.deck = make([]string, len(s.discard))
		copy(s.deck, s.discard)

		// s.deck = s.discard
		s.discard = s.discard[:0]

		if len(s.discard) != 0 {
			panic("Coders fault")
		}
	}

	if len(s.deck) < num {
		// Chapel strategy can trash a lot of cards.
		num = len(s.deck)
	}

	cards := make([]string, num)
	copy(cards, s.deck[0:num])
	s.deck = s.deck[num:]

	if len(s.deck)+len(s.discard) != numcards-num {
		panic("Coders fault again")
	}

	return cards
}

func (s *State) DrawHand() {
	if len(s.hand) != 0 {
		panic("Already have hand")
		// fmt.Println("Already have hand")
		return
	}
	cards := s.drawCards(5)
	s.hand = make([]string, len(cards))
	copy(s.hand, cards)
}

func (s *State) AddToHand(num int) {
	cards := s.drawCards(num)
	s.hand = append(s.hand, cards...)
}

func (s *State) CardInHand(name string) bool {
	for _, card := range s.hand {
		if card == name {
			return true
		}
	}
	return false
}

func (s *State) CopyHand() []string {
	cp := make([]string, len(s.hand))
	copy(cp, s.hand)
	return cp
}

func (s *State) TrashFromHand(indices []int) string {
	var trashed []string
	for _, idx := range indices {
		if idx >= len(s.hand) {
			panic("Invalid index")
		}
		trashed = append(trashed, s.hand[idx])
		s.hand[idx] = ""
	}

	var newh []string
	for _, card := range s.hand {
		if card != "" {
			newh = append(newh, card)
		}
	}
	s.hand = newh

	sort.Sort(sort.StringSlice(trashed))
	return strings.Join(trashed, ",")
}

func (s *State) TrashUselessCards() {
	var newh []string
	for _, card := range s.hand {
		if card == "copper" || card == "estate" {
			continue
		}
		newh = append(newh, card)
	}
	s.hand = newh
}

func (s *State) AddCardAndDiscardHand(c string) {
	s.victory = -1
	s.numProv = -1
	s.discard = append(s.discard, c)
	s.Picks = append(s.Picks, c)
	s.Discard()
}

func (s *State) Discard() {
	s.discard = append(s.discard, s.hand...)
	s.hand = s.hand[:0]
	if len(s.hand) != 0 {
		panic("Hand should be zero.")
	}
}

func (s *State) Value() int {
	total := 0
	for _, card := range s.hand {
		total += GetValue(card)
	}
	return total
}

func (s *State) TotalVictory() int {
	if s.victory != -1 {
		return s.victory
	}

	total := 0
	for _, card := range s.discard {
		total += GetVictory(card)
	}
	for _, card := range s.deck {
		total += GetVictory(card)
	}
	for _, card := range s.hand {
		total += GetVictory(card)
	}
	s.victory = total
	return total
}

func (s *State) NumProvinces() int {
	if s.numProv != -1 {
		return s.numProv
	}
	s.numProv = s.TotalCardsByName("province")
	return s.numProv
}

func (s *State) TotalCards() int {
	total := 0
	total += len(s.discard)
	total += len(s.hand)
	total += len(s.deck)
	return total
}

func (s *State) TotalCardsByName(name string) int {
	total := 0
	for _, card := range s.discard {
		if card == name {
			total += 1
		}
	}
	for _, card := range s.deck {
		if card == name {
			total += 1
		}
	}
	for _, card := range s.hand {
		if card == name {
			total += 1
		}
	}
	return total
}
