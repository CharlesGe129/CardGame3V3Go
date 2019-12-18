package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Players         [6]Player
	FinishedPlayers map[int]struct{}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame() (g Game) {
	for i := 0; i < 3; i++ {
		g.Players[2*i].Team = 1
	}
	g.Players[0].Type = PlayerTypeUser
	g.FinishedPlayers = make(map[int]struct{})
	return
}

func (g *Game) Start() {
	g.AssignCards()
	var curShot Shot
	curPlayer := rand.Intn(6)
	bigPlayer := curPlayer
	numPasses := g.ResetNumPasses()
	for !g.isFinished() {
		if numPasses == 0 {
			curShot = Shot{}
			g.ShowCards()
			numPasses = g.ResetNumPasses()
		}
		shot := g.Players[curPlayer].NextShot(curShot)
		if shot.Type != ShotTypePass {
			curShot = shot
			bigPlayer = curPlayer
			numPasses = g.ResetNumPasses()
		} else {
			numPasses -= 1
		}
		fmt.Printf("Player%d: %s, numPasses=%d\n", curPlayer, shot.Cards, numPasses)
		if g.Players[curPlayer].IsFinished() {
			fmt.Printf("Player%d finishes\n", curPlayer)
			g.FinishedPlayers[curPlayer] = struct{}{}
			bigPlayer = g.NextPlayer(bigPlayer)
		}
		curPlayer = g.NextPlayer(curPlayer)
	}
}

func (g *Game) ResetNumPasses() int {
	return 6 - 1 - len(g.FinishedPlayers)
}

func (g *Game) isFinished() bool {
	_, ok1 := g.FinishedPlayers[0]
	_, ok2 := g.FinishedPlayers[2]
	_, ok3 := g.FinishedPlayers[4]
	if ok1 && ok2 && ok3 {
		fmt.Printf("Team 1 wins! \n%v\n", g.FinishedPlayers)
		return true
	}
	_, ok1 = g.FinishedPlayers[1]
	_, ok2 = g.FinishedPlayers[3]
	_, ok3 = g.FinishedPlayers[5]
	if ok1 && ok2 && ok3 {
		fmt.Printf("Team 2 wins! \n%v\n", g.FinishedPlayers)
		return true
	}
	return false
}

func (g *Game) NextPlayer(cur int) int {
	for {
		if cur == 5 {
			cur = 0
		} else {
			cur += 1
		}
		if _, ok := g.FinishedPlayers[cur]; !ok {
			return cur
		}
	}
}

func initialCard(num int) Cards {
	n := uint32(num)
	return Cards{
		Card{
			Num:   n,
			Color: SPADE,
		},
		Card{
			Num:   n,
			Color: HEART,
		},
		Card{
			Num:   n,
			Color: CLUB,
		},
		Card{
			Num:   n,
			Color: DIAMOND,
		},
	}
}

func initialCards() (cards Cards) {
	for num := 0; num < 3; num++ {
		for i := 3; i <= 15; i++ {
			cards = append(cards, initialCard(i)...)
		}
		cards = append(cards, Card{Num: 21})
		cards = append(cards, Card{Num: 22})
	}
	return
}

func (g *Game) AssignCards() {
	cards := initialCards()
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	for i := 0; i < len(cards); i++ {
		g.Players[i%6].AddCard(cards[i])
	}
	g.ShowCards()
}

func (g *Game) ShowCards() {
	fmt.Println("========== all cards ==========")
	for i := 0; i < len(g.Players); i++ {
		if _, ok := g.FinishedPlayers[i]; !ok {
			fmt.Printf("Player%d: ", i)
			g.Players[i].ShowCards()
		}
	}
	fmt.Println("===============================")
}
