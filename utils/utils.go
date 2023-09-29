package utils

// UrlBuilder util takes possible parts of url and returns full url
func UrlBuilder(proto string, domain string, port string) string {
	var url string
	if len(proto) > 0 {
		url = url + proto
	} else {
		url = url + "http://"
	}
	url = url + domain
	if len(port) > 0 {
		url = url + ":" + port
	}
	return url
}

// IntToBool util takes an int, returns bool (0=false, 1=true)
func IntToBool(i int) bool {
	return i != 0
}

// BoolToInt util takes a bool, returns an int 0=false, 1=true)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
