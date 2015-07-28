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
	"sort"
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "Path to profile file")
var numiter = flag.Int("numiter", 10, "Number of iterations to run")

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
		ds.Picks = append(ds.Picks, "NONE")
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
		st := n.S.NewCopy()
		st.TrashUselessCards()
		// If I trash all possible combinations of cards in hand, the number of
		// states generated are way too many to run within a single machine having
		// 25GB of RAM. So, instead, just delete Copper + Estate, aka Useless Cards.
		result = append(result, BuyPhase(st, n.B, n.Moves)...)
	}

	return result
}

func MainLoop() (int, string) {
	var board Board
	board.Init(2)

	var state State
	state.Init()

	n := Node{S: state, B: board}
	h := new(Nodes)
	heap.Init(h)
	heap.Push(h, n)

	repetitive := false
	considered := make(map[string]bool)
	attempt := 0
	for {
		attempt += 1

		tp := heap.Pop(h)
		node := tp.(Node)

		ps := node.S.PickState()
		if already := considered[ps]; already {
			// fmt.Printf("Already considered: %s\n", ps)
			repetitive = true
			continue
		} else {
			considered[ps] = true
		}

		if attempt%1000 == 0 {
			fmt.Printf("Iter [%d] Moves: %d VP: %d Picks: %s\n",
				attempt, node.Moves, node.S.TotalVictory(), strings.Join(node.S.Picks, ","))
		}

		if node.S.NumProvinces() >= 4 {
			fmt.Printf("Reached 4 provinces in %d moves. Vic Points: %d\n",
				node.Moves, node.S.TotalVictory())
			node.S.Print()
			fmt.Println("Was repetitive", repetitive)

			return node.Moves, node.S.PickState()
		}
		next := PlayTurn(node)
		for _, nn := range next {
			ps = nn.S.PickState()
			if already := considered[ps]; !already {
				heap.Push(h, nn)
			}
		}
	}
}

func main() {
	// rand.Seed(time.Now().UnixNano())
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

	var moves []int
	var sols []string
	for i := 0; i < *numiter; i++ {
		rand.Seed(int64(i))
		m, s := MainLoop()
		moves = append(moves, m)
		sols = append(sols, s)
	}

	dm := make([]int, len(moves))
	copy(dm, moves)
	sort.Sort(sort.IntSlice(dm))
	fmt.Println("Found solution in moves:", dm)
	for i, sol := range sols {
		fmt.Println(moves[i], sol)
	}
	fmt.Println("DONE")
}
