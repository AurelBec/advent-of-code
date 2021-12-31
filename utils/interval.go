package utils

// Interval represents a simple Interval of values
type Interval[K number] struct {
	Min, Max K
}

// NewInterval is a quick way to get an Interval without worrying about specification and value assignation
func NewInterval[K number](min, max K) Interval[K] {
	return Interval[K]{Min: min, Max: max}
}

// Len returns the Interval range
func (interval Interval[K]) Len() K {
	return interval.Max - interval.Min + 1
}

// Merge merges the other value (number or Interval) into this one
// note: overlap is not checked
func (lhs *Interval[K]) Merge(value any) {
	switch value := value.(type) {
	case K:
		lhs.Min = Min(lhs.Min, value)
		lhs.Max = Max(lhs.Max, value)
	case Interval[K]:
		lhs.Min = Min(lhs.Min, value.Min)
		lhs.Max = Max(lhs.Max, value.Max)
	}
}

// Shift moves interval's bounds by the specified delta value
func (interval *Interval[K]) Shift(delta K) {
	interval.Min += delta
	interval.Max += delta
}

// Touches tells if the 2 Interval are at most delta unit apart
func (lhs Interval[K]) Touches(rhs Interval[K], delta K) bool {
	return Max(lhs.Min-delta, rhs.Min) <= Min(lhs.Max+delta, rhs.Max)
}

// Overlaps tells if the 2 Interval have a common part (i.e. touches with a 0 delta)
func (lhs Interval[K]) Overlaps(rhs Interval[K]) bool {
	return lhs.Touches(rhs, 0)
}

// Intersection returns the intersecting interval between 2 ones
func (lhs Interval[K]) Intersection(rhs Interval[K]) Interval[K] {
	return NewInterval(Max(lhs.Min, rhs.Min), Min(lhs.Max, rhs.Max))
}

// Split separates interval in two around value
// parameter to specify if value is sent to left or right
// if value is not present, initial interval is sent to left, right is empty
func (interval Interval[K]) Split(value K, left bool) (Interval[K], Interval[K]) {
	if !interval.Contains(value) {
		return interval, Interval[K]{}
	} else if left {
		return NewInterval(interval.Min, value), NewInterval(value+1, interval.Max)
	} else {
		return NewInterval(interval.Min, value-1), NewInterval(value, interval.Max)
	}
}

// Contains returns whether this Interval fully contains the other value (number or Interval)
func (interval Interval[K]) Contains(value any) bool {
	switch value := value.(type) {
	case K:
		return interval.Min <= value && value <= interval.Max
	case Interval[K]:
		return interval.Min <= value.Min && value.Max <= interval.Max
	}
	return false
}

// Intervals represents an array of Interval
// no interval should be overlapping an other, they are merged if it's the case
type Intervals[K number] struct {
	UnorderedArray[Interval[K]]
}

// Insert inserts a new Interval into the range
// if the new interval does not touch any other, it's inserted
// else, it is merged with the previous, which is removed and the merged one is inserted etc.
// delta parameter allow to define tolerance for the touch
func (intervals *Intervals[K]) Insert(value any, delta ...K) {
	switch value := value.(type) {
	case K:
		intervals.Insert(NewInterval(value, value), delta...)
	case Interval[K]:
		deltaK := K(0)
		if len(delta) > 0 {
			deltaK = delta[0]
		}

		for i, interval := range intervals.UnorderedArray {
			// merge intervals if touching: remove the old one and retry insertion with the merge
			if interval.Touches(value, deltaK) {
				intervals.Remove(i)
				value.Merge(interval)
				intervals.Insert(value, delta...)
				return
			}
		}

		// finish by inserting the range
		intervals.UnorderedArray.Insert(value)
	}
}
