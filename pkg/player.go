package pkg

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	PlayerTypeUser     PlayerType = "user"
	PlayerTypeNormalAI PlayerType = "ai"
)

type PlayerType string

type Player struct {
	Cards []Card
	Team  uint32
	Type  PlayerType
}

func (p *Player) ShowCards() {
	fmt.Printf("%s, len=%d\n", Cards(p.Cards), len(p.Cards))
	fmt.Printf("All 5 combos: ")
	cardsList := p.FormFive()
	for _, cards := range cardsList {
		fmt.Printf("%s ", cards)
	}
	fmt.Println("")
}

func (p *Player) AddCard(card Card) {
	p.Cards = append(p.Cards, card)
	if len(p.Cards) == 27 {
		sort.Sort(NumSorter(p.Cards))
	}
}

func (p *Player) NextShot(curShot Shot) Shot {
	if p.Type == PlayerTypeUser {
		for {
			next, err := p.ShotByInput(curShot)
			if err != nil {
				fmt.Println("Oops, wrong cards! Please try again:")
			} else {
				return next
			}
		}
	}
	if curShot.Type == ShotTypePass {
		return p.NewRoundShot()
	} else if p.CheckFriendShot(curShot) {
		return Shot{
			Team: p.Team,
		}
	} else {
		return p.ShotByType(curShot)
	}
}

func (p *Player) CheckFriendShot(curShot Shot) bool {
	if curShot.Team != p.Team {
		return false
	}
	cards := curShot.Cards
	if curShot.Type < 5 {
		return !(3 <= cards[0].Num && cards[0].Num <= 9)
	}
	level, large, err := cards.Get5Level()
	if err != nil {
		panic(err)
	}
	switch level {
	case 0, 1, 2:
		return false
	case 3:
		return !(3 <= large && large <= 9)
	case 4, 5:
		return true
	default:
		panic(fmt.Errorf("bad 5 level: %d", level))
	}
}

func (p *Player) ShotByInput(curShot Shot) (Shot, error) {
	fmt.Printf("Current cards: ")
	p.ShowCards()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please type your next shot, friend=%v: \n", curShot.Team == p.Team)
	cardStr, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	cardStr = strings.Split(cardStr, "\n")[0]
	if strings.HasPrefix("pass", strings.ToLower(cardStr)) {
		return Shot{Team: p.Team}, nil
	}
	cardList := CardStrToCards(cardStr)
	if curShot.Type == ShotTypePass || len(cardList) == int(curShot.Type) && p.ValidateCards(cardList) && curShot.CheckLarger(cardList) {
		p.RemoveCards(cardList)
		return Shot{
			Cards: cardList,
			Type:  ShotType(len(cardList)),
			Team:  p.Team,
		}, nil
	}
	return p.ShotByInput(curShot)
}

func (p *Player) ValidateCards(shotCards Cards) bool {
	switch len(shotCards) {
	case 1, 2, 3:
		n := shotCards[0].Num
		for _, card := range shotCards {
			if n != card.Num {
				return false
			}
		}
	case 5:
	default:
		return false
	}
	curCards := shotCards.Copy()
	for _, shotCard := range shotCards {
		if !curCards.Contains(shotCard) {
			return false
		}
		curCards.Delete(shotCard)
	}
	return true
}

func (p *Player) ShotByType(curShot Shot) Shot {
	if curShot.Type == ShotTypeFive {
		cards5Combos := p.FormFive()
		for _, cards := range cards5Combos {
			if len(cards) != 5 {
				break
			}
			if curShot.CheckLarger(cards) {
				p.RemoveCards(cards)
				return Shot{
					Cards: cards,
					Type:  5,
					Team:  p.Team,
				}
			}
		}
	} else {
		// type 1, 2, 3
		splitCards := Cards(p.Cards).SplitInGroups()
		for cardType, v := range splitCards {
			if ShotType(cardType) != curShot.Type {
				continue
			}
			for _, cardsStr := range v {
				cards := CardStrToCards(cardsStr)
				if curShot.CheckLarger(cards) {
					p.RemoveCards(cards)
					return Shot{
						Cards: cards,
						Type:  ShotType(len(cards)),
						Team:  p.Team,
					}
				}
			}
		}
	}
	return Shot{
		Team: p.Team,
	}
}

func (p *Player) NewRoundShot() Shot {
	cardsFive := p.FormFive()
	if cardsFive != nil {
		// type 5
		if cardsFive[0].String() == "小" || cardsFive[0].String() == "大" {
			fmt.Println(123)
		}
		if len(cardsFive[0]) == 5 {
			cards := cardsFive[0]
			p.RemoveCards(cards)
			return Shot{
				Cards: cards,
				Type:  5,
				Team:  p.Team,
			}
		}
	}
	// type 1, 2, 3
	splitCards := Cards(p.Cards).SplitInGroups()
	for cardType, v := range splitCards {
		for _, cardsStr := range v {
			cards := CardStrToCards(cardsStr)
			p.RemoveCards(cards)
			var t ShotType
			if cardType == 6 {
				t = ShotType(len(cardsStr) / 3)
			} else {
				t = ShotType(len(cardsStr))
			}
			return Shot{
				Cards: cards,
				Type:  t,
				Team:  p.Team,
			}
		}
	}
	panic("NewRoundShot() panic")
	return Shot{}
}

func (p *Player) RemoveCards(cards Cards) {
	for _, card := range cards {
		p.Cards = Cards(p.Cards).Delete(card)
	}
}

func (p *Player) IsFinished() bool {
	return len(p.Cards) == 0
}

func (p *Player) FormFive() (cardsList []Cards) {
	splitCards := Cards(p.Cards).SplitInGroups()
	var cardsRemains []Cards
	var l int
	// 1 + 4
	if len(splitCards[1]) >= len(splitCards[4]) {
		l = len(splitCards[4])
	} else {
		l = len(splitCards[1])
		for i := l; i < len(splitCards[4]); i++ {
			cardsRemains = append(cardsRemains, CardStrToCards(splitCards[4][i]))
		}
	}
	for i := 0; i < l; i++ {
		cardsList = append(cardsList, CardStrToCards(splitCards[1][i]+splitCards[4][i]))
	}
	// 2 + 3
	if len(splitCards[2]) >= len(splitCards[3]) {
		l = len(splitCards[3])
	} else {
		l = len(splitCards[2])
		for i := l; i < len(splitCards[3]); i++ {
			cardsRemains = append(cardsRemains, CardStrToCards(splitCards[3][i]))
		}
	}
	for i := 0; i < l; i++ {
		cardsList = append(cardsList, CardStrToCards(splitCards[2][i]+splitCards[3][i]))
	}
	// 5
	for _, cardsStr := range splitCards[5] {
		cardsList = append(cardsList, CardStrToCards(cardsStr))
	}
	// 6
	for _, cardsStr := range splitCards[6] {
		cardsList = append(cardsList, CardStrToCards(cardsStr))
	}
	// remain
	for _, cards := range cardsRemains {
		cardsList = append(cardsList, cards)
	}
	return
}
