package HLF_ABAC

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("reference_chaincode")

type Resource struct {
	RID  string `json:"resourceID"`
	Data string `json:"data"`
}

func AddResource(ctx contractapi.TransactionContextInterface, resourceID string, data string) (string, error) {

	if len(resourceID) == 0 {
		return "", fmt.Errorf("Please enter valid Resource ID")
	}

	resource := Resource{
		RID:  resourceID,
		Data: data,
	}

	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("AddResource", resourceJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(resourceID, resourceJSON)
}

func UpdateResource(ctx contractapi.TransactionContextInterface, resourceID string, data string) (string, error) {

	if len(resourceID) == 0 {
		return "", fmt.Errorf("Please enter valid Resource ID")
	}

	resourceJSON, err := ctx.GetStub().GetState(resourceID)

	if resourceJSON == nil {
		return "", fmt.Errorf("No Resource with given ID")
	}

	resource := Resource{
		RID:  resourceID,
		Data: data,
	}

	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}

	ctx.GetStub().SetEvent("UpdateResource", resourceJSON)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(resourceID, resourceJSON)
}

func GetResource(ctx contractapi.TransactionContextInterface, resourceID string) *Resource {
	if len(resourceID) == 0 {
		return nil
	}

	resourceJSON, err := ctx.GetStub().GetState(resourceID)
	ctx.GetStub().SetEvent("GetResource", resourceJSON)

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
