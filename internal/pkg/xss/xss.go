package xss

import "github.com/microcosm-cc/bluemonday"

var Policy *bluemonday.Policy

func init() {
	Policy = bluemonday.UGCPolicy()
}

// Sanitize takes a string that contains a HTML fragment or document and applies the given policy allowlist.
// It returns an HTML string that has been sanitized by the policy or an empty string if an error has
// occurred (most likely as a consequence of extremely malformed input)
func Sanitize(input string) string {
	return Policy.Sanitize(input)
}
