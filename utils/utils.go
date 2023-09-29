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
