package utils

func GetAllowMethodsString(methods ...string) string {
	var result string
	var separator string = "|"

	for index, method := range methods {
		result += method

		if index < len(methods)-1 {
			result += separator
		}
	}

	return result
}
