package generator

import (
	"math"
	"math/rand"
	"testing"
)

func TestWeightedChoiceEmpty(t *testing.T) {
	choices := map[string]int{}

	wc := newWeightedChoice(choices)
	if got := wc.next(); got != "" {
		t.Errorf("when choices are empty, it should return empty string, got %s", got)
	}
}

func TestWeightedChoice_next(t *testing.T) {
	rand.Seed(42)

	const (
		rounds = 10e6
		eps    = 0.01
	)

	choices := map[string]int{
		"first":  100,
		"second": 50,
		"third":  80,
	}

	wc := newWeightedChoice(choices)

	counts := map[string]int{}
	for i := 0; i < rounds; i++ {
		n := wc.next()
		counts[n]++
	}

	// test if the resulting distribution diverges from probabilistic distribution more then the allowed epsilon
	for k := range choices {
		g := float64(counts[k]) / float64(rounds)
		if math.Abs(g-float64(choices[k])/float64(wc.sum)) > eps {
			t.Errorf("failed next method test element %s: got %v expected %v +/- %v", k, g, choices[k], eps)
		}
	}
}
