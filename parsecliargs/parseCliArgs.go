package parsecliargs

import (
	"fmt"
	"os"
	"strings"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/utils"
)

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

func GetTicketTitleAndDescription(cliArgsPassed CliArgs, clipboardContent string) (title string, description string) {
	if cliArgsPassed.ParseFromClipboard {
		return parseTitleAndDescriptionFromText(clipboardContent)
	} else {
		title := cliArgsPassed.TicketTitle
		description := cliArgsPassed.TicketDescription
		if cliArgsPassed.ParseFromClipboard {
			description = clipboardContent
		}
		return title, description
	}
}

func verifyThatProjectNameHasBeenPassedInArguments(args []string) {
	if len(args) < 1 {
		utils.ThrowCustomError("You need to pass a valid shortcut that matches a key in your ~/.jiraticketcreator file. If you are unsure what that means, please refer to the readme.")
	}
}

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

type CliArgs struct {
	Project                        constants.Project
	TicketTitle                    string
	TicketDescription              string
	ParseFromClipboard             bool
	CreateKnownSDETBugNotification bool
	SelfAssign                     bool
	TransitioningJiraID            string
	PriorityJiraID                 string
	Labels                         []string
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

var optionArgumentsThatShouldBeFollowedByAString = []string{"--priority", "-p", "--transition", "-t", "--label", "-l"}

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

func shouldSelfAssignTicket(args []string) bool {
	_, foundShortform := utils.Find(args, "--self")
	_, foundLong := utils.Find(args, "--self-assign")
	_, foundShort := utils.Find(args, "-s")
	return foundShortform || foundLong || foundShort
}

func ValidateCommandLineArguments() CliArgs {
	var cliArgumentsPassed = CliArgs{}

	args := os.Args[1:]
	verifyThatProjectNameHasBeenPassedInArguments(args)
	verifyThatRequiredArgumentsHaveBeenPassedWithOptions(args)

	cliArgumentsPassed.Project = getProject(args[0])
	cliArgumentsPassed.Labels = getLabels(args, cliArgumentsPassed.Project)
	cliArgumentsPassed.TicketTitle = getTicketTitle(args)
	cliArgumentsPassed.TicketDescription = getTicketDescription(args)
	cliArgumentsPassed.ParseFromClipboard = shouldUseClipboardContentAsDescription(args)
	cliArgumentsPassed.CreateKnownSDETBugNotification = shouldCreateKnownSDETBugNotification(args)
	cliArgumentsPassed.SelfAssign = shouldSelfAssignTicket(args)
	cliArgumentsPassed.TransitioningJiraID = getTransitionID(args, cliArgumentsPassed.Project)
	cliArgumentsPassed.PriorityJiraID = getPriorityID(args, cliArgumentsPassed.Project)
	return cliArgumentsPassed
}
