package main

import (
	"insurance/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	policyContract := new(contracts.PolicyContract)
	accidentRecordContract := new(contracts.AccidentRecordContract)

	chaincode, err := contractapi.NewChaincode(policyContract, accidentRecordContract)

	if err != nil {
		log.Panicf("could not create chaincode %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("could not start chaincode %v", err)
	}
}
