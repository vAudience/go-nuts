package gonuts

func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringSliceIndexOf(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func StringSliceRemoveString(max int, sourceSlice []string, stringToRemove string) []string {
	if max == 0 {
		return sourceSlice
	}
	found := []int{}
	for i, a := range sourceSlice {
		if a == stringToRemove {
			found = append(found, i)
		}
	}
	var resultSlice []string = sourceSlice
	for _, idx := range found {
		if max == -1 || max > 0 {
			resultSlice = StringSliceRemoveIndex(resultSlice, idx)
		}
	}
	return resultSlice
}

func StringSliceRemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func StringSlicesHaveSameContent(source []string, compare []string) bool {
	if len(source) != len(compare) {
		return false
	}
	for _, item := range source {
		found := false
		for _, cmp := range compare {
			if item == cmp {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
