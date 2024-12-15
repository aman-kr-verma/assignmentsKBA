package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var EventPayload string
var chaincodeEventsList []string
var blockEventsList []string
var pvtEventsList []string

func ChaincodeEventListener(organization string, channelName string, chaincodeName string, wg *sync.WaitGroup) {
	defer wg.Done()
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n*** Start Chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start Chaincode event listening: %w", err))
	}

	for event := range events {
		fmt.Println("Received chaincode event", event.EventName)
		decodedPayload := string(event.Payload)

		chaincodeEventsList = append(chaincodeEventsList, decodedPayload)
	}

}

func BlockEventListener(organization string, channelName string, wg *sync.WaitGroup) {
	defer wg.Done()
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\n*** Start Block event listening")

	events, err := network.BlockEvents(ctx, client.WithStartBlock(4))
	if err != nil {
		panic(fmt.Errorf("failed to start Block event listening: %w", err))
	}

	for event := range events {
		fmt.Println("Received block number", event.GetHeader().GetNumber())
		decodedPayload := string(event.GetHeader().GetNumber())

		blockEventsList = append(blockEventsList, decodedPayload)
	}

}

func PvtblockListener(organization string, channelName string, wg *sync.WaitGroup) {
	defer wg.Done()
	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("***started listening to private txn in blocks***")
	events, err := network.BlockAndPrivateDataEvents(ctx, client.WithStartBlock(1))

	if err != nil {
		panic(fmt.Errorf("failed to start Block event listening: %w", err))
	}
	// .GetHeader().GetNumber()

	for event := range events {
		if event.GetPrivateDataMap() != nil {
			//fmt.Println("private data event is:",event)
			fmt.Printf("Received block: %d with pvt data \n", event.GetBlock().GetHeader().GetNumber())
		}
		decodedPayload := string(event.GetBlock().Header.Number)

		pvtEventsList = append(pvtEventsList, decodedPayload)
	}

}

func getChaincodeEvents() []string {
	return chaincodeEventsList
}

func getBlockEvents() []string {
	return blockEventsList
}

func getpvtEvents() []string {
	return pvtEventsList
}
