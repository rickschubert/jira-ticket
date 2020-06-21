package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkifyJiraDescriptionObjects(t *testing.T) {
	inputObjects := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you go",
				},
			},
		},
	}
	expected := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: ",
				},
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
					Marks: []mark{
						{
							Type: "link",
							Attrs: Attr{
								Href: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
							},
						},
					},
				},
				{
					Type: "text",
					Text: " and there you go",
				},
			},
		},
	}
	assert.Equal(t, linkifyJiraApiDescriptionObjects(inputObjects), expected)
}

func TestLinkifyJiraDescriptionObjectsMultipleParagraphs(t *testing.T) {
	inputObjects := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you go",
				},
			},
		},
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "hardBreak",
				},
			},
		},
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you go",
				},
			},
		},
	}
	expected := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: ",
				},
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
					Marks: []mark{
						{
							Type: "link",
							Attrs: Attr{
								Href: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
							},
						},
					},
				},
				{
					Type: "text",
					Text: " and there you go",
				},
			},
		},
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "hardBreak",
				},
			},
		},
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "See failure on Jenkins: ",
				},
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
					Marks: []mark{
						{
							Type: "link",
							Attrs: Attr{
								Href: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
							},
						},
					},
				},
				{
					Type: "text",
					Text: " and there you go",
				},
			},
		},
	}
	assert.Equal(t, linkifyJiraApiDescriptionObjects(inputObjects), expected)
}

func TestLinkifyJiraDescriptionObjectsLinksAtBeginningAndEnd(t *testing.T) {
	inputObjects := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/ and there you go https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
				},
			},
		},
	}
	expected := []content{
		{
			Type: "paragraph",
			Content: []paragraphContent{
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
					Marks: []mark{
						{
							Type: "link",
							Attrs: Attr{
								Href: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
							},
						},
					},
				},
				{
					Type: "text",
					Text: " and there you go ",
				},
				{
					Type: "text",
					Text: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
					Marks: []mark{
						{
							Type: "link",
							Attrs: Attr{
								Href: "https://jenkins.tray.io/job/qa-api-utils-tests/6548/allure/#suites/7deabbaf120515942d030aa4a12b42ab/8caffa3c0bf480ed/",
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, linkifyJiraApiDescriptionObjects(inputObjects), expected)
}
