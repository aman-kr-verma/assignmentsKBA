package main

import "fmt"

func main() {
	result := submitTxnFn(
		"manufacturer",
		"electronicschannel",
		"electronicsSC",
		"ElectronicsContract",
		"invoke",
		make(map[string][]byte),
		"CreateElectronicItem",
		"Item09",
		"TV",
		"black",
		"22/07/2024",
		"LG",
		"LG52in",
		"Electronic Manufacturers LTD",
	)

	// privateData := map[string][]byte{
	// 	"materialType":      []byte("Iron"),
	// 	"quantity":          []byte("1 Ton"),
	// 	"supplierId":        []byte("SUP001"),
	// 	"dateOfManufacture": []byte("22/02/2023"),
	// 	"model":             []byte("MOD001"),
	// }

	// result := submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "private", privateData, "CreateRawMaterial", "MAT06")

	// result := submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "query", make(map[string][]byte), "ReadRawMaterial", "MAT01")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetAllElectronicItems")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicItemsByRange","Item01","Item04")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicItemHistory","Item01")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicsItemsWithPagination","3","")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "RawMaterialContract", "query", make(map[string][]byte), "GetAllRawMaterials")

	// result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "RawMaterialContract", "query", make(map[string][]byte), "GetRawMaterialsByRange","MAT01","MAT04")

	// result := submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "private", make(map[string][]byte), "DeleteRawMaterial", "MAT05")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "GetMatchingOrders", "Car-06")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "MatchOrder", "Car-06", "ORD-03")

	// result := submitTxnFn("mvd", "autochannel", "KBA-Automobile", "CarContract", "invoke", make(map[string][]byte), "RegisterCar", "Car-06", "Dani", "KL-01-CD-01")

	// result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "ReadCar", "Car-06")

	fmt.Println(result)

}
