package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("amcc")

type Subject struct {
	SubjectID   string            `json:"subjectID"`
	SubjectName string            `json:"subjectName"`
	Attributes  map[string]string `json:"attributes"`
}

type Object struct {
	ObjectID   string            `json:"objectID"`
	Attributes map[string]string `json:"attributes"`
}

func (s *SmartContract) AddSubject(ctx contractapi.TransactionContextInterface, subjectID string, subjectName string, attributes string) (string, error) {

	if len(subjectID) == 0 {
		return "", fmt.Errorf("Please enter valid Subject ID")
	}

	attrs := strings.Split(attributes, ",")

	var attrsmap = make(map[string]string)

	for i := 0; i < len(attrs); i++ {
		attrPair := strings.Split(attrs[i], ":")
		attrsmap[attrPair[0]] = attrPair[1]
	}

	subject := Subject{
		SubjectID:   subjectID,
		SubjectName: subjectName,
		Attributes:  attrsmap,
	}

	subjectJSON, err := json.Marshal(subject)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("CreateSubject", subjectJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(subjectID, subjectJSON)
}

func (s *SmartContract) GetSubject(ctx contractapi.TransactionContextInterface, subjectID string) *Subject {
	if len(subjectID) == 0 {
		return nil
	}

	subjectJSON, err := ctx.GetStub().GetState(subjectID)

	if err != nil {
		return nil
	}

	if subjectJSON == nil {
		return nil
	}

	subject := new(Subject)
	_ = json.Unmarshal(subjectJSON, subject)

	return subject

}

func (s *SmartContract) GetAllSubjects(ctx contractapi.TransactionContextInterface) ([]*Subject, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("subject", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var subjects []*Subject
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var subject Subject
		err = json.Unmarshal(queryResponse.Value, &subject)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, &subject)
	}

	return subjects, nil
}

func (s *SmartContract) AddObject(ctx contractapi.TransactionContextInterface, objectID string, attributes string) (string, error) {

	if len(objectID) == 0 {
		return "", fmt.Errorf("Please enter valid Object ID")
	}

	attrs := strings.Split(attributes, ",")

	var attrsmap = make(map[string]string)

	for i := 0; i < len(attrs); i++ {
		attrPair := strings.Split(attrs[i], ":")
		attrsmap[attrPair[0]] = attrPair[1]
	}

	object := Object{
		ObjectID:   objectID,
		Attributes: attrsmap,
	}

	objectJSON, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("CreateObject", objectJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(objectID, objectJSON)
}

func (s *SmartContract) GetObject(ctx contractapi.TransactionContextInterface, objectID string) *Object {
	if len(objectID) == 0 {
		return nil
	}

	objectJSON, err := ctx.GetStub().GetState(objectID)

	if err != nil {
		return nil
	}

	if objectJSON == nil {
		return nil
	}

	object := new(Object)
	_ = json.Unmarshal(objectJSON, object)

	return object
}

func (s *SmartContract) GetAllObjects(ctx contractapi.TransactionContextInterface) ([]*Object, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("object", "policy")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var objects []*Object
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var object Object
		err = json.Unmarshal(queryResponse.Value, &object)
		if err != nil {
			return nil, err
		}
		objects = append(objects, &object)
	}

	return objects, nil
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
