package common

// GetKeyByValue gets a key by its value. The first key with a matching
// value is returned. If value is not found, def is returned
func GetKeyByValue(mapping map[string]string, value, def string) string {
	for key, _value := range mapping {
		if _value == value {
			return key
		}
	}
	return def
}
