package orderedset

// Strings is an ordered set for strings.
type Strings struct {
	elements []string
	index    map[string]interface{}
}

// Get returns the array of this set.
// Do not change the returned array.
func (set *Strings) Get() []string {
	return set.elements
}

// Append adds an item in order.
// NOTE: this is not concurrent safe.
func (set *Strings) Append(s string) {
	_, exists := set.index[s]
	if exists {
		return
	}
	if set.index == nil {
		set.index = make(map[string]interface{})
	}
	set.index[s] = nil
	set.elements = append(set.elements, s)
}
