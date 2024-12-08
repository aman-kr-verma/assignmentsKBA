package main

import (
	"electronics/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	electronicItemContract := new(contracts.ElectronicsContract)
	rawMaterialContract := new(contracts.RawMaterialContract)

	chaincode, err := contractapi.NewChaincode(electronicItemContract, rawMaterialContract)

	if err != nil {
		log.Panicf("could not create chaincode %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("could not start chaincode %v", err)
	}
}
