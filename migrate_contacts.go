package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/k-capehart/go-salesforce"
)

type Contact struct {
	Id        string `soql:"selectColumn,fieldName=Id" json:"Id"`
	AccountId string `soql:"selectColumn,fieldName=AccountId" json:"AccountId"`
}

type ContactQueryCriteria struct {
	AccountId string `soql:"equalsOperator,fieldName=AccountId"`
}

type ContactSoqlQuery struct {
	SelectClause Contact              `soql:"selectClause,tableName=Contact"`
	WhereClause  ContactQueryCriteria `soql:"whereClause"`
}

func main() {
	args := os.Args
	if len(args) < 3 {
		panic(errors.New("expected 2 command line arguments in addition to program"))
	}
	srcAccount := os.Args[1]
	targetAccount := os.Args[2]

	sf, err := salesforce.Init(salesforce.Creds{
		Domain:         {YOUR SF DOMAIN},
		ConsumerKey:    {YOUR CONNECTED APP CONSUMER KEY},
		ConsumerSecret: {YOUR CONNECTED APP CONSUMER SECRET},
	})
	if err != nil {
		panic(err)
	}

	contacts := []Contact{}
	contactSoqlQuery := ContactSoqlQuery{
		SelectClause: Contact{},
		WhereClause: ContactQueryCriteria{
			AccountId: srcAccount,
		},
	}
	err = sf.QueryStruct(contactSoqlQuery, &contacts)
	if err != nil {
		panic(err)
	}

	for i := range contacts {
		contacts[i].AccountId = targetAccount
	}
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	err = sf.UpdateCollection("Contact", contacts, 200)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Print("successfully updated " + strconv.Itoa(len(contacts)) + " contacts")
}
