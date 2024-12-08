package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AccidentRecordContract struct {
	contractapi.Contract
}

type AccidentRecord struct {
	AssetType                       string `json:"assetType"`
	AccidentID                      string `json:"accidentID"`
	Date                            string `json:"date"`
	Location                        string `json:"location"`
	VehiclesNum                     string `json:"vehicleNum"`
	InjuryDetails                   string `json:"injuryDetails"`
	PolicyIDUnderClaimConsideration string `json:"policyIDUnderClaimConsideration"`
	ClaimStatus                     string `json:"claimStatus"`
}

const collectionName = "AccidentRecordCollection"

func (a *AccidentRecordContract) AccidentRecordExists(ctx contractapi.TransactionContextInterface, accidentID string) (bool, error) {

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, accidentID)

	if err != nil {
		return false, fmt.Errorf("could not fetch the private data hash. %s", err)
	}

	return data != nil, nil
}

func (a *AccidentRecordContract) CreateAccidentRecord(ctx contractapi.TransactionContextInterface, accidentID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("could not fetch client identity. %s", err)
	}

	if clientOrgID == "police-insuranceClaimPostAccident-com" {
		exists, err := a.AccidentRecordExists(ctx, accidentID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset (accident record) %s already exists", accidentID)
		}

		var accidentRecord AccidentRecord

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("could not fetch transient data. %s", err)
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of date, location, vehicleNum, injuryDetails, policyIDUnderClaimConsideration,claimStatus")
		}

		date, exists := transientData["date"]
		if !exists {
			return "", fmt.Errorf("the date was not specified in transient data. Please try again")
		}
		accidentRecord.Date = string(date)

		location, exists := transientData["location"]
		if !exists {
			return "", fmt.Errorf("the location was not specified in transient data. Please try again")
		}
		accidentRecord.Location = string(location)

		vehicleNum, exists := transientData["vehicleNum"]
		if !exists {
			return "", fmt.Errorf("the vehicle number was not specified in transient data. Please try again")
		}
		accidentRecord.VehiclesNum = string(vehicleNum)

		injuryDetails, exists := transientData["injuryDetails"]
		if !exists {
			return "", fmt.Errorf("the injuryDetails was not specified in transient data. Please try again")
		}
		accidentRecord.InjuryDetails = string(injuryDetails)

		policyIDUnderClaimConsideration, exists := transientData["policyIDUnderClaimConsideration"]
		if !exists {
			return "", fmt.Errorf("the policyIDUnderClaimConsideration was not specified in transient data. Please try again")
		}
		accidentRecord.PolicyIDUnderClaimConsideration = string(policyIDUnderClaimConsideration)

		claimStatus, exists := transientData["claimStatus"]
		if !exists {
			return "", fmt.Errorf("the claimStatus was not specified in transient data. Please try again")
		}
		accidentRecord.ClaimStatus = string(claimStatus)

		accidentRecord.AssetType = "AccidentRecord"
		accidentRecord.AccidentID = accidentID

		bytes, _ := json.Marshal(accidentRecord)
		err = ctx.GetStub().PutPrivateData(collectionName, accidentID, bytes)
		if err != nil {
			return "", fmt.Errorf("could not able to write the data")
		}
		return fmt.Sprintf("accident record with id %v added successfully", accidentID), nil
	} else {
		return fmt.Sprintf("accident record cannot be created by organisation with MSPID %v ", clientOrgID), nil
	}
}

func (a *AccidentRecordContract) ReadAccidentRecord(ctx contractapi.TransactionContextInterface, accidentID string) (*AccidentRecord, error) {
	exists, err := a.AccidentRecordExists(ctx, accidentID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset (accident record) %s does not exist", accidentID)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, accidentID)
	if err != nil {
		return nil, fmt.Errorf("could not get the private data. %s", err)
	}
	var accidentRecord AccidentRecord

	err = json.Unmarshal(bytes, &accidentRecord)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private data collection data to type Order")
	}

	return &accidentRecord, nil

}

// DeleteOrder deletes an instance of Order from the private data collection
func (a *AccidentRecordContract) DeleteAccidentRecord(ctx contractapi.TransactionContextInterface, accidentID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("could not read the client identity. %s", err)
	}

	if clientOrgID == "police-insuranceClaimPostAccident-com" {

		exists, err := a.AccidentRecordExists(ctx, accidentID)

		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset (accident record) %s does not exist", accidentID)
		}

		return ctx.GetStub().DelPrivateData(collectionName, accidentID)
	} else {
		return fmt.Errorf("organisation with %v cannot delete the order", clientOrgID)
	}
}
