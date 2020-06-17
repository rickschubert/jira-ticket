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

func SplitTextOnLinks(text string) []textChunk {
	var linksInText []string
	var textWithPlaceHolders string

	textWithPlaceHolders = text
	linkRegex := regexp.MustCompile(`http(?:s)?://.+?(\s|$)`)
	results := linkRegex.FindAllStringSubmatch(text, -1)
	for i, result := range results {
		for _, match := range result {
			if strings.Trim(match, " ") != "" {
				linkOnly := strings.TrimRight(match, " ")
				linksInText = append(linksInText, linkOnly)
				textWithPlaceHolders = strings.Replace(textWithPlaceHolders, linkOnly, generatePlaceHolder(i), 1)
			}
		}
	}

	var chunkedUpText []string
	chunkedUpText = []string{textWithPlaceHolders}

	for i := 0; i < len(linksInText); i++ {
		placeholder := generatePlaceHolder(i)
		lastChunk := chunkedUpText[len(chunkedUpText)-1]
		splits := strings.Split(lastChunk, placeholder)
		chunkedUpText[len(chunkedUpText)-1] = splits[0]
		chunkedUpText = append(chunkedUpText, placeholder)
		chunkedUpText = append(chunkedUpText, splits[1])
	}

	var chunkObjects []textChunk
	var amountOfLinks int

	for _, chunk := range chunkedUpText {
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
