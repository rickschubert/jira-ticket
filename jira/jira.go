package jira

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/utils"
)

type CreateNewticketInput struct {
	Title       string
	Description string
	ProjectId   string
	IssueType   string
	Labels      []string
}

type update struct{}

type issuetype struct {
	Id string `json:"id"`
}

type project struct {
	Id string `json:"id"`
}

type reporter struct {
	Id string `json:"id"`
}

type paragraphContent struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type"`
}

type content struct {
	Type    string             `json:"type"`
	Content []paragraphContent `json:"content"`
}

type ticketDescription struct {
	Type    string    `json:"type"`
	Version int       `json:"version"`
	Content []content `json:"content"`
}

type fields struct {
	Summary     string            `json:"summary"`
	IssueType   issuetype         `json:"issuetype"`
	Project     project           `json:"project"`
	Reporter    reporter          `json:"reporter"`
	Description ticketDescription `json:"description"`
	Labels      []string          `json:"labels,omitempty"`
}

type createNewTicketJiraAPIInput struct {
	Update update `json:"update"`
	Fields fields `json:"fields"`
}

type NewTicket struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	ApiLink string `json:"self"`
	Link    string
}

type credentials struct {
	UserId   string `json:"JIRA_USER_ID"`
	ApiKey   string `json:"JIRA_API_KEY"`
	UserName string `json:"JIRA_USER_NAME"`
}

var requestClient = resty.New()

func lineIsNotEmpty(line string) bool {
	return strings.TrimSpace(line) != ""
}

func splitPlainTextDescriptionIntoJiraApiObjects(plainTextDescription string) []content {
	individualLines := strings.Split(plainTextDescription, "\n")
	var apiObjectsPerLine []paragraphContent
	for _, line := range individualLines {
		if lineIsNotEmpty(line) {
			apiObjectsPerLine = append(apiObjectsPerLine, paragraphContent{
				Type: "text",
				Text: line,
			})
		}
		apiObjectsPerLine = append(apiObjectsPerLine, paragraphContent{
			Type: "hardBreak",
		})
	}
	return []content{
		{
			Type:    "paragraph",
			Content: apiObjectsPerLine,
		},
	}
}

func convertDescription(plainTextDescription string) ticketDescription {
	splitUpParagraphs := splitPlainTextDescriptionIntoJiraApiObjects(plainTextDescription)
	return ticketDescription{
		Type:    "doc",
		Version: 1,
		Content: splitUpParagraphs,
	}
}

func CreateNewTicket(input CreateNewticketInput) NewTicket {
	settings := constants.GetSettings()
	descriptionChunks := convertDescription(input.Description)
	newTicketInput := createNewTicketJiraAPIInput{
		Update: update{},
		Fields: fields{
			Reporter:    reporter{Id: settings.UserId},
			Summary:     input.Title,
			IssueType:   issuetype{Id: input.IssueType},
			Project:     project{Id: input.ProjectId},
			Description: descriptionChunks,
		},
	}
	newTicketInput.Fields.Labels = input.Labels

	inputJson, errMarshalling := json.MarshalIndent(newTicketInput, "", "    ")
	fmt.Println(string(inputJson))
	utils.HandleErrorStrictly(errMarshalling)

	resp, err := requestClient.R().
		SetBasicAuth(settings.UserName, settings.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(string(inputJson)).
		Post(fmt.Sprintf("%s/rest/api/3/issue", settings.JiraBaseUrl))
	utils.HandleErrorStrictly(err)

	if resp.StatusCode() != 201 {
		panic(fmt.Sprintf("We were expecting status code 201 but instead received %v. The response says:\n\n%v", resp.StatusCode(), resp))
	}

	var ticket NewTicket
	unmarshallingError := json.Unmarshal(resp.Body(), &ticket)
	utils.HandleErrorStrictly(unmarshallingError)
	ticket.Link = fmt.Sprintf("%s/browse/%s", settings.JiraBaseUrl, ticket.Key)
	return ticket
}
