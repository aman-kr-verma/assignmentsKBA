package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PolicyContract struct {
	contractapi.Contract
}

type HistoryQueryResult struct {
	Record    *Policy `json:"record"`
	TxId      string  `json:"txId"`
	Timestamp string  `json:"timestamp"`
	IsDelete  bool    `json:"isDelete"`
}

type PaginatedQueryResult struct {
	Records             []*Policy `json:"records"`
	FetchedRecordsCount int32     `json:"fetchedRecordsCount"`
	Bookmark            string    `json:"bookmark"`
}

type Policy struct {
	AssetType        string  `json:"assetType"`
	PolicyID         string  `json:"policyID"`
	Coverage         string  `json:"coverage"`
	PeriodOfCoverage float32 `json:"periodOfCoverage"`
	Approved         bool    `json:"approved"`
	ClaimAmount      float64 `json:"claimAmount"`
}

// PolicyExists returns true when asset with given ID exists in world state
func (p *PolicyContract) PolicyExists(ctx contractapi.TransactionContextInterface, policyID string) (bool, error) {
	data, err := ctx.GetStub().GetState(policyID)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)

	}
	return data != nil, nil
}

// CreateCar creates a new instance of Car
func (c *PolicyContract) CreatePolicy(ctx contractapi.TransactionContextInterface, policyID string, coverage string, periodOfCoverage float32, claimAmount float64) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	// if clientOrgID == "insuranceCompany-insuranceClaimPostAccident-com" {
	if clientOrgID == "Org1MSP" {

		exists, err := c.PolicyExists(ctx, policyID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the policy, %s already exists", policyID)
		}

		policy := Policy{
			AssetType:        "Insurance Policy",
			PolicyID:         policyID,
			Coverage:         coverage,
			PeriodOfCoverage: periodOfCoverage,
			Approved:         false,
			ClaimAmount:      claimAmount,
		}

		bytes, _ := json.Marshal(policy)

		err = ctx.GetStub().PutState(policyID, bytes)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("successfully added policy %v", policyID), nil
		}

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// ReadCar retrieves an instance of Policy from the world state
func (p *PolicyContract) ReadPolicy(ctx contractapi.TransactionContextInterface, policyID string) (*Policy, error) {

	bytes, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the policy %s does not exist", policyID)
	}

	var policy Policy

	err = json.Unmarshal(bytes, &policy)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Car")
	}
	return &policy, nil
}

// approve func will approve the policy
func (p *PolicyContract) ApprovePolicy(ctx contractapi.TransactionContextInterface, policyID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	// if clientOrgID == "government-insuranceClaimPostAccident-com" {
	if clientOrgID == "Org2MSP" {

		policy, err := p.ReadPolicy(ctx, policyID)
		if err != nil {
			return "", fmt.Errorf("unable to approve policy (in reading phase): %s", err)
		}

		policy.Approved = true

		bytes, err := json.Marshal(policy)
		if err != nil {
			return "", fmt.Errorf("failed to marshal updated policy: %v", err)
		}

		err = ctx.GetStub().PutState(policyID, bytes)
		if err != nil {
			return "", fmt.Errorf("failed to update policy in world state: %v", err)
		}
		return "Approved", nil
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

func (p *PolicyContract) GetAllPolicies(ctx contractapi.TransactionContextInterface) ([]*Policy, error) {

	queryString := `{"selector":{"assetType":"Insurance Policy"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the query result. %s", err)
	}
	defer resultsIterator.Close()
	return policyResultIteratorFunction(resultsIterator)
}
func policyResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Policy, error) {
	var policies []*Policy
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not fetch the details of the result iterator. %s", err)
		}
		var policy Policy
		err = json.Unmarshal(queryResult.Value, &policy)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal the data. %s", err)
		}
		policies = append(policies, &policy)
	}

	return policies, nil
}

func (p *PolicyContract) GetPolicyHistory(ctx contractapi.TransactionContextInterface, policyID string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(policyID)
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

		var policy Policy
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &policy)
			if err != nil {
				return nil, err
			}
		} else {
			policy = Policy{
				PolicyID: policyID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &policy,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

func (p *PolicyContract) GetPoliciesWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	queryString := `{"selector":{"assetType":"Insurance Policy"}}`

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("could not get the policy records. %s", err)
	}
	defer resultsIterator.Close()

	policies, err := policyResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, fmt.Errorf("could not return the policy records %s", err)
	}

	return &PaginatedQueryResult{
		Records:             policies,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}
