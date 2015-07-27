package main

import "fmt"

var allcards []string

func CardInit() {
	a := [...]string{"smithy", "chapel", "estate", "dutchy", "province", "copper", "silver", "gold"}

	allcards = a[:]
	fmt.Printf("All Cards: %+v\n", allcards)
}

func GetCost(name string) int {
	switch name {
	case "estate":
		return 2
	case "dutchy":
		return 5
	case "province":
		return 8
	case "copper":
		return 0
	case "silver":
		return 3
	case "gold":
		return 6
	case "smithy":
		return 4
	case "chapel":
		return 2
	default:
		panic("Invalid name")
	}
	return 0
}

func GetValue(name string) int {
	switch name {
	case "copper":
		return 1
	case "silver":
		return 2
	case "gold":
		return 3
	default:
		return 0
	}
	return 0
}

func GetVictory(name string) int {
	switch name {
	case "estate":
		return 1
	case "dutchy":
		return 3
	case "province":
		return 6
	default:
		return 0
	}
	return 0
}

func GetCard(name string) (c string, r bool) {
	for _, card := range allcards {
		if card == name {
			return card, true
		}
	}
	return c, false
}

func CardNames() []string {
	return allcards
}
