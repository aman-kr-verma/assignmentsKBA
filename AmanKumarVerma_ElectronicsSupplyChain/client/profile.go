package main

type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

var profile = map[string]Config{

	"manufacturer": {
		CertPath:     "../electronics_network/organizations/peerOrganizations/manufacturer.electronics.com/users/User1@manufacturer.electronics.com/msp/signcerts/cert.pem",
		KeyDirectory: "../electronics_network/organizations/peerOrganizations/manufacturer.electronics.com/users/User1@manufacturer.electronics.com/msp/keystore/",
		TLSCertPath:  "../electronics_network/organizations/peerOrganizations/manufacturer.electronics.com/peers/peer0.manufacturer.electronics.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.manufacturer.electronics.com",
		MSPID:        "ManufacturerMSP",
	},

	"dealer": {
		CertPath:     "../electronics_network/organizations/peerOrganizations/dealer.electronics.com/users/User1@dealer.electronics.com/msp/signcerts/cert.pem",
		KeyDirectory: "../electronics_network/organizations/peerOrganizations/dealer.electronics.com/users/User1@dealer.electronics.com/msp/keystore/",
		TLSCertPath:  "../electronics_network/organizations/peerOrganizations/dealer.electronics.com/peers/peer0.dealer.electronics.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.dealer.electronics.com",
		MSPID:        "DealerMSP",
	},

	"supplier": {
		CertPath:     "../electronics_network/organizations/peerOrganizations/supplier.electronics.com/users/User1@supplier.electronics.com/msp/signcerts/cert.pem",
		KeyDirectory: "../electronics_network/organizations/peerOrganizations/supplier.electronics.com/users/User1@supplier.electronics.com/msp/keystore/",
		TLSCertPath:  "../electronics_network/organizations/peerOrganizations/supplier.electronics.com/peers/peer0.supplier.electronics.com/tls/ca.crt",
		PeerEndpoint: "localhost:11051",
		GatewayPeer:  "peer0.supplier.electronics.com",
		MSPID:        "SupplierMSP",
	},
}