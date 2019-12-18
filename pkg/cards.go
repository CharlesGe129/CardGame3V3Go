package pkg

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type NumSorter []Card

func (s NumSorter) Len() int           { return len(s) }
func (s NumSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s NumSorter) Less(i, j int) bool { return s[i].Num < s[j].Num }

type Cards []Card

func (c Cards) SplitInGroups() map[int][]string {
	result := make(map[int][]string)

	cards := c
	sort.Sort(NumSorter(cards))
	cards = append(cards, Card{Num: 1})

	curCard := ""
	curCount := 0
	for _, card := range cards {
		name := card.Name()
		if curCount == 0 {
			curCard = name
			curCount += 1
		} else if name != "" && strings.HasPrefix(curCard, name) {
			curCount += 1
			curCard += name
			if curCount == 5 {
				result[5] = append(result[5], curCard)
				curCount = 0
			}
		} else {
			var idx int
			if strings.HasPrefix(curCard, "小") || strings.HasPrefix(curCard, "大") {
				idx = 6
			} else {
				idx = curCount
			}
			result[idx] = append(result[idx], curCard)
			curCard = name
			curCount = 1
		}
	}
	return result
}

func (c *Cards) Get5Level() (level uint32, large uint32, err error) {
	cards := *c
	if len(cards) != 5 {
		return 0, 0, fmt.Errorf("bad 5-cards %q", c)
	}
	sort.Sort(NumSorter(cards))

	// straight, flush, full house, four, flush straight, five
	cardsStr := cards.String()
	largeCount := 0
	largeNum := uint32(0)
	for _, card := range cards {
		count := strings.Count(cardsStr, card.Name())
		if largeCount < count {
			largeCount = count
			large = card.Num
		}
		if largeNum < card.Num {
			largeNum = card.Num
		}
	}
	// five, four, full house
	switch largeCount {
	case 5:
		return 5, large, nil
	case 4:
		return 3, large, nil
	case 3:
		return 2, large, nil
	}
	if c.checkStraight() && c.checkFlush() {
		return 4, largeNum, nil
	} else if c.checkFlush() {
		return 1, largeNum, nil
	} else if c.checkStraight() {
		return 0, largeNum, nil
	} else {
		return 0, 0, fmt.Errorf("bad 5-cards %q", c)
	}
}

func (c *Cards) checkStraight() bool {
	cards := *c
	sort.Sort(NumSorter(cards))
	return len(cards) == 5 &&
		cards[0].Num+1 == cards[1].Num &&
		cards[1].Num+1 == cards[2].Num &&
		cards[2].Num+1 == cards[3].Num &&
		cards[3].Num+1 == cards[4].Num
}

func (c *Cards) checkFlush() bool {
	cards := *c
	sort.Sort(NumSorter(cards))
	return len(cards) == 5 &&
		cards[0].Color == cards[1].Color &&
		cards[1].Color == cards[2].Color &&
		cards[2].Color == cards[3].Color &&
		cards[3].Color == cards[4].Color
}

func (c Cards) String() (str string) {
	for _, card := range c {
		str += card.Name()
	}
	return
}

func (c *Cards) Larger(o *Cards) (bool, error) {
	cards := *c
	sort.Sort(NumSorter(cards))
	others := *o
	sort.Sort(NumSorter(others))
	if len(cards) != len(others) {
		return false, fmt.Errorf("cards length not equal: %q, %q", cards, others)
	}
	level1, large1, err := cards.validate()
	if err != nil {
		return false, err
	}
	level2, large2, err := others.validate()
	if err != nil {
		return false, err
	}
	if level1 > level2 {
		return true, nil
	} else if level1 < level2 {
		return false, nil
	} else {
		return large1 > large2, nil
	}
}

func (c *Cards) validate() (level uint32, large uint32, err error) {
	cards := *c
	sort.Sort(NumSorter(cards))
	switch len(cards) {
	case 1, 2, 3:
		single := cards[0].Num
		for _, card := range cards {
			if card.Num != single {
				return 0, 0, fmt.Errorf("bad cards %q", cards)
			}
		}
		return 0, single, nil
	case 5:
		return c.Get5Level()
	default:
		return 0, 0, fmt.Errorf("bad cards: %q", cards)
	}
}

func (c Cards) Contains(target Card) bool {
	for _, card := range c {
		if card.Num == target.Num && card.Color == target.Color {
			return true
		}
	}
	return false
}

func (c Cards) Delete(target Card) (rs Cards) {
	var deleted bool
	for _, card := range c {
		// TODO: card color
		if !deleted && card.Num == target.Num {
			deleted = true
		} else {
			rs = append(rs, card)
		}
	}
	return
}

func (c *Cards) Copy() (rs Cards) {
	for _, card := range *c {
		rs = append(rs, Card{
			Num:   card.Num,
			Color: card.Color,
		})
	}
	return
}

func CardStrToCards(cardStr string) (cards Cards) {
	// TODO: Card Color
	for _, char := range cardStr {
		var card Card
		if 50 < char && char <= 57 {
			// 3~9
			n, err := strconv.ParseUint(string(char), 10, 32)
			if err != nil {
				panic(err)
			}
			card.Num = uint32(n)
		} else {
			switch char {
			case 48:
				// 10
				card.Num = 10
			case 50:
				// 2
				card.Num = 15
			case 74:
				// J
				card.Num = 11
			case 81:
				// Q
				card.Num = 12
			case 75:
				// K
				card.Num = 13
			case 65:
				// A
				card.Num = 14
			case 23567:
				// 小
				card.Num = 21
			case 22823:
				// 大
				card.Num = 22
			}
		}
		cards = append(cards, card)
	}
	sort.Sort(NumSorter(cards))
	return
}
