package parsecliargs

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
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             false,
		TicketTitle:                    "this will be the title",
		CreateKnownSDETBugNotification: false,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded", "this will be the title"}

	actual := ValidateCommandLineArguments()

	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsRegularInvoking() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             false,
		TicketTitle:                    "this will be the title",
		CreateKnownSDETBugNotification: false,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded", "this will be the title"}

	actual := ValidateCommandLineArguments()

	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldCreateKnownSDETNotification() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldCreateKnownSDETNotificationAndSelfAssign() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot", "--self"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldCreateKnownSDETNotificationAndSelfAssignLong() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot", "--self-assign"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldCreateKnownSDETNotificationAndSelfAssignShort() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     true,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot", "-s"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboard() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: false,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboardAndCreateKnownSDETNotification() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestArgsShouldParseFullTicketFromClipboardAndCreateKnownSDETNotificationReversed() {
	expected := CliArgs{
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		ParseFromClipboard:             true,
		TicketTitle:                    "",
		CreateKnownSDETBugNotification: true,
		SelfAssign:                     false,
	}
	os.Args = []string{"jira-ticket", "embedded", "--sdet-bot"}
	actual := ValidateCommandLineArguments()
	assert.Equal(suite.T(), expected, actual)
}

func (suite *testSuite) TestTicketTitleAndDescriptionRetrievalWithoutTakingClipboardContentIntoDescription() {
	cliArgumentsRetrieved := CliArgs{
		TicketTitle:        "This should be the ticket title!!!",
		ParseFromClipboard: false,
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		CreateKnownSDETBugNotification: false,
		SelfAssign:                     false,
	}
	clipboardContent := ""
	title, description := GetTicketTitleAndDescription(cliArgumentsRetrieved, clipboardContent)
	assert.Equal(suite.T(), "This should be the ticket title!!!", title)
	assert.Equal(suite.T(), "", description)
}

func (suite *testSuite) TestParsingTicketTitleAndDescriptionFromClipboard() {
	cliArgumentsRetrieved := CliArgs{
		TicketTitle:        "",
		ParseFromClipboard: true,
		Project: constants.Project{
			Shortcut:  "embedded",
			Id:        "10059",
			IssueType: "10004",
			Labels:    []string(nil),
		},
		CreateKnownSDETBugNotification: false,
		SelfAssign:                     false,
	}
	clipboardContent := "This will be the title\n\nThis will be the description"
	title, description := GetTicketTitleAndDescription(cliArgumentsRetrieved, clipboardContent)
	assert.Equal(suite.T(), "This will be the title", title)
	assert.Equal(suite.T(), "This will be the description", description)
}
