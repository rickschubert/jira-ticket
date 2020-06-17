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

func matchContainsActualText(match string) bool {
	return strings.Trim(match, " ") != ""
}

func extractAllLinksFromTextAndReplaceThemWithPlaceholders(text string) (linksInText []string, textWithPlaceHolders string) {
	var links []string
	textWithPlaceHolders = text

	linkRegex := regexp.MustCompile(`http(?:s)?://.+?(\s|$)`)
	results := linkRegex.FindAllStringSubmatch(text, -1)
	for i, result := range results {
		for _, match := range result {
			if matchContainsActualText(match) {
				linkOnly := strings.TrimRight(match, " ")
				links = append(links, linkOnly)
				textWithPlaceHolders = strings.Replace(textWithPlaceHolders, linkOnly, generatePlaceHolder(i), 1)
			}
		}
	}

	return links, textWithPlaceHolders
}

func splitStringOnMarkers(text string, markers []string) []string {
	var chunkedUpText []string
	chunkedUpText = []string{text}

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
			if isPlaceholder {
				chunkObjects = append(chunkObjects, textChunk{
					IsLink: true,
					Text:   linksInText[amountOfLinks],
				})
				amountOfLinks++
			} else {
				chunkObjects = append(chunkObjects, textChunk{
					IsLink: false,
					Text:   chunk,
				})
			}
		}
	}

	return chunkObjects
}

func SplitTextOnLinks(text string) []textChunk {
	links, textWithPlaceHolders := extractAllLinksFromTextAndReplaceThemWithPlaceholders(text)
	chunkedUpText := splitStringOnMarkers(textWithPlaceHolders, links)
	return convertListOfStringsIntoTextChunksShowingIfLinkOrNot(chunkedUpText, links)
}
