package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("adcc")

type Subject struct {
	SubjectID   string            `json:"subjectID"`
	SubjectName string            `json:"subjectName"`
	Attributes  map[string]string `json:"attributes"`
}

type Object struct {
	ObjectID   string            `json:"objectID"`
	Attributes map[string]string `json:"attributes"`
}

type Policy struct {
	PolicyID    string            `json:"policyID"`
	SubjectAttr map[string]string `json:"subjectAttr"`
	ObjectAttr  map[string]string `json:"objectAttr"`
	Rules       map[string]string `json:"rules"`
}

type Resource struct {
	ObjectID string `json:"subjectID"`
	Data     string `json:"data"`
}

func (s *SmartContract) CheckAccess(ctx contractapi.TransactionContextInterface, subject *Subject, object *Object, op string, policies []*Policy) (bool, error) {
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30
	num := rand.Intn(max-min+1) + min
	if num%2 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *SmartContract) GetResource(ctx contractapi.TransactionContextInterface, subjectID string, objectID string, operation string) (string, error) {

	if len(subjectID) == 0 {
		return "", fmt.Errorf("Please enter valid Subject ID")
	}

	if len(objectID) == 0 {
		return "", fmt.Errorf("Please enter valid Object ID")
	}

	params := []string{"GetSubject", subjectID}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("amcc", queryArgs, "mychannel")

	fmt.Printf("response: %s", string(response.Payload))

	var subject *Subject
	_ = json.Unmarshal(response.Payload, &subject)

	fmt.Printf("subject ID : %s", (*subject).SubjectID)

	params = []string{"GetObject", objectID}
	queryArgs = make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response = ctx.GetStub().InvokeChaincode("amcc", queryArgs, "mychannel")

	var object *Object
	_ = json.Unmarshal(response.Payload, &object)

	fmt.Printf("object ID : %s", (*object).ObjectID)

	params = []string{"GetAllPolicies"}
	queryArgs = make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response = ctx.GetStub().InvokeChaincode("pmcc", queryArgs, "mychannel")

	var policies []*Policy
	_ = json.Unmarshal(response.Payload, &response)

	allow, err := s.CheckAccess(ctx, subject, object, operation, policies)

	if err != nil {
		return "", fmt.Errorf("Error while verifying access")
	}

	if allow {

		params := []string{"GetResource", objectID}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}

		response := ctx.GetStub().InvokeChaincode("rmcc", queryArgs, "mychannel")

		return string(response.Payload), nil
	} else {
		return "No Access", fmt.Errorf("Access denied for the object: %s for subject: %s", objectID, subjectID)
	}
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating asset management chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting asset management chain code: %s", err.Error())
	}

}
