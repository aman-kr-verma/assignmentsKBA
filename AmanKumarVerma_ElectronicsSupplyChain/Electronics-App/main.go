package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ElectronicItem struct {
	ItemId            string `json:"itemId"`
	ItemType          string `json:"itemType"`
	Model             string `json:"model"`
	Make              string `json:"make"`
	Color             string `json:"color"`
	DateOfManufacture string `json:"dateOfManufacture"`
	OwnedBy           string `json:"ownedBy"`
}

type ElectronicItemData struct {
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

type RawMaterial struct {
	MaterialId        string `json:"materialId"`
	MaterialType      string `json:"materialType"`
	Quantity          string `json:"quantity"`
	SupplierId        string `json:"supplierId"`
	DateOfManufacture string `json:"dateOfManufacture"`
	Model             string `json:"model"`
}

type RawMaterialData struct {
	AssetType         string `json:"assetType"`
	MaterialId        string `json:"materialId"`
	MaterialType      string `json:"materialType"`
	Quantity          string `json:"quantity"`
	SupplierId        string `json:"supplierId"`
	DateOfManufacture string `json:"dateOfManufacture"`
	Model             string `json:"model"`
}

type ItemHistory struct {
	Record    *ElectronicItem `json:"record"`
	TxId      string          `json:"txId"`
	Timestamp string          `json:"timestamp"`
	IsDelete  bool            `json:"isDelete"`
}

func main() {
	router := gin.Default()

	var wg sync.WaitGroup
	wg.Add(1)
	go ChaincodeEventListener("manufacturer", "electronicschannel", "electronicsSC", &wg)
	wg.Add(1)
	go BlockEventListener("manufacturer", "electronicschannel", &wg)
	wg.Add(1)
	go PvtblockListener("manufacturer", "electronicschannel", &wg)

	// Get Electronic items
	router.GET("/api/electronicItems", func(ctx *gin.Context) {
		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetAllElectronicItems")

		var items []ElectronicItemData

		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &items); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"data": items})
	})

	// Add Electronic item
	router.POST("/api/electronicItem", func(ctx *gin.Context) {
		var req ElectronicItem
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("electronic item response %s", req)
		submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "invoke",
			make(map[string][]byte), "CreateElectronicItem", req.ItemId, req.ItemType, req.Color, req.DateOfManufacture, req.Make, req.Model, req.OwnedBy)

		ctx.JSON(http.StatusOK, req)
	})

	// Read electronic item
	router.GET("/api/electronicItem/:id", func(ctx *gin.Context) {
		itemId := ctx.Param("id")

		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "ReadElectronicItem", itemId)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// Read Items by range
	router.GET("/api/electronicItem/Range", func(ctx *gin.Context) {
		from := ctx.Query("from")
		to := ctx.Query("to")

		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicItemsByRange", from, to)

		ctx.JSON(http.StatusOK, gin.H{"items in range ": result})
	})

	// Electronic item history
	router.GET("/api/electronicItem/history/:id", func(ctx *gin.Context) {
		itemID := ctx.Param("id")
		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicItemHistory", itemID)

		var items []ItemHistory

		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &items); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"data": items})

	})

	// Electronic item delete
	router.DELETE("/api/electronicItem/delete/:id", func(ctx *gin.Context) {
		itemId := ctx.Param("id")

		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "invoke", make(map[string][]byte), "DeleteElectronicItem", itemId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// Electronic Items pagination
	router.GET("/api/electronicItem/pagination", func(ctx *gin.Context) {
		pageSize := ctx.Query("pageSize")
		bookmark := ctx.DefaultQuery("bookmark", "")

		result := submitTxnFn("manufacturer", "electronicschannel", "electronicsSC", "ElectronicsContract", "query", make(map[string][]byte), "GetElectronicsItemsWithPagination", pageSize, bookmark)

		//return bookmark
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// Electronic item chaincode event
	router.GET("/api/chaincodeEvent", func(ctx *gin.Context) {
		result := getChaincodeEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"electronic item chaincode Event": result})

	})

	// Electronic item block event
	router.GET("/api/blockEvent", func(ctx *gin.Context) {
		result := getBlockEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"electronic item block Event": result})

	})

	// Electronic item privateData event
	router.GET("/api/pvtBlockEvent", func(ctx *gin.Context) {
		result := getpvtEvents()
		fmt.Println("result:", result)

		ctx.JSON(http.StatusOK, gin.H{"electronic item private data block Event": result})

	})

	//Get all raw materials
	router.GET("/api/rawMaterials", func(ctx *gin.Context) {

		result := submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "query", make(map[string][]byte), "GetAllRawMaterials")

		var rawMaterials []RawMaterialData

		if len(result) > 0 {
			if err := json.Unmarshal([]byte(result), &rawMaterials); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"data": rawMaterials})

	})

	router.POST("/api/rawMaterial", func(ctx *gin.Context) {
		var req RawMaterial
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		fmt.Printf("raw material  %s", req)

		privateData := map[string][]byte{
			"materialType":      []byte(req.MaterialType),
			"quantity":          []byte(req.Quantity),
			"supplierId":        []byte(req.SupplierId),
			"dateOfManufacture": []byte(req.DateOfManufacture),
			"model":             []byte(req.Model),
		}

		submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "private", privateData, "CreateRawMaterial", req.MaterialId)
		ctx.JSON(http.StatusOK, req)
	})

	router.GET("/api/rawMaterial/:id", func(ctx *gin.Context) {
		materialId := ctx.Param("id")

		result := submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "query", make(map[string][]byte), "ReadRawMaterial", materialId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// Raw material delete
	router.DELETE("/api/rawMaterial/delete/:id", func(ctx *gin.Context) {
		itemId := ctx.Param("id")

		submitTxnFn("supplier", "electronicschannel", "electronicsSC", "RawMaterialContract", "private", make(map[string][]byte), "DeleteRawMaterial", itemId)

		ctx.JSON(http.StatusOK, gin.H{"raw material deleted successfully ": itemId})

	})
	router.Run("localhost:8080")

}
