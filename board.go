package main

type Board struct {
	cards []string
}

func (b *Board) NewCopy() Board {
	db := Board{}
	db.cards = make([]string, len(b.cards))
	copy(db.cards, b.cards)
	return db
}

func (b *Board) Init(players int) {
	num := 8
	if players > 2 {
		num = 12
	}
	for i := 0; i < num; i++ {
		b.cards = append(b.cards, "estate")
		b.cards = append(b.cards, "dutchy")
		b.cards = append(b.cards, "province")
		// b.cards = append(b.cards, "smithy")
		b.cards = append(b.cards, "chapel")
	}
	for i := 0; i < 30; i++ {
		b.cards = append(b.cards, "copper")
		b.cards = append(b.cards, "silver")
		b.cards = append(b.cards, "gold")
	}
}

func (b *Board) Get(name string) (rcard string, rok bool) {
	idx := -1
	for i, card := range b.cards {
		if card == name {
			idx = i
			rcard = card
			break
		}
	}
	if idx == -1 {
		return rcard, false
	}
	l := len(b.cards) - 1
	b.cards[idx], b.cards[l] = b.cards[l], b.cards[idx]
	b.cards = b.cards[0:l]
	return rcard, true
}
