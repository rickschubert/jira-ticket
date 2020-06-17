package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextOnLinks(t *testing.T) {
	input := "See failure on Jenkins: https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you https://jenkins.tray.io/job/qa-api-utied/ go https://jenkins.tray.io/"
	expectedOutput := []textChunk{
		{
			IsLink: false,
			Text:   "See failure on Jenkins: ",
		},
		{
			IsLink: true,
			Text:   "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
		},
		{
			IsLink: false,
			Text:   " and there you ",
		},
		{
			IsLink: true,
			Text:   "https://jenkins.tray.io/job/qa-api-utied/",
		},
		{
			IsLink: false,
			Text:   " go ",
		},
		{
			IsLink: true,
			Text:   "https://jenkins.tray.io/",
		},
	}
	assert.Equal(t, SplitTextOnLinks(input), expectedOutput)
}
