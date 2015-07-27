// Create a simulator to test out big money strategy.
package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "Path to profile file")

type Node struct {
	S     State
	B     Board
	Moves int
}
type Nodes []Node

func (n Nodes) Len() int { return len(n) }
func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
func (n Nodes) Less(i, j int) bool {
	mi := n[i].Moves - n[i].S.NumProvinces()
	mj := n[j].Moves - n[j].S.NumProvinces()
	return mi < mj
}

func (n *Nodes) Push(x interface{}) {
	node := x.(Node)
	*n = append(*n, node)
}

func (n *Nodes) Pop() interface{} {
	old := *n
	l := len(old) - 1
	node := old[l]
	*n = old[0:l]
	return node
}

func shuffle(cards []string) {
	for i := range cards {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func addCard(b *Board, s *State, name string) bool {
	card, ok := GetCard(name)
	if !ok {
		panic("Invalid card")
		return false
	}

	if s.Value() < GetCost(card) {
		return false
	}

	card, ok = b.Get(name)
	if !ok {
		return false
	}

	s.AddCardAndDiscardHand(card)
	return true
}

func BuyPhase(s State, b Board, moves int) []Node {
	var result []Node
	{
		// No purchase made.
		ds := s.NewCopy()
		db := b.NewCopy()
		ds.Discard()
		nn := Node{S: ds, B: db, Moves: moves + 1}
		result = append(result, nn)
	}

	for _, name := range CardNames() {
		ds := s.NewCopy()
		db := b.NewCopy()
		if ok := addCard(&db, &ds, name); ok {
			nn := Node{S: ds, B: db, Moves: moves + 1}
			result = append(result, nn)
		}
	}
	return result
}

type TrashFunc func(picked []int)

func doPick(picked, left []int, fn TrashFunc) {
	if len(picked) > 0 {
		fn(picked)
	}
	if len(left) == 0 {
		return
	}

	for i := 0; i < len(left); i++ {
		if len(picked) > 0 {
			last := picked[len(picked)-1]
			if left[i] < last {
				continue
			}
		}

		cp := make([]int, len(picked)+1)
		copy(cp, picked)
		cp[len(picked)] = left[i]

		l := len(left) - 1
		cl := make([]int, l)
		copy(cl[0:i], left[:i])
		copy(cl[i:], left[i+1:])

		doPick(cp, cl, fn)
	}
}
func PickCardsToTrash(num int, fn TrashFunc) {
	var picked []int
	left := make([]int, num)
	for i := 0; i < num; i++ {
		left[i] = i
	}
	doPick(picked, left, fn)
}

func PlayTurn(n Node) []Node {
	// Draw Hand
	n.S.DrawHand()

	var result []Node
	// Play no Actions. Direct to buy phase scenario.
	result = append(result, BuyPhase(n.S, n.B, n.Moves)...)

	// Play Actions scenario.
	if has := n.S.CardInHand("smithy"); has {
		st := n.S.NewCopy()
		st.AddToHand(3)
		result = append(result, BuyPhase(st, n.B, n.Moves)...)
		// fmt.Printf("Cards in hand now: %s\n", n.S.StringHand())
	}

	if has := n.S.CardInHand("chapel"); has {
		hand := n.S.CopyHand()
		cidx := -1
		for idx, card := range hand {
			if card == "chapel" {
				cidx = idx
			}
		}
		if cidx == -1 {
			panic("Should happen")
		}

		alreadytrashed := make(map[string]bool)
		trashFn := func(picked []int) {
			for _, idx := range picked {
				if cidx == idx {
					return
				}
			}

			st := n.S.NewCopy()
			tr := st.TrashFromHand(picked)
			if al := alreadytrashed[tr]; !al {
				result = append(result, BuyPhase(st, n.B, n.Moves)...)
				alreadytrashed[tr] = true
			}
		}
		PickCardsToTrash(len(hand), trashFn)
	}

	return result
}

func MainLoop() {
	var board Board
	board.Init(2)

	var state State
	state.Init()

	n := Node{S: state, B: board}
	h := new(Nodes)
	heap.Init(h)
	heap.Push(h, n)

	attempt := 0
	for {
		attempt += 1

		tp := heap.Pop(h)
		node := tp.(Node)

		if attempt%1000 == 0 {
			fmt.Printf("Iter [%d] Moves: %d VP: %d Picks: %s\n",
				attempt, node.Moves, node.S.TotalVictory(), strings.Join(node.S.Picks, ","))
		}

		if node.S.NumProvinces() >= 4 {
			fmt.Printf("Reached 4 provinces in %d moves. Vic Points: %d\n",
				node.Moves, node.S.TotalVictory())
			node.S.Print()

			break
		}
		next := PlayTurn(node)
		for _, nn := range next {
			heap.Push(h, nn)
		}
		// q = append(q, next...)
		// sort.Sort(Nodes(q))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	CardInit()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	MainLoop()
	fmt.Println("DONE")
}
