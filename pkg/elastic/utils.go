package elastic

func stringToInterface(s []string) (result []interface{}) {
	result = make([]interface{}, len(s))
	for i, v := range s {
		result[i] = v
	}
	return
}
