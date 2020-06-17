package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextOnLinks(t *testing.T) {
	input := "See failure on Jenkins: https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you https://jenkins.tray.io/job/qa-api-utied/ go https://jenkins.tray.io/"
	expectedOutput := []textChunk{
		{
			isLink: false,
			Text:   "See failure on Jenkins: ",
		},
		{
			isLink: true,
			Text:   "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
		},
		{
			isLink: false,
			Text:   " and there you ",
		},
		{
			isLink: true,
			Text:   "https://jenkins.tray.io/job/qa-api-utied/",
		},
		{
			isLink: false,
			Text:   " go ",
		},
		{
			isLink: true,
			Text:   "https://jenkins.tray.io/",
		},
	}
	assert.Equal(t, SplitTextOnLinks(input), expectedOutput)
}
