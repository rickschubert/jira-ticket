package jira

import (
	"encoding/json"
	"fmt"

	"github.com/rickschubert/jira-ticket/constants"
	"github.com/rickschubert/jira-ticket/utils"
)

type assignTicketPostBody struct {
	AccountId string `json:"accountId"`
}

func SelfAssignTicket(ticketInfo NewTicket) {
	settings := constants.GetSettings()
	postBody := assignTicketPostBody{
		AccountId: settings.UserId,
	}

	inputJson, errMarshalling := json.MarshalIndent(postBody, "", "    ")
	utils.HandleErrorStrictly(errMarshalling)

	resp, err := requestClient.R().
		SetBasicAuth(settings.UserName, settings.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(string(inputJson)).
		Put(fmt.Sprintf("%s/rest/api/3/issue/%s/assignee", settings.JiraBaseUrl, ticketInfo.Key))
	utils.HandleErrorStrictly(err)

	if resp.StatusCode() != 204 {
		panic(fmt.Sprintf("We were expecting status code 204 but instead received %v when self-assigning the ticket. The response says:\n\n%v", resp.StatusCode(), resp))
	}
}
