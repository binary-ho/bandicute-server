package util

import (
	"html"
	"regexp"
)

func SanitizeContent(content string) string {
	// Remove HTML tags
	removedHtmlTags := removeHtmlTags(content)

	// Decode HTML entities
	decodedEntities := html.UnescapeString(removedHtmlTags)

	// Retain structural elements (like bullet points, new lines)
	normalizedWhitespace := regexp.MustCompile(`\s+`).ReplaceAllString(decodedEntities, " ")

	// Remove unwanted characters but retain some context (e.g., "-", "•", "*")
	cleanedContent := regexp.MustCompile(`[^\p{Hangul}\p{L}\d\s.,!?•*()-]`).ReplaceAllString(normalizedWhitespace, "")

	return cleanedContent
}

// removeHtmlTags removes all HTML tags from the input string but retains line breaks for better context.
func removeHtmlTags(content string) string {
	decoded := html.UnescapeString(content)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(decoded, "\n") // Replace tags with newlines to retain structure.
}
