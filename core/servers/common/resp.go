package common

// SingleObjectResp creates standard response format for only a single object
func SingleObjectResp(objectType string, attr map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"type":       objectType,
			"id":         attr["id"],
			"attributes": attr,
		},
	}
}

// MultiObjectResp creates standard response format for multiple objects
func MultiObjectResp(objectType string, attrs []map[string]interface{}) map[string]interface{} {
	var data = make([]map[string]interface{}, 0)
	for _, attr := range attrs {
		data = append(data, map[string]interface{}{
			"type":       objectType,
			"id":         attr["id"],
			"attributes": attr,
		})
	}
	return map[string]interface{}{
		"data": data,
	}
}
