package email

import (
	"slices"
	"strings"
)

// MaskEmail replaces all but the first character in each part of an email address with asterisks,
// e.g. j***_d**@e******.c**
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	localPart := parts[0]
	domainPart := parts[1]

	return MaskString(localPart) + "@" + MaskString(domainPart)
}

// MaskString replaces all but the first character in each part of a string with asterisks,
// where parts are delineated by '-' (hyphen), '.' (period), or '_' (underscore)
func MaskString(s string) string {
	if s == "" {
		return ""
	}
	const chars = "-._"
	newStrings := []string{}
	parts := strings.FieldsFunc(s, func(r rune) bool { return slices.Contains([]rune(chars), r) })
	for _, p := range parts {
		masked := "*"
		if len(p) > 1 {
			masked = p[0:1] + strings.Repeat("*", len(p)-1)
		}
		newStrings = append(newStrings, masked)
	}

	i := strings.IndexAny(s, chars)
	if i == -1 {
		return newStrings[0]
	}
	return strings.Join(newStrings, s[i:i+1])
}
