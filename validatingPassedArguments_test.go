package main

import (
	"os"
	"testing"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	originalOSArguments []string
}

func (suite *testSuite) SetupTest() {
	suite.originalOSArguments = os.Args
}

func (suite *testSuite) TearDownTestSuite() {
	os.Args = suite.originalOSArguments
}

func TestSuiteValidatingArguments(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestArgsRegularInvokingAndAttachingClipboardContentToDescription() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             false,
		ticketTitle:                    "this will be the title",
		createKnownSDETBugNotification: false,
	}
	os.Args = []string{"jira-ticket", "embedded", "this will be the title"}

	actual := validateCommandLineArguments()

	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsRegularInvoking() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             false,
		ticketTitle:                    "this will be the title",
		createKnownSDETBugNotification: false,
	}
	os.Args = []string{"jira-ticket", "embedded", "this will be the title"}

	actual := validateCommandLineArguments()

	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldCreateKnownSDETNotification() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             true,
		ticketTitle:                    "",
		createKnownSDETBugNotification: true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := validateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboard() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             true,
		ticketTitle:                    "",
		createKnownSDETBugNotification: false,
	}
	os.Args = []string{"jira-ticket", "embedded"}
	actual := validateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboardAndCreateKnownSDETNotification() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             true,
		ticketTitle:                    "",
		createKnownSDETBugNotification: true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := validateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboardAndCreateKnownSDETNotificationReversed() {
	expected := cliArgs{
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		parseFromClipboard:             true,
		ticketTitle:                    "",
		createKnownSDETBugNotification: true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := validateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestTicketTitleAndDescriptionRetrievalWithoutTakingClipboardContentIntoDescription() {
	cliArgumentsRetrieved := cliArgs{
		ticketTitle:        "This should be the ticket title!!!",
		parseFromClipboard: false,
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		createKnownSDETBugNotification: false,
	}
	clipboardContent := ""
	title, description := getTicketTitleAndDescription(cliArgumentsRetrieved, clipboardContent)
	assert.Equal(suite.T(), "This should be the ticket title!!!", title)
	assert.Equal(suite.T(), "", description)
}

func (suite *testSuite) TestParsingTicketTitleAndDescriptionFromClipboard() {
	cliArgumentsRetrieved := cliArgs{
		ticketTitle:        "",
		parseFromClipboard: true,
		project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		createKnownSDETBugNotification: false,
	}
	clipboardContent := "This will be the title\n\nThis will be the description"
	title, description := getTicketTitleAndDescription(cliArgumentsRetrieved, clipboardContent)
	assert.Equal(suite.T(), "This will be the title", title)
	assert.Equal(suite.T(), "This will be the description", description)
}
