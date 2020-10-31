package main

import (
	"testing"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectSuccess(t *testing.T) {
	assert.Equal(t, constants.Project{
		Shortcut:  "builderfe",
		Id:        "10131",
		IssueType: "10004",
		Labels:    []string{"Frontend"},
		Assignee:  "5d19daa472f6850cd226fe1d",
	}, getProject("builderfe"))
}

func TestGetProjectWithTransitionsSuccess(t *testing.T) {
	expectedMap := make(map[string]string)
	expectedMap["inprogress"] = "21"
	assert.Equal(t, constants.Project{
		Shortcut:    "qa",
		Id:          "10064",
		IssueType:   "10084",
		Transitions: expectedMap,
	}, getProject("qa"))
}
