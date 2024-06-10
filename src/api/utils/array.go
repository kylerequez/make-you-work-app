package utils

func HasValues(a []interface{}, b []string) bool {
	for _, strB := range b {
		for _, v := range a {
			if strVal, ok := v.(string); ok && strVal == strB {
				return true
			}
		}
	}
	return false
}
