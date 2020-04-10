package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Constants struct {
	RickUserId string
}


func main() {
	consts := Constants{
		RickUserId: "5b9f82a2f226b370980f271a",
	}

	fmt.Println("Hello22")

	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetBasicAuth("rick@tray.io", "JeQzZGNxkRUMR2NUeR7602A0").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(fmt.Sprintf(`{
  "update": {},
  "fields": {
    "summary": "Main order flow broken",
    "issuetype": {
      "id": "10271"
    },
    "project": {
      "id": "10149"
    },
    "description": {
      "type": "doc",
      "version": 1,
      "content": [
        {
          "type": "paragraph",
          "content": [
            {
              "text": "Order entry fails when selecting supplier.",
              "type": "text"
            }
          ]
        }
      ]
    },
    "reporter": {
      "id": "%s"
    },
    "assignee": {
      "id": "%s"
    }
  }
}`, consts.RickUserId, consts.RickUserId)).
		Post("https://trayio.atlassian.net/rest/api/3/issue")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Body()))
}
