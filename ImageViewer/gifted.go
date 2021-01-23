package main

import "github.com/disintegration/gift"

// gifted extends gift.GIFT with a Remove and a Replace function
type gifted struct {
	*gift.GIFT
}

// remove removes the filter passed into the function
func (g *gifted) remove(filter gift.Filter) bool {
	for i, f := range g.Filters {
		if f == filter {
			copy(g.Filters[i:], g.Filters[i+1:])     // Shift a[i+1:] left one index.
			g.Filters[len(g.Filters)-1] = nil        // Erase last element (write zero value).
			g.Filters = g.Filters[:len(g.Filters)-1] // Truncate slice.
			return true
		}
	}
	return false
}

// replace replaces the old with the new filter. If the old filter doesn't exist, it applys the new filter.
func (g *gifted) replace(old, new gift.Filter) bool {
	for i, f := range g.Filters {
		if f == old {
			g.Filters[i] = new
			return true
		}
	}
	g.Add(new)
	return false
}
