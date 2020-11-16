package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextOnLinks(t *testing.T) {
	input := "See failure on system: https://system.host.com/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you https://system.host.com/job/qa-api-utied/ go https://system.host.com/"
	expectedOutput := []textChunk{
		{
			IsLink: false,
			Text:   "See failure on system: ",
		},
		{
			IsLink: true,
			Text:   "https://system.host.com/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
		},
		{
			IsLink: false,
			Text:   " and there you ",
		},
		{
			IsLink: true,
			Text:   "https://system.host.com/job/qa-api-utied/",
		},
		{
			IsLink: false,
			Text:   " go ",
		},
		{
			IsLink: true,
			Text:   "https://system.host.com/",
		},
	}
	assert.Equal(t, SplitTextOnLinks(input), expectedOutput)
}

func TestSplitTextOnLinksLinkOnly(t *testing.T) {
	input := "https://system.host.com/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/"
	expectedOutput := []textChunk{
		{
			IsLink: true,
			Text:   "https://system.host.com/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
		},
	}
	assert.Equal(t, SplitTextOnLinks(input), expectedOutput)
}

func TestSplitTextOnLinksTextOnly(t *testing.T) {
	input := "See failure on system"
	expectedOutput := []textChunk{
		{
			IsLink: false,
			Text:   "See failure on system",
		},
	}
	assert.Equal(t, SplitTextOnLinks(input), expectedOutput)
}
