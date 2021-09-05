package utils


func RemoveDuplicateElement(slice []string) []string {
	result := make([]string, 0, len(slice))
	temp := map[string]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}


func StrInArray(str string, strArr []string) bool {
	for _, v := range strArr {
		if str == v {
			return true
		}
	}
	return false
}
