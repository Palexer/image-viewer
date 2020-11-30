package main

import "github.com/disintegration/gift"

// Gifted extends gift.GIFT with a Remove and a Replace function
type Gifted struct {
	*gift.GIFT
}

// Remove removes the filter passed into the function
func (g *Gifted) Remove(filter gift.Filter) bool {
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

// Replace replaces the old with the new filter
func (g *Gifted) Replace(old, new gift.Filter) bool {
	for i, f := range g.Filters {
		if f == old {
			g.Filters[i] = new
			return true
		}
	}
	return false
}
