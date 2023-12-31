package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	AllDifferent = iota
	OnePair
	TwoPairs
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Game struct {
	hand []int
	bid  int
}

type Games []Game

func (g Games) Len() int {
	return len(g)
}

// sort games of the samve val
func (g Games) Less(i, j int) bool {
	for k := range g[i].hand {
		if g[i].hand[k] > g[j].hand[k] {
			return false
		}

		if g[i].hand[k] < g[j].hand[k] {
			return true
		}
	}

	panic("unexpected equality")
}

func (g Games) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(bytes), "\n")
	games, err := linesToGames(lines)
	if err != nil {
		log.Fatal(err)
	}

	gamesByHandVal := groupGamesByHandValue(games)
	for v := range gamesByHandVal {
		sort.Sort(gamesByHandVal[v])
	}

	fmt.Println(countPoints(gamesByHandVal))

}

func linesToGames(lines []string) (games Games, err error) {
	games = make([]Game, len(lines))

	for i := range games {
		split := strings.Split(lines[i], " ")

		games[i].bid, err = strconv.Atoi(split[1])
		if err != nil {
			return games, err
		}

		for _, char := range strings.Split(split[0], "") {
			var val int

			switch char {
			case "A":
				val = 14
			case "K":
				val = 13
			case "Q":
				val = 12
			case "J":
				val = 11
			case "T":
				val = 10
			default:
				val, err = strconv.Atoi(char)
				if err != nil {
					return games, err
				}
			}

			games[i].hand = append(games[i].hand, val)
		}
	}

	return games, nil
}

func getHandValue(hand []int) int {
	seeds := countSeeds(hand)

	if len(seeds) == 1 {
		return FiveOfAKind
	}

	if len(seeds) == 5 {
		return AllDifferent
	}

	if len(seeds) == 4 {
		return OnePair
	}

	if len(seeds) == 2 {
		if _, ok := mapContains(seeds, 4); ok {
			return FourOfAKind
		}
		return FullHouse
	}

	if len(seeds) == 3 {
		if _, ok := mapContains(seeds, 3); ok {
			return ThreeOfAKind
		}
		return TwoPairs
	}

	log.Fatal("unknown case", seeds)
	return 0
}

func countSeeds[V comparable](slice []V) map[V]int {
	set := map[V]int{}

	for _, v := range slice {
		if _, ok := set[v]; ok {
			set[v]++
		} else {
			set[v] = 1
		}
	}

	return set
}

func groupGamesByHandValue(games Games) map[int]Games {
	grouped := map[int]Games{}

	for i := range games {
		val := getHandValue(games[i].hand)
		if _, ok := grouped[val]; ok {
			grouped[val] = append(grouped[val], games[i])
		} else {
			grouped[val] = []Game{games[i]}
		}
	}

	return grouped
}

func countPoints(orederdGames map[int]Games) int {
	total := 0
	counter := 1
	for i := 0; i < 7; i++ {
		for j := range orederdGames[i] {
			total += (counter * orederdGames[i][j].bid)
			counter++
		}
	}

	return total
}

func mapContains[K, V comparable](m map[K]V, v V, excluding ...K) (K, bool) {
	for k := range m {
		if slices.Contains(excluding, k) {
			continue
		}

		if val, ok := m[k]; ok {
			if val == v {
				return k, true
			}
		}
	}

	var zero K
	return zero, false
}
