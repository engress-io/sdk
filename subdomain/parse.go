package subdomain

// ParseSubdomain splits a subdomain label into prefix and an 8-character
// lowercase alphanumeric suffix. The suffix is the last 8 characters.
func ParseSubdomain(label string) (prefix, suffix string, ok bool) {
	if len(label) < 9 {
		return "", "", false
	}
	suffix = label[len(label)-8:]
	for _, c := range suffix {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return "", "", false
		}
	}
	return label[:len(label)-8], suffix, true
}