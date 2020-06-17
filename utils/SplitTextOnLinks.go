package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type textChunk struct {
	IsLink bool
	Text   string
}

const placeholderMarker = "_--__-_"

func generatePlaceHolder(index int) string {
	return fmt.Sprintf("%s%d%s", placeholderMarker, index, placeholderMarker)
}

func extractLinksFromTextAndReplaceWithPlaceholders(text string) (textWithPlaceHolders string, linksInText []string) {
	var links []string
	textWithPlaceHolders = text

	linkRegex := regexp.MustCompile(`http(?:s)?://.+?(\s|$)`)
	results := linkRegex.FindAllStringSubmatch(text, -1)
	for i, result := range results {
		for _, match := range result {
			link := strings.Trim(match, " ")
			if link != "" {
				links = append(links, link)
				textWithPlaceHolders = strings.Replace(textWithPlaceHolders, link, generatePlaceHolder(i), 1)
			}
		}
	}

	return textWithPlaceHolders, links
}

func splitStringOnMarkers(text string, markers []string) []string {
	chunkedUpText := []string{text}

	for i := 0; i < len(markers); i++ {
		placeholder := generatePlaceHolder(i)
		lastChunk := chunkedUpText[len(chunkedUpText)-1]
		splits := strings.Split(lastChunk, placeholder)
		chunkedUpText[len(chunkedUpText)-1] = splits[0]
		chunkedUpText = append(chunkedUpText, placeholder)
		chunkedUpText = append(chunkedUpText, splits[1])
	}

	return chunkedUpText
}

func convertListOfStringsIntoTextChunksShowingIfLinkOrNot(textChunks []string, linksInText []string) []textChunk {
	var chunkObjects []textChunk
	var amountOfLinks int

	for _, chunk := range textChunks {
		if chunk != "" {
			isPlaceholder, _ := regexp.Match(placeholderMarker, []byte(chunk))
			var currentTextChunk textChunk
			if isPlaceholder {
				currentTextChunk = textChunk{
					IsLink: true,
					Text:   linksInText[amountOfLinks],
				}
				amountOfLinks++
			} else {
				currentTextChunk = textChunk{
					IsLink: false,
					Text:   chunk,
				}
			}
			chunkObjects = append(chunkObjects, currentTextChunk)
		}
	}

	return chunkObjects
}

func SplitTextOnLinks(text string) []textChunk {
	textWithPlaceHolders, links := extractLinksFromTextAndReplaceWithPlaceholders(text)
	chunkedUpText := splitStringOnMarkers(textWithPlaceHolders, links)
	return convertListOfStringsIntoTextChunksShowingIfLinkOrNot(chunkedUpText, links)
}
