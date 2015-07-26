// Create a simulator to test out big money strategy.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
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
	mi := n[i].Moves
	mj := n[j].Moves
	if mi == mj {
		vi := n[i].S.TotalVictory()
		vj := n[j].S.TotalVictory()
		return vi < vj
	}
	return mi > mj
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

func PlayTurn(n Node) []Node {
	var result []Node
	n.S.DrawHand()

	for _, name := range CardNames() {
		ds := n.S.NewCopy()
		db := n.B.NewCopy()
		if ok := addCard(&db, &ds, name); ok {
			nn := Node{S: ds, B: db, Moves: n.Moves + 1}
			result = append(result, nn)
		}
	}
	return result
}

func MainLoop() {
	var board Board
	board.Init(2)

	var state State
	state.Init()

	// considered := make(map[string]bool)

	n := Node{S: state, B: board}
	var q []Node
	q = append(q, n)

	attempt := 0
	for {
		if attempt > 5000 {
			break
		}
		fmt.Printf("\nAttempt: %d\n", attempt)
		attempt += 1

		lq := len(q) - 1
		node := q[lq]
		q = q[:lq]
		node.S.Print()
		fmt.Printf("In moves: %d\n", node.Moves)

		/*
			ps := node.S.PickState()
			if already := considered[ps]; already {
				fmt.Println("Already considered", ps)
				continue
			} else {
				fmt.Println("New consideration", ps)
				considered[ps] = true
			}
		*/

		if node.S.TotalVictory() >= 24 {
			fmt.Printf("Reached victory in %d moves\n", node.Moves)

			break
		}
		next := PlayTurn(node)
		q = append(q, next...)
		sort.Sort(Nodes(q))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	CardInit()

	flag.Parse()
	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	MainLoop()
	fmt.Println("DONE")
}
