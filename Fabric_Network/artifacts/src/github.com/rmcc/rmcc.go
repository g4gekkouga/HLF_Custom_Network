package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("rmcc")

type Resource struct {
	ObjectID string `json:"objectID"`
	Data     string `json:"data"`
}

func (s *SmartContract) AddResource(ctx contractapi.TransactionContextInterface, objectID string, data string) (string, error) {

	if len(objectID) == 0 {
		return "", fmt.Errorf("Please enter valid Object ID")
	}

	resource := Resource{
		ObjectID: objectID,
		Data:     data,
	}

	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("CreateResource", resourceJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState("resource_"+objectID, resourceJSON)
}

func (s *SmartContract) GetResource(ctx contractapi.TransactionContextInterface, objectID string) *Resource {
	if len(objectID) == 0 {
		return nil
	}

	resourceJSON, err := ctx.GetStub().GetState("resource_" + objectID)

	if err != nil {
		return nil
	}

	if resourceJSON == nil {
		return nil
	}

	resource := new(Resource)
	_ = json.Unmarshal(resourceJSON, resource)

	return resource

}

func (s *SmartContract) GetAllResources(ctx contractapi.TransactionContextInterface) ([]*Resource, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("resource", "subject")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var resources []*Resource
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var resource Resource
		err = json.Unmarshal(queryResponse.Value, &resource)
		if err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}

	return resources, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating resource management chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting resource management chain code: %s", err.Error())
	}

}
