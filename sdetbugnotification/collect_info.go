package sdetbugnotification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Songmu/prompter"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/cucumberfeatureparser"
	"github.com/rickschubert/jira-ticket/utils"
)

type knownIssuesWorkflowInputSchema struct {
	Bucket      string `json:"bucket"`
	Feature     string `json:"feature"`
	Step        string `json:"step"`
	Error       string `json:"error"`
	Jira        string `json:"jira"`
	Cause       string `json:"cause,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func promptForStep() string {
	step := prompter.Prompt("Enter the step where the feature fails", "")
	if step == "" {
		log.Fatal("You need to enter a valid step")
	}
	return step
}

func promptForError() string {
	err := prompter.Prompt("Enter the error which is shown as failure", "")
	if err == "" {
		log.Fatal("You need to enter a message for that error")
	}
	return err
}

func promptForEnvironment() string {
	env := prompter.Prompt("Enter the environment where the feature fails (optional)", "")
	return env
}

func promptForCause() string {
	cause := prompter.Prompt("Enter the cause of the error (optional)", "")
	return cause
}

func collectInfo(ticketID string, title string, description string) knownIssuesWorkflowInputSchema {
	fullText := title + " " + description
	mentionedFeature, err := cucumberfeatureparser.GetFeatureInfoOfTextMentioningFeature(fullText)
	utils.HandleErrorStrictly(err)

	color.Yellow(fmt.Sprintf("Please enter some more information to attach to the SDET Known Issues Notification. We managed to retrieve the name of the feature (%s) and the name of the bucket (%s), but we don't know yet some other vital information.", mentionedFeature.Name, mentionedFeature.Bucket))
	step := promptForStep()
	errAppearingWithFeature := promptForError()
	environment := promptForEnvironment()
	cause := promptForCause()

	return knownIssuesWorkflowInputSchema{
		Bucket:      mentionedFeature.Bucket,
		Feature:     mentionedFeature.Name,
		Jira:        ticketID,
		Step:        step,
		Error:       errAppearingWithFeature,
		Environment: environment,
		Cause:       cause,
	}
}

func CreateNotification(ticketID string, title string, description string) {
	bugInfo := collectInfo(ticketID, title, description)
	settings := constants.GetSettings()
	inputJSON, errMarshalling := json.MarshalIndent(bugInfo, "", "    ")
	utils.HandleErrorStrictly(errMarshalling)
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(string(inputJSON)).
		Post(settings.KnownIssueWorkflowUrl)
	utils.HandleErrorStrictly(err)

	if resp.StatusCode() != 200 {
		log.Fatal(fmt.Sprintf("We were expecting status code 200 but instead received %v. The response says:\n\n%v", resp.StatusCode(), resp))
	}

	color.Green("We have sent off a new SDET notification for this, thanks!")
}
