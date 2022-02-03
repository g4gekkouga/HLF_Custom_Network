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

var logger = flogging.MustGetLogger("pmcc")

type Policy struct {
	PolicyID    string            `json:"policyID"`
	SubjectAttr map[string]string `json:"subjectAttr"`
	ObjectAttr  map[string]string `json:"objectAttr"`
	Rules       map[string]string `json:"rules"`
}

func (s *SmartContract) AddPolicy(ctx contractapi.TransactionContextInterface, policyID string, subjectAttr string, objectAttr string, rules string) (string, error) {

	if len(policyID) == 0 {
		return "", fmt.Errorf("Please enter valid Policy ID")
	}

	attrs := strings.Split(subjectAttr, ",")

	var subjAttrsmap = make(map[string]string)

	for i := 0; i < len(attrs); i++ {
		attrPair := strings.Split(attrs[i], ":")
		subjAttrsmap[attrPair[0]] = attrPair[1]
	}

	attrs = strings.Split(objectAttr, ",")

	var objAttrsmap = make(map[string]string)

	for i := 0; i < len(attrs); i++ {
		attrPair := strings.Split(attrs[i], ":")
		objAttrsmap[attrPair[0]] = attrPair[1]
	}

	attrs = strings.Split(rules, ",")

	var rulesmap = make(map[string]string)

	for i := 0; i < len(attrs); i++ {
		attrPair := strings.Split(attrs[i], ":")
		rulesmap[attrPair[0]] = attrPair[1]
	}

	policy := Policy{
		PolicyID:    policyID,
		SubjectAttr: subjAttrsmap,
		ObjectAttr:  objAttrsmap,
		Rules:       rulesmap,
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("CreatePolicy", policyJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(policyID, policyJSON)
}

func (s *SmartContract) GetPolicy(ctx contractapi.TransactionContextInterface, policyID string) (*Policy, error) {
	if len(policyID) == 0 {
		return nil, fmt.Errorf("Please enter valid Policy ID")
	}

	policyJSON, err := ctx.GetStub().GetState(policyID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if policyJSON == nil {
		return nil, fmt.Errorf("%s does not exist", policyID)
	}

	policy := new(Policy)
	_ = json.Unmarshal(policyJSON, policy)

	return policy, nil

}

func (s *SmartContract) GetAllPolicies(ctx contractapi.TransactionContextInterface) []*Policy {

	resultsIterator, err := ctx.GetStub().GetStateByRange("policy", "resource")
	if err != nil {
		return nil
	}
	defer resultsIterator.Close()

	var policies []*Policy
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil
		}

		var policy Policy
		err = json.Unmarshal(queryResponse.Value, &policy)
		if err != nil {
			return nil
		}
		policies = append(policies, &policy)
	}

	return policies
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
