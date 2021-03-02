package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/Songmu/prompter"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/jira"
	"github.com/rickschubert/jira-ticket/parsecliargs"
	"github.com/rickschubert/jira-ticket/sdetbugnotification"
)

func openLinkInBrowser(link string) {
	color.Green(fmt.Sprintf("Your new ticket has been successfully created! The link is %s \nWe will now try to open the new ticket for you in the browser. (Probably doesn't work on Windows.)", link))
	cmd := exec.Command("open", link)
	cmd.Run()
}

func getClipboardContent() string {
	content, err := clipboard.ReadAll()
	if err != nil {
		return ""
	} else {
		return content
	}
}

func getNewTicketInput(cliArgsPassed parsecliargs.CliArgs, clipboardContent string) jira.CreateNewticketInput {
	var newTicketInfo jira.CreateNewticketInput
	newTicketInfo.Labels = cliArgsPassed.Labels
	newTicketInfo.ProjectId = cliArgsPassed.Project.Id
	newTicketInfo.IssueType = cliArgsPassed.Project.IssueType
	newTicketInfo.AssigneeUserId = cliArgsPassed.Project.Assignee
	newTicketInfo.PriorityId = cliArgsPassed.PriorityJiraID

	title, description := parsecliargs.GetTicketTitleAndDescription(cliArgsPassed, clipboardContent)
	newTicketInfo.Title = title
	newTicketInfo.Description = description

	return newTicketInfo
}

func main() {
	cliArgsPassed := parsecliargs.ValidateCommandLineArguments()
	clipboardContent := getClipboardContent()
	if cliArgsPassed.ParseFromClipboard {
		showClipboardAndAskIfOkay(clipboardContent)
	}

	newTicketInput := getNewTicketInput(cliArgsPassed, clipboardContent)
	if cliArgsPassed.SelfAssign {
		newTicketInput.AssigneeUserId = constants.GetSettings().UserId
	}
	if cliArgsPassed.TransitioningJiraID != "" {
		newTicketInput.TransitionId = cliArgsPassed.TransitioningJiraID
	}
	if cliArgsPassed.PriorityJiraID != "" {
		newTicketInput.PriorityId = cliArgsPassed.PriorityJiraID
	}
	ticketInfo := jira.CreateNewTicket(newTicketInput)

	openLinkInBrowser(ticketInfo.Link)

	if cliArgsPassed.CreateKnownSDETBugNotification {
		sdetbugnotification.CreateNotification(ticketInfo.Key, newTicketInput.Title, newTicketInput.Description)
	}
}

func showClipboardAndAskIfOkay(clipboardContent string) {
	promptMsg := fmt.Sprintf("The following text is in your clipboard. Are you sure you want to attach this to the ticket or parse a ticket out of it?\n\n%s\n\n", clipboardContent)
	if !prompter.YN(promptMsg, false) {
		log.Fatal("Exiting program.")
	}
}
