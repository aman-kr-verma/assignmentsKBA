package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type RawMaterialContract struct {
	contractapi.Contract
}

type RawMaterial struct {
	AssetType         string `json:"assetType"`
	MaterialId        string `json:"materialId"`
	MaterialType      string `json:"materialType"`
	Quantity          string `json:"quantity"`
	SupplierId        string `json:"supplierId"`
	DateOfManufacture string `json:"dateOfManufacture"`
	Model             string `json:"model"`
}

const collectionName string = "RawMaterialCollection"

// RawMaterialExists returns true when asset with given ID exists in private data collection
func (r *RawMaterialContract) RawMaterialExists(ctx contractapi.TransactionContextInterface, materialID string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, materialID)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

// CreateRawMaterial creates an intance of raw material, with the provided details
func (r *RawMaterialContract) CreateRawMaterial(ctx contractapi.TransactionContextInterface, materialID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "SupplierMSP" {

		exists, err := r.RawMaterialExists(ctx, materialID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the raw material %s already exists", materialID)
		}

		var material RawMaterial

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of materialType, quantity, supplierId, dateOfManufacture, model")
		}

		materialType, exists := transientData["materialType"]
		if !exists {
			return "", fmt.Errorf("the materialType was not specified in transient data. Please try again")
		}
		material.MaterialType = string(materialType)

		quantity, exists := transientData["quantity"]
		if !exists {
			return "", fmt.Errorf("the quantity was not specified in transient data. Please try again")
		}
		material.Quantity = string(quantity)

		supplierId, exists := transientData["supplierId"]
		if !exists {
			return "", fmt.Errorf("the supplierId was not specified in transient data. Please try again")
		}
		material.SupplierId = string(supplierId)

		dateOfManufacture, exists := transientData["dateOfManufacture"]
		if !exists {
			return "", fmt.Errorf("the dateOfManufacture was not specified in transient data. Please try again")
		}
		material.DateOfManufacture = string(dateOfManufacture)

		model, exists := transientData["model"]
		if !exists {
			return "", fmt.Errorf("the model was not specified in transient data. Please try again")
		}
		material.Model = string(model)

		material.AssetType = "rawmaterial"
		material.MaterialId = string(materialID)

		bytes, _ := json.Marshal(material)
		err = ctx.GetStub().PutPrivateData(collectionName, string(materialID), bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("raw material with id %v added successfully", materialID), nil
	} else {
		return fmt.Sprintf("raw material cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

// ReadRawMaterial retrieves an instance of Raw Material from the private data collection
func (r *RawMaterialContract) ReadRawMaterial(ctx contractapi.TransactionContextInterface, materialID string) (*RawMaterial, error) {
	exists, err := r.RawMaterialExists(ctx, materialID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", materialID)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, materialID)
	if err != nil {
		return nil, fmt.Errorf("could not get the private data. %s", err)
	}
	var material RawMaterial

	err = json.Unmarshal(bytes, &material)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private data collection data to type Raw Material")
	}

	return &material, nil

}

// DeleteRawMaterial deletes an instance of Raw Material from the private data collection
func (r *RawMaterialContract) DeleteRawMaterial(ctx contractapi.TransactionContextInterface, materialID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	if clientOrgID == "SupplierMSP" {

		exists, err := r.RawMaterialExists(ctx, materialID)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset %s does not exist", materialID)
		}

		return ctx.GetStub().DelPrivateData(collectionName, materialID)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the order", clientOrgID)
	}
}

// GetAllRawMaterials fetches all the raw materials
func (r *RawMaterialContract) GetAllRawMaterials(ctx contractapi.TransactionContextInterface) ([]*RawMaterial, error) {
	queryString := `{"selector":{"assetType":"rawmaterial"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return RawMaterialResultIteratorFunction(resultsIterator)
}

// GetRawMaterialsByRange fetches raw materials by range of material ids provided
func (r *RawMaterialContract) GetRawMaterialsByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*RawMaterial, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)

	if err != nil {
		return nil, fmt.Errorf("could not fetch the private data by range. %s", err)
	}
	defer resultsIterator.Close()

	return RawMaterialResultIteratorFunction(resultsIterator)

}

// iterator function

func RawMaterialResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*RawMaterial, error) {
	var materials []*RawMaterial
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of result iterator. %s", err)
		}
		var material RawMaterial
		err = json.Unmarshal(queryResult.Value, &material)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		materials = append(materials, &material)
	}

	return materials, nil
}
