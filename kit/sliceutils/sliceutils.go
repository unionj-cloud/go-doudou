package sliceutils

func StringSlice2InterfaceSlice(strSlice []string) []interface{} {
	ret := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v
	}
	return ret
}

func InterfaceSlice2StringSlice(strSlice []interface{}) []string {
	ret := make([]string, len(strSlice))
	for i, v := range strSlice {
		ret[i] = v.(string)
	}
	return ret
}

func Contains(src []interface{}, test interface{}) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

func StringContains(src []string, test string) bool {
	for _, item := range src {
		if item == test {
			return true
		}
	}
	return false
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
