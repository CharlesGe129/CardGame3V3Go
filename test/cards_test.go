package test

import (
	"testing"

	"CardGame3V3Go/pkg"
	"github.com/stretchr/testify/require"
)

func TestCards_SplitInGroups(t *testing.T) {
	cards := pkg.Cards{
		pkg.Card{Num: 3}, pkg.Card{Num: 4}, pkg.Card{Num: 5},
	}
	expected := [6][]string{
		{"3", "4", "5"},
	}
	require.Equal(t, expected, cards.SplitInGroups())

	cards = pkg.Cards{
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
		pkg.Card{Num: 4}, pkg.Card{Num: 4}, pkg.Card{Num: 4}, pkg.Card{Num: 4},
	}
	expected = [6][]string{
		nil,
		nil,
		{"333"},
		{"4444"},
	}
	require.Equal(t, expected, cards.SplitInGroups())

	cards = pkg.Cards{
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
		pkg.Card{Num: 21}, pkg.Card{Num: 21}, pkg.Card{Num: 22},
	}
	expected = [6][]string{
		nil,
		nil,
		{"333"},
		nil,
		nil,
		{"小小", "大"},
	}
	require.Equal(t, expected, cards.SplitInGroups())

	cards = pkg.Cards{
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
		pkg.Card{Num: 3}, pkg.Card{Num: 3}, pkg.Card{Num: 3},
	}
	expected = [6][]string{
		nil,
		{"33"},
		nil,
		nil,
		{"33333", "33333"},
	}
	require.Equal(t, expected, cards.SplitInGroups())
}

func TestCards_Get5Level(t *testing.T) {
	cards := pkg.Cards{
		pkg.Card{Num: 3}, pkg.Card{Num: 4}, pkg.Card{Num: 5},
	}
	level, large, err := cards.Get5Level()
	require.Error(t, err)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 4, Color: pkg.SPADE},
		pkg.Card{Num: 5, Color: pkg.HEART},
		pkg.Card{Num: 6, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.Error(t, err)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 4, Color: pkg.SPADE},
		pkg.Card{Num: 5, Color: pkg.HEART},
		pkg.Card{Num: 6, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(0), level)
	require.Equal(t, uint32(7), large)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 4, Color: pkg.SPADE},
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 6, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(1), level)
	require.Equal(t, uint32(7), large)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.HEART},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(2), level)
	require.Equal(t, uint32(7), large)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.HEART},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(3), level)
	require.Equal(t, uint32(7), large)

	cards = pkg.Cards{
		pkg.Card{Num: 3, Color: pkg.SPADE},
		pkg.Card{Num: 4, Color: pkg.SPADE},
		pkg.Card{Num: 5, Color: pkg.SPADE},
		pkg.Card{Num: 6, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(4), level)
	require.Equal(t, uint32(7), large)

	cards = pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.HEART},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	level, large, err = cards.Get5Level()
	require.NoError(t, err)
	require.Equal(t, uint32(5), level)
	require.Equal(t, uint32(7), large)
}

func TestCards_TestLarger(t *testing.T) {
	cards := pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	others := pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	isLarger, err := cards.Larger(&others)
	require.NoError(t, err)
	require.False(t, isLarger)

	cards = pkg.Cards{
		pkg.Card{Num: 5, Color: pkg.SPADE},
	}
	others = pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	isLarger, err = cards.Larger(&others)
	require.NoError(t, err)
	require.False(t, isLarger)

	cards = pkg.Cards{
		pkg.Card{Num: 8, Color: pkg.SPADE},
	}
	others = pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	isLarger, err = cards.Larger(&others)
	require.NoError(t, err)
	require.True(t, isLarger)

	cards = pkg.Cards{
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
	}
	others = pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	isLarger, err = cards.Larger(&others)
	require.NoError(t, err)
	require.True(t, isLarger)

	cards = pkg.Cards{
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
	}
	others = pkg.Cards{
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
		pkg.Card{Num: 7, Color: pkg.SPADE},
	}
	isLarger, err = cards.Larger(&others)
	require.NoError(t, err)
	require.True(t, isLarger)

	cards = pkg.Cards{
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
		pkg.Card{Num: 8, Color: pkg.SPADE},
	}
	others = pkg.Cards{
		pkg.Card{Num: 9, Color: pkg.SPADE},
		pkg.Card{Num: 9, Color: pkg.SPADE},
		pkg.Card{Num: 9, Color: pkg.SPADE},
		pkg.Card{Num: 9, Color: pkg.SPADE},
		pkg.Card{Num: 5, Color: pkg.SPADE},
	}
	isLarger, err = cards.Larger(&others)
	require.NoError(t, err)
	require.True(t, isLarger)
}
