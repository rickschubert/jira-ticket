package main

import (
	"testing"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectSuccess(t *testing.T) {
	assert.Equal(t, getProject("builder"), constants.Project{
		Shortcut:  "builder",
		Id:        "10131",
		IssueType: "10004",
		Labels:    []string{"Frontend"},
		Assignee:  "5d19daa472f6850cd226fe1d",
	})
}
