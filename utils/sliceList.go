package utils


const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

type EqualFunc func(a, b interface{}) bool
type KeyFunc func(item interface{}) string

type SliceList struct {
	eFunc	EqualFunc
	kFunc 	KeyFunc
	items 	[]interface{}
	size	int
}

func NewSliceList(f EqualFunc, k KeyFunc) *SliceList {
	return &SliceList{eFunc:f, kFunc:k}
}

// Add appends a value at the end of the list
func (list *SliceList) Add(value interface{}) (index int) {
	list.growBy(1)
	list.items[list.size] = value
	index = list.size
	list.size++

	return
}

func (list *SliceList) Contains(item interface{}) bool {

	for _, it := range list.items {
		if list.eFunc(item, it) {
			return true
		}
	}
	return false
}

// Remove removes one or more elements from the list with the supplied indices.
func (list *SliceList) Remove(index int) {

	if !list.withinRange(index) {
		return
	}

	list.items[index] = nil                                    // cleanup reference
	copy(list.items[index:], list.items[(index + 1):list.size]) // shift to the left by one (slow operation, need ways to optimize this)
	list.size--

	list.shrink()
}


func (list *SliceList) Get(index int) (interface{}, bool) {

	if !list.withinRange(index) {
		return nil, false
	}

	return list.items[index], true
}


func (list *SliceList) Values() []interface{} {
	newElements := make([]interface{}, list.size, list.size)
	copy(newElements, list.items[:list.size])
	return newElements
}


//TODO: take a good look at this
func (list *SliceList) Keys() map[string]bool  {
	values := list.Values()
	keys := make(map[string]bool)
	for x := 0; x < len(values); x++ {
		v := list.kFunc(values[x])
		if "" == v {
			continue
		}
		keys[v] = true
	}

	return keys
}

// Check that the index is within bounds of the list
func (list *SliceList) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

func (list *SliceList) resize(cap int) {
	newElements := make([]interface{}, cap, cap)
	copy(newElements, list.items)
	list.items = newElements
}

func (list *SliceList) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.items)
	if list.size + n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity + n))
		list.resize(newCapacity)
	}
}


// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (list *SliceList) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(list.items)
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}

