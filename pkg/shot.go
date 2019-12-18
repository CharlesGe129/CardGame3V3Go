package pkg

const (
	ShotTypePass  ShotType = 0
	ShotTypeOne   ShotType = 1
	ShotTypeTwo   ShotType = 2
	ShotTypeThree ShotType = 3
	ShotTypeFive  ShotType = 5
)

type ShotType uint32

type Shot struct {
	Cards Cards
	Type  ShotType
	Team  uint32
}

func (s Shot) String() string {
	switch s.Type {
	case ShotTypePass:
		return "pass"
	case ShotTypeOne, ShotTypeTwo, ShotTypeThree, ShotTypeFive:
		return s.Cards.String()
	default:
		return "shot type not yet implemented"
	}
}

func (s *Shot) CheckLarger(nextCards Cards) bool {
	isLarger, err := nextCards.Larger(&s.Cards)
	if err != nil {
		panic(err)
	}
	return isLarger
}
