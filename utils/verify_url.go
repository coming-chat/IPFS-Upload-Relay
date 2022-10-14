package utils

import "net/url"

// VerifyURL : Return if a URL is valid (not empty & matches URL spec)
func VerifyURL(u string) bool {
	if u == "" {
		return false
	}
	_, err := url.Parse(u)
	if err != nil {
		return false
	}

	return true
}
