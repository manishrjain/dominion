# Dominion
Compute strategies to win Dominion - the popular [deck-building game](https://en.wikipedia.org/wiki/Dominion_(card_game)).

Use this tool to find the ideal strategy to win Dominion, given a bunch of action cards (or without action cards). The program would iterate over all possible states using [A\* algorithm](https://en.wikipedia.org/wiki/A*_search_algorithm), to find the solution with the minimum moves. In addition, it runs multiple times with different `rand.Seed` values so you can easily determine a solution which can consistenly guarantee wins.

Supported Action cards:
- Smithy
- Chapel

### Installation
```go
go get github.com/manishrjain/dominion
go install github.com/manishrjain/dominion
```

### Usage
```bash
$ dominion
```

### Typical Output
Output with `chapel` and `smithy` action cards, with target to achieve 4 `provinces` (for 2 player game), with 10 iterations:
```
Found solution in moves: [10 11 11 11 11 11 11 11 11 12]  // <- Shows sorted list of number of moves required by each winning strategy.
11 silver smithy chapel gold copper estate NONE province province province province  // FORMAT: Number of Moves, followed by Card picked during each move.
11 smithy silver NONE gold province smithy province silver estate province province  // NONE means no card was picked in that move.
11 chapel smithy copper NONE gold copper gold province province province province
11 silver chapel smithy NONE copper gold silver province province province province
10 smithy silver gold gold smithy province silver province province province
11 chapel smithy silver NONE silver gold province province silver province province
11 smithy silver NONE gold gold gold province NONE province province province
11 chapel silver smithy NONE gold NONE province silver province province province
12 chapel smithy silver NONE NONE gold smithy gold province province province province
11 silver silver smithy gold gold gold province NONE province province province
```

**What this means:** This shows that opening with `silver smithy chapel` looks like the best move, showing up in 5 different solutions. After this, a `smithy silver` is the next best move, showing up in 3 different solutions. In no scenario does `big money + (only) chapel` come up as solution. This also [settles the debate](http://forum.dominionstrategy.com/index.php?topic=636.0) about `big money + chapel` v/s `big money + smithy`.

### Modify
You can choose which cards you want to play with, by modifying the code in `board.go`.
```go
// Comment any of these cards out to see a strategy without that card.
b.cards = append(b.cards, "smithy")
b.cards = append(b.cards, "chapel")
```

### TODO
- The program is pretty basic right now, needs support for more action cards.
- Currently cards are defined in both `card.go`, and `board.go`. We need a better way to express all the cards in a cosolidated location.
- Better way to add code logic required by new cards.
