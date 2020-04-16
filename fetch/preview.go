package fetch

// Preview fetches and extracts Details for a URL
func Preview(url string) (*Details, error) {
	return any(url)
}
