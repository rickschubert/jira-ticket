package utils

import (
	"fmt"
	"log"
)

func HandleErrorStrictly(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleErrorStrictlyWithMessage(err error, msg string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s\n\n%s", msg, err))
	}
}
