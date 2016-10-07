package structures

// @todo: see if we can make this more generic using interfaces

type UniqueSlice struct {
	data []string
}

/**
 * Create a new UniqueSlice and populate it with the provided slice
 */
func NewUniqueSlice(slice []string) *UniqueSlice {
	ret := new(UniqueSlice)

	for _, v := range slice {
		ret.Append(v)
	}

	return ret
}

// Append an element to the slice
func (slice *UniqueSlice) Append(item string) bool {
	for i := range slice.data {
		if slice.data[i] == item {
			return false
		}
	}

	slice.data = append(slice.data, item)

	return true;
}

// Get the element at an index
func (slice *UniqueSlice) Get(index int) string {
	return slice.data[index]
}

func (slice *UniqueSlice) Set(index int, value string) {
	slice.data[index] = value
}

// Return all the data as a []string
func (slice *UniqueSlice) ToSlice() []string {
	return slice.data
}
// Get the number of elements in UniqueSlice
func (slice *UniqueSlice) Length() int {
	return len(slice.data)
}

func (slice *UniqueSlice) IndexOf(item string) int {
	for i := range slice.data {
		if slice.data[i] == item {
			return i
		}
 	}

	return -1
}
