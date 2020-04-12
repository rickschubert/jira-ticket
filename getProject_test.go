package main

import (
	"testing"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectSuccess(t *testing.T) {
	assert.Equal(t, getProject("embedded"), constants.Project{
		Shortcut:  "embedded",
		Id:        "10059",
		IssueType: "10004",
		Labels:    []string(nil),
	})
}
