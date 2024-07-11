package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"log"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
    "github.com/rapolunagarjuna/receipt-processor-challenge/controllers"
    "github.com/rapolunagarjuna/receipt-processor-challenge/models"
	"bytes"
)

type MockReceiptService struct {
	mock.Mock
}

func (m *MockReceiptService) AddNewReceipt(receipt *models.Receipt) (string, int64) {
	args := m.Called(receipt)
	return args.String(0), args.Get(1).(int64)
}

func (m *MockReceiptService) GetReceipt(id string) (int64, bool) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Bool(1)
}


func TestProcessReceiptValidReceipt(t *testing.T) {
    router := gin.Default()
    mockService := MockReceiptService{}
	mockService.On("AddNewReceipt", mock.Anything).Return("1", int64(100))
	mockService.On("GetReceipt", "1").Return(int64(100), true)

    receiptController := controllers.ReceiptController{ReceiptService: &mockService}


    router.POST("/receipts/process", receiptController.ProcessReceipt)

    validReceipt := models.Receipt{
        Retailer: "Test Retailer",
        PurchaseDate: "2023-01-01",
        PurchaseTime: "12:00",
        Total: "10.00",
        Items: []models.Item{
            {ShortDescription: "Item 1", Price: "10.00"},
            {ShortDescription: "Item 2", Price: "20.00"},
        },
    }
    jsonBody, _ := json.Marshal(validReceipt)
	log.Println(string(jsonBody))

    req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")


    rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var response map[string]string
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "1", response["id"])
}

/*
Testing for 400 error code when the retailer is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidRetailer(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
    err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the purchase date is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidPurchaseDate(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the date is not in the format YYYY-MM-DD
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidDateFormat(t *testing.T) {
	router := gin.Default()
    mockService := MockReceiptService{}
	mockService.On("AddNewReceipt", mock.Anything).Return("1", int64(100))
	mockService.On("GetReceipt", "1").Return(int64(100), true)

    receiptController := controllers.ReceiptController{ReceiptService: &mockService}


    router.POST("/receipts/process", receiptController.ProcessReceipt)
	
	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "01-01-2023",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the purchase time is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidPurchaseTime(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the time is not in the format HH:MM
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidTimeFormat(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	mockService.On("AddNewReceipt", mock.Anything).Return("1", int64(100))
	mockService.On("GetReceipt", "1").Return(int64(100), true)

	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the total is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidTotal(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the total is not in the format ^\\d+\\.\\d{2}$
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidTotalFormat(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	mockService.On("AddNewReceipt", mock.Anything).Return("1", int64(100))
	mockService.On("GetReceipt", "1").Return(int64(100), true)

	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "10",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the items array is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidItems(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the short description of the item is empty or not provided
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidShortDescription(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "", Price: "10.00"},
			{ShortDescription: "Item 2", Price: "20.00"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 400 error code when the price of the item is not in the format ^\\d+\\.\\d{2}$
Also testing for the error message in the response
*/

func TestProcessReceiptInvalidItemPriceFormat(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	mockService.On("AddNewReceipt", mock.Anything).Return("1", int64(100))
	mockService.On("GetReceipt", "1").Return(int64(100), true)

	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.POST("/receipts/process", receiptController.ProcessReceipt)

	invalidReceipt := models.Receipt{
		Retailer: "Test Retailer",
		PurchaseDate: "2023-01-01",
		PurchaseTime: "12:00",
		Total: "10.00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10"},
			{ShortDescription: "Item 2", Price: "20"},
		},
	}
	jsonBody, _ := json.Marshal(invalidReceipt)

	req := httptest.NewRequest("POST", "http://example.com/receipts/process", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "The receipt is invalid", response["description"])
}

/*
Testing for 404 error code when the receipt id is not found
Also testing for the error message in the response
*/

func TestGetReceiptIdIsNotFound(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	mockService.On("GetReceipt", "1").Return(int64(0), false)

	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.GET("/receipts/:id/points", receiptController.GetReceiptPoints)

	req := httptest.NewRequest("GET", "http://example.com/receipts/1/points", nil)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "No receipt found for that id", response["description"])
}

/*
Testing for 200 success code when the receipt id is found
Also testing for the points in the response
*/

func TestGetReceiptIdIsFound(t *testing.T) {
	router := gin.Default()
	mockService := MockReceiptService{}
	mockService.On("GetReceipt", "1").Return(int64(100), true)

	receiptController := controllers.ReceiptController{ReceiptService: &mockService}

	router.GET("/receipts/:id/points", receiptController.GetReceiptPoints)

	req := httptest.NewRequest("GET", "http://example.com/receipts/1/points", nil)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	var response map[string]int64
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, int64(100), response["points"])
}
