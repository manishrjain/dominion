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
			fmt.Printf("Iter [%d] Moves: %d Victory Points: %d\n",
				attempt, node.Moves, node.S.TotalVictory())
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
