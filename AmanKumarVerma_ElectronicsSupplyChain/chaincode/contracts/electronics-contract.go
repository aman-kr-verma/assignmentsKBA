package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ElectronicsContract struct {
	contractapi.Contract
}

type EventData struct {
	Type  string
	Model string
}

type HistoryQueryResult struct {
	Record    *ElectronicItem `json:"record"`
	TxId      string          `json:"txId"`
	Timestamp string          `json:"timestamp"`
	IsDelete  bool            `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records             []*ElectronicItem `json:"records"`
	FetchedRecordsCount int32             `json:"fetchedRecordsCount"`
	Bookmark            string            `json:"bookmark"`
}

type ElectronicItem struct {
	AssetType         string `json:"assetType"`
	ItemId            string `json:"itemId"`
	ItemType          string `json:"itemType"`
	Model             string `json:"model"`
	Make              string `json:"make"`
	Color             string `json:"color"`
	DateOfManufacture string `json:"dateOfManufacture"`
	OwnedBy           string `json:"ownedBy"`
	Status            string `json:"status"`
}

// ElectronicItemExists returns true when asset with given ID exists in world state
func (e *ElectronicsContract) ElectronicItemExists(ctx contractapi.TransactionContextInterface, itemID string) (bool, error) {
	data, err := ctx.GetStub().GetState(itemID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)

	}
	return data != nil, nil
}

// CreateElectronicItem creates a new instance of electronic item
func (e *ElectronicsContract) CreateElectronicItem(ctx contractapi.TransactionContextInterface, itemID string, itemType string, model string, make string, color string, dateOfManufacture string, manufacturerName string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "ManufacturerMSP" {

		exists, err := e.ElectronicItemExists(ctx, itemID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the electronic item, %s already exists", itemID)
		}

		item := ElectronicItem{
			AssetType:         "electronics",
			ItemId:            itemID,
			ItemType:          itemType,
			Color:             color,
			DateOfManufacture: dateOfManufacture,
			Make:              make,
			Model:             model,
			OwnedBy:           manufacturerName,
			Status:            "In Factory",
		}

		// fmt.Println("Create electronic item data ======= ", item)

		bytes, _ := json.Marshal(item)

		err = ctx.GetStub().PutState(itemID, bytes)
		if err != nil {
			return "", err
		} else {
			addItemEventData := EventData{
				Type:  "Item creation",
				Model: model,
			}
			eventDataByte, _ := json.Marshal(addItemEventData)
			ctx.GetStub().SetEvent("CreateElectronicItem", eventDataByte)

			return fmt.Sprintf("successfully added electronic item %v", itemID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// ReadElectronicItem retrieves an instance of electronic item from the world state
func (e *ElectronicsContract) ReadElectronicItem(ctx contractapi.TransactionContextInterface, itemID string) (*ElectronicItem, error) {

	bytes, err := ctx.GetStub().GetState(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the electronic item %s does not exist", itemID)
	}

	var item ElectronicItem

	err = json.Unmarshal(bytes, &item)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Electronic Item")
	}
	return &item, nil
}

func (e *ElectronicsContract) DeleteElectronicItem(ctx contractapi.TransactionContextInterface, itemID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "ManufacturerMSP" {
		exists, err := e.ElectronicItemExists(ctx, itemID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if !exists {
			return "", fmt.Errorf("the electronic item, %s does not exist", itemID)
		}

		err = ctx.GetStub().DelState(itemID)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("electronic item with id %v is deleted from the world state.", itemID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (e *ElectronicsContract) GetAllElectronicItems(ctx contractapi.TransactionContextInterface) ([]*ElectronicItem, error) {

	queryString := `{"selector":{"assetType":"electronics"}, "sort":[{ "color": "desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return itemResultIteratorFunction(resultsIterator)
}

func itemResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*ElectronicItem, error) {
	var items []*ElectronicItem
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var item ElectronicItem
		err = json.Unmarshal(queryResult.Value, &item)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (e *ElectronicsContract) GetElectronicItemsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*ElectronicItem, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the data by range. %s", err)
	}
	defer resultsIterator.Close()

	return itemResultIteratorFunction(resultsIterator)
}

func (e *ElectronicsContract) GetElectronicItemHistory(ctx contractapi.TransactionContextInterface, itemID string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(itemID)
	if err != nil {
		return nil, fmt.Errorf("could not get the data. %s", err)
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get the value of resultsIterator. %s", err)
		}

		var item ElectronicItem
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &item)
			if err != nil {
				return nil, err
			}
		} else {
			item = ElectronicItem{
				ItemId: itemID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &item,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

func (e *ElectronicsContract) GetElectronicsItemsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"assetType":"electronics"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the electronic item records. %s", err)
	}
	defer resultsIterator.Close()

	items, err := itemResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the electronic items records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             items,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
