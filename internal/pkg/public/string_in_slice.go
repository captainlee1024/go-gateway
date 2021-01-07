package public

// InStringSlice 判断字符串是否在string切片中
func InStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}
