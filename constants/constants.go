package constants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/rickschubert/jira-ticket/utils"
)

const SETTINGS_FILE_RELATIVE_PATH = "~/.jiraticketcreator"

type Project struct {
	Shortcut  string   `json:"shortcut"`
	Id        string   `json:"id"`
	IssueType string   `json:"defaultIssueType"`
	Labels    []string `json:"labels"`
}

type Settings struct {
	UserId                string    `json:"JIRA_USER_ID"`
	ApiKey                string    `json:"JIRA_API_KEY"`
	UserName              string    `json:"JIRA_USER_NAME"`
	JiraBaseUrl           string    `json:"JIRA_BASE_URL"`
	FeatureFolderPath     string    `json:"FEATURE_FOLDER,omitempty"`
	Projects              []Project `json:"SHORTCUTS,omitempty"`
	KnownIssueWorklfowUrl string    `json:"KNOWN_ISSUE_WORKFLOW_URL"`
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

func GetSettings() Settings {
	settingsFileContent := readSettingsFileContent()
	var settings Settings
	err := json.Unmarshal(settingsFileContent, &settings)
	utils.HandleErrorStrictly(err)
	return settings
}
