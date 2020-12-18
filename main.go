package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/cucumberfeatureparser"
	"github.com/rickschubert/jira-ticket/jira"
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

func openLinkInBrowser(link string) {
	color.Green(fmt.Sprintf("Your new ticket has been successfully created! The link is %s \nWe will now try to open the new ticket for you in the browser. (Probably doesn't work on Windows.)", link))
	cmd := exec.Command("open", link)
	cmd.Run()
}

type cliArgs struct {
	project                        constants.Project
	ticketTitle                    string
	ticketDescription              string
	parseFromClipboard             bool
	createKnownSDETBugNotification bool
	selfAssign                     bool
	transitioningJiraID            string
	priorityJiraID                 string
	labels                         []string
}

func shouldUseClipboardContentAsDescription(args []string) bool {
	var positionalArguments []string
	var isOptionToAFlag bool
	for _, arg := range args {
		if isOptionToAFlag {
			continue
		}
		_, isAFlagThatRequiresOption := utils.Find(optionArgumentsThatShouldBeFollowedByAString, arg)
		if !isArgumentANonPositionalOptionalArgument(arg) {
			positionalArguments = append(positionalArguments, arg)
		}
		// No need to check here again if all flags come with option as we already do so earlier
		if isAFlagThatRequiresOption {
			isOptionToAFlag = true
		} else {
			isOptionToAFlag = false
		}
	}
	return len(positionalArguments) < 2
}

func shouldCreateKnownSDETBugNotification(args []string) bool {
	_, found := utils.Find(args, "--sdet-bot")
	return found
}

func getTransitionID(args []string, project constants.Project) string {
	idxLong, foundLong := utils.Find(args, "--transition")
	idxShort, foundShort := utils.Find(args, "-t")
	var transitionID string
	if foundLong || foundShort {
		var transitionTitlePassedInCLIArguments string
		if foundLong {
			transitionTitlePassedInCLIArguments = args[idxLong+1]
		} else {
			transitionTitlePassedInCLIArguments = args[idxShort+1]
		}

		id, found := project.Transitions[transitionTitlePassedInCLIArguments]
		if found {
			transitionID = id
		} else {
			utils.ThrowCustomError(fmt.Sprintf("You specified in the command line arguments that you want to transition the ticket using the ID mapped to the key '%s' in your settings file for project '%s' - but such a transition title doesn't exist in your settings.", transitionTitlePassedInCLIArguments, project.Shortcut))
		}
	}
	return transitionID
}

func getProject(desiredProject string) constants.Project {
	settings := constants.GetSettings()
	desiredProjectKeyLowerCased := strings.ToLower(desiredProject)

	var projectToReturn constants.Project
	for _, project := range settings.Projects {
		if strings.ToLower(project.Shortcut) == desiredProjectKeyLowerCased {
			projectToReturn = project
			break
		}
	}
	return projectToReturn
}

func getLabels(args []string, project constants.Project) []string {
	var labelsPassedInCLIArguments []string
	for idx, arg := range args {
		// This works on the assumption that after `--label` or `-l`, a label
		// is passed. If this is not the case, it will continue but with errors.
		if arg == "--label" || arg == "-l" {
			labelsPassedInCLIArguments = append(labelsPassedInCLIArguments, args[idx+1])
		}
	}
	return append(project.Labels, labelsPassedInCLIArguments...)
}

func mapPriorityWordToPriorityID(priorityPassed string) string {
	prioritiesMap := make(map[string]string)
	prioritiesMap["critical"] = "1"
	prioritiesMap["high"] = "2"
	prioritiesMap["medium"] = "3"
	prioritiesMap["low"] = "4"
	prioritiesMap["lowest"] = "5"
	id, found := prioritiesMap[strings.ToLower(priorityPassed)]
	if found {
		return id
	}
	return ""
}

func getPriorityID(args []string, proj constants.Project) string {
	var priorityID string
	priorityID = mapPriorityWordToPriorityID(proj.Priority)
	for idx, arg := range args {
		// This works on the assumption that after `--priority` or `-p`, a priority ID
		// is passed, i.e. 1 (Critical) or 5 (Lowest)
		if arg == "--priority" || arg == "-p" {
			nextArg := args[idx+1]
			priorityIDDerivedFromWord := mapPriorityWordToPriorityID(nextArg)
			if priorityIDDerivedFromWord != "" {
				priorityID = priorityIDDerivedFromWord
			} else {
				priorityID = nextArg
			}
		}
	}
	return priorityID
}

func isArgumentANonPositionalOptionalArgument(arg string) bool {
	if len(arg) == 2 {
		firstCharacterOfArgument := arg[0:1]
		return firstCharacterOfArgument == "-"
	} else if len(arg) > 2 {
		firstTwoCharactersOfArgument := arg[0:2]
		return firstTwoCharactersOfArgument == "--"
	} else {
		return false
	}
}

func getTicketTitle(args []string) string {
	if len(args) < 2 {
		return ""
	}
	if isArgumentANonPositionalOptionalArgument(args[1]) {
		return ""
	} else {
		return args[1]
	}
}

func getTicketDescription(args []string) string {
	if len(args) < 3 {
		return ""
	}
	if isArgumentANonPositionalOptionalArgument(args[2]) {
		return ""
	} else {
		return args[2]
	}
}

func shouldSelfAssignTicket(args []string) bool {
	_, foundShortform := utils.Find(args, "--self")
	_, foundLong := utils.Find(args, "--self-assign")
	_, foundShort := utils.Find(args, "-s")
	return foundShortform || foundLong || foundShort
}

var optionArgumentsThatShouldBeFollowedByAString = []string{"--priority", "-p", "--transition", "-t", "--label", "-l"}

