package constants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/rickschubert/jira-ticket/utils"
)

const SETTINGS_FILE_RELATIVE_PATH = "~/.jiraticketcreator"

type Project struct {
	Shortcut  string   `json:"shortcut"`
	Id        string   `json:"id"`
	IssueType string   `json:"defaultIssueType"`
	Labels    []string `json:"labels"`
	Assignee  string   `json:"assignee"`
	// Object key: transitionTitle | Object value: Jira transition ID
	Transitions map[string]string `json:"transitions"`
	Priority    string            `json:"priority"`
}

type Settings struct {
	UserId                string    `json:"JIRA_USER_ID"`
	ApiKey                string    `json:"JIRA_API_KEY"`
	UserName              string    `json:"JIRA_USER_NAME"`
	JiraBaseUrl           string    `json:"JIRA_BASE_URL"`
	FeatureFolderPath     string    `json:"FEATURE_FOLDER,omitempty"`
	Projects              []Project `json:"SHORTCUTS,omitempty"`
	KnownIssueWorkflowUrl string    `json:"KNOWN_ISSUE_WORKFLOW_URL"`
}

func getAbsolutePathToSettingsFile() string {
	absolutePath, err := homedir.Expand(SETTINGS_FILE_RELATIVE_PATH)
	utils.HandleErrorStrictly(err)
	return absolutePath
}

func readSettingsFileContent() []byte {
	settingsFilePath := getAbsolutePathToSettingsFile()
	settingsFile, err := ioutil.ReadFile(settingsFilePath)
	utils.HandleErrorStrictlyWithMessage(err, fmt.Sprintf("We have been unable to read the credentials file. This should be located at %s", settingsFilePath))
	return settingsFile
}

func getMockSettings() Settings {
	mockTransitions := make(map[string]string)
	mockTransitions["p"] = "21"
	mockTransitions["prog"] = "21"
	mockTransitions["progress"] = "21"
	mockTransitions["inprogress"] = "21"
	dir := os.Getenv("ROOT_DIRECTORY")
	featureFixturesPath := path.Join(dir, "test_fixtures/features")

	return Settings{
		ApiKey:                "mockAPIKey",
		FeatureFolderPath:     featureFixturesPath,
		KnownIssueWorkflowUrl: "https://mockUrl.com",
		JiraBaseUrl:           "https://mycompany.atlassian.net",
		UserId:                "1",
		UserName:              "MyUserName",
		Projects: []Project{
			{
				Assignee:  "5d19daa472f6850cd226fe1d",
				Id:        "10131",
				IssueType: "10004",
				Labels:    []string{"Frontend"},
				Shortcut:  "builderfe",
			},
			{
				Id:          "10064",
				IssueType:   "10084",
				Transitions: mockTransitions,
				Shortcut:    "qa",
			},
			{
				Shortcut:  "embedded",
				Id:        "10059",
				IssueType: "10004",
			},
		},
	}
}

func GetSettings() Settings {
	if os.Getenv("UNIT_TESTS") == "true" {
		return getMockSettings()
	}
	settingsFileContent := readSettingsFileContent()
	var settings Settings
	err := json.Unmarshal(settingsFileContent, &settings)
	utils.HandleErrorStrictly(err)
	return settings
}
