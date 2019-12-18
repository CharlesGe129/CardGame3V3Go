package pkg

const (
	SPADE   CardColor = "spade"
	HEART   CardColor = "heart"
	CLUB    CardColor = "club"
	DIAMOND CardColor = "diamond"
)

var (
	mapCardName = map[uint32]string{
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "0",
		11: "J",
		12: "Q",
		13: "K",
		14: "A",
		15: "2",
		21: "小",
		22: "大",
	}
)

type CardColor string

type Card struct {
	Num   uint32
	Color CardColor
}

func (c *Card) Name() string {
	return mapCardName[c.Num]
}

func (c *Card) Cmp(other *Card) uint32 {
	return c.Num - other.Num
}