func verifyThatRequiredArgumentsHaveBeenPassedWithOptions(args []string) {
	for idx, arg := range args {
		_, isAFlagThatRequiresOption := utils.Find(optionArgumentsThatShouldBeFollowedByAString, arg)
		errorMsg := fmt.Sprintf("You passed the argument \"%s\" but didn't follow it with an option. Please provide a value along with it.", arg)
		if isAFlagThatRequiresOption {
			var nextArg string
			if len(args) == idx+1 {
				utils.ThrowCustomError(errorMsg)
			} else {
				nextArg = args[idx+1]
			}
			if isArgumentANonPositionalOptionalArgument(nextArg) {
				utils.ThrowCustomError(errorMsg)
			}
		}
	}
}

func verifyThatProjectNameHasBeenPassedInArguments(args []string) {
	if len(args) < 1 {
		utils.ThrowCustomError("You need to pass a valid shortcut that matches a key in your ~/.jiraticketcreator file. If you are unsure what that means, please refer to the readme.")
	}
}

func validateCommandLineArguments() cliArgs {
	var cliArgumentsPassed = cliArgs{}

	args := os.Args[1:]
	verifyThatProjectNameHasBeenPassedInArguments(args)
	verifyThatRequiredArgumentsHaveBeenPassedWithOptions(args)

	cliArgumentsPassed.project = getProject(args[0])
	cliArgumentsPassed.labels = getLabels(args, cliArgumentsPassed.project)
	cliArgumentsPassed.ticketTitle = getTicketTitle(args)
	cliArgumentsPassed.ticketDescription = getTicketDescription(args)
	cliArgumentsPassed.parseFromClipboard = shouldUseClipboardContentAsDescription(args)
	cliArgumentsPassed.createKnownSDETBugNotification = shouldCreateKnownSDETBugNotification(args)
	cliArgumentsPassed.selfAssign = shouldSelfAssignTicket(args)
	cliArgumentsPassed.transitioningJiraID = getTransitionID(args, cliArgumentsPassed.project)
	cliArgumentsPassed.priorityJiraID = getPriorityID(args, cliArgumentsPassed.project)
	return cliArgumentsPassed
}

func getClipboardContent() string {
	content, err := clipboard.ReadAll()
	if err != nil {
		return ""
	} else {
		return content
	}
}

func parseTitleAndDescriptionFromText(text string) (title string, description string) {
	lineBreak := "\n"
	splitAtLineBreak := strings.SplitAfterN(text, lineBreak, 2)
	if len(splitAtLineBreak) > 1 {
		title := strings.TrimRight(splitAtLineBreak[0], lineBreak)
		description := strings.TrimLeft(splitAtLineBreak[1], lineBreak)
		return title, description
	} else {
		return text, ""
	}
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

func collectInformationToCreateKnownSdetBugNotification(ticketID string, title string, description string) knownIssuesWorkflowInputSchema {
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

func createKnownSdetBugNotification(bugInfo knownIssuesWorkflowInputSchema) {
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

func getNewTicketInput(cliArgsPassed cliArgs, clipboardContent string) jira.CreateNewticketInput {
	var newTicketInfo jira.CreateNewticketInput
	newTicketInfo.Labels = cliArgsPassed.labels
	newTicketInfo.ProjectId = cliArgsPassed.project.Id
	newTicketInfo.IssueType = cliArgsPassed.project.IssueType
	newTicketInfo.AssigneeUserId = cliArgsPassed.project.Assignee
	newTicketInfo.PriorityId = cliArgsPassed.priorityJiraID

	title, description := getTicketTitleAndDescription(cliArgsPassed, clipboardContent)
	newTicketInfo.Title = title
	newTicketInfo.Description = description

	return newTicketInfo
}

func main() {
	cliArgsPassed := validateCommandLineArguments()
	clipboardContent := getClipboardContent()
	if cliArgsPassed.parseFromClipboard {
		showClipboardAndAskIfOkay(clipboardContent)
	}

	newTicketInput := getNewTicketInput(cliArgsPassed, clipboardContent)
	if cliArgsPassed.selfAssign {
		newTicketInput.AssigneeUserId = constants.GetSettings().UserId
	}
	if cliArgsPassed.transitioningJiraID != "" {
		newTicketInput.TransitionId = cliArgsPassed.transitioningJiraID
	}
	if cliArgsPassed.priorityJiraID != "" {
		newTicketInput.PriorityId = cliArgsPassed.priorityJiraID
	}
	ticketInfo := jira.CreateNewTicket(newTicketInput)

	openLinkInBrowser(ticketInfo.Link)

	if cliArgsPassed.createKnownSDETBugNotification {
		bugInfo := collectInformationToCreateKnownSdetBugNotification(ticketInfo.Key, newTicketInput.Title, newTicketInput.Description)
		createKnownSdetBugNotification(bugInfo)
	}
}

func showClipboardAndAskIfOkay(clipboardContent string) {
	promptMsg := fmt.Sprintf("The following text is in your clipboard. Are you sure you want to attach this to the ticket or parse a ticket out of it?\n\n%s\n\n", clipboardContent)
	if !prompter.YN(promptMsg, false) {
		log.Fatal("Exiting program.")
	}
}

func getTicketTitleAndDescription(cliArgsPassed cliArgs, clipboardContent string) (title string, description string) {
	if cliArgsPassed.parseFromClipboard {
		return parseTitleAndDescriptionFromText(clipboardContent)
	} else {
		title := cliArgsPassed.ticketTitle
		description := cliArgsPassed.ticketDescription
		if cliArgsPassed.parseFromClipboard {
			description = clipboardContent
		}
		return title, description
	}
}
