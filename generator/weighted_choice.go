package generator

import "math/rand"

// weightedChoice is a structure that provides random choice from the map src of choices, fairly distributing the probability
// according to the choice's weight.
type weightedChoice struct {
	src map[string]int
	sum int
}

// newWeightedChoice returns a new Weighted choice, summing up the choices weights inside to speed up next function.
func newWeightedChoice(src map[string]int) *weightedChoice {
	sum := 0
	for _, v := range src {
		sum += v
	}
	return &weightedChoice{
		sum: sum,
		src: src,
	}
}

// next returns a random choice, respecting the weights of the available choices.
func (wc *weightedChoice) next() (k string) {
	if len(wc.src) == 0 {
		return ""
	}

	r := rand.Intn(wc.sum)

	var acum, v int
	for k, v = range wc.src {
		if r >= acum && r <= (acum+v) {
			return k
		}
		acum += v
	}

	return k
}
