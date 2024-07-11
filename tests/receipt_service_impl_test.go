package tests

import (
	"testing"
	"github.com/rapolunagarjuna/receipt-processor-challenge/services"
	"github.com/rapolunagarjuna/receipt-processor-challenge/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetReceipt(id string) (int64, bool) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Bool(1)
}

func (m *MockDB) AddNewReceipt(points int64) string {
	args := m.Called(points)
	return args.String(0)
}

/*
	testing whether the service is working as expected
	when a new receipt is added

	here we are adding a new receipt and checking whether the points are calculated correctly
	receipt taken 
	{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
			},{
			"shortDescription": "Emils Cheese Pizza",
			"price": "12.25"
			},{
			"shortDescription": "Knorr Creamy Chicken",
			"price": "1.26"
			},{
			"shortDescription": "Doritos Nacho Cheese",
			"price": "3.35"
			},{
			"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			"price": "12.00"
			}
		],
		"total": "35.35"
	Total Points: 28
		Breakdown:
			6 points - retailer name has 6 characters
			10 points - 4 items (2 pairs @ 5 points each)
			3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
						item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
			3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
						item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
			6 points - purchase day is odd
		+ ---------
		= 28 points
	}
*/
func TestAddNewReceiptTargetExample(t *testing.T) {
	assert := assert.New(t)
	dbMock := MockDB{}
	dbMock.On("AddNewReceipt", mock.Anything).Return("1")
	dbMock.On("GetReceipt", "1").Return(int64(100), true)

	
	receipt := models.Receipt{
		Retailer: "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total: "35.35",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price: "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price: "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price: "1.26",
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price: "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price: "12.00",
			},
		},
	}
	receiptService := services.ReceiptServiceImpl{
		DB: &dbMock,
	}

	id, points := receiptService.AddNewReceipt(&receipt)
	assert.Equal(int64(28), points)
	assert.Equal("1", id)
}

/*
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}

Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
*/
func TestAddNewReceiptMandMExample(t *testing.T) {
	assert := assert.New(t)
	dbMock := MockDB{}
	dbMock.On("AddNewReceipt", mock.Anything).Return("1")
	dbMock.On("GetReceipt", "1").Return(int64(100), true)


	receipt := models.Receipt{
		Retailer: "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{
				ShortDescription: "Gatorade",
				Price: "2.25",
			},{
				ShortDescription: "Gatorade",
				Price: "2.25",
			},{
				ShortDescription: "Gatorade",
				Price: "2.25",
			},{
				ShortDescription: "Gatorade",
				Price: "2.25",
			},
		},
		Total: "9.00",
	}

	receiptService := services.ReceiptServiceImpl{
		DB: &dbMock,
	}

	id, points := receiptService.AddNewReceipt(&receipt)
	assert.Equal(int64(109), points)
	assert.Equal("1", id)
}

func TestGetReceipt(t *testing.T) {
	dbMock := &MockDB{}
	dbMock.On("AddNewReceipt", mock.Anything).Return("1")
	dbMock.On("GetReceipt", "1").Return(int64(100), true)

	receiptService := services.ReceiptServiceImpl{
		DB: dbMock,
	}

	points, ok := receiptService.GetReceipt("1")

	assert.Equal(t, int64(100), points)
	assert.True(t, ok)
}

func TestGetReceiptNotFound(t *testing.T) {
	dbMock := &MockDB{}
	dbMock.On("AddNewReceipt", mock.Anything).Return("1")
	dbMock.On("GetReceipt", "1").Return(int64(0), false)

	receiptService := services.ReceiptServiceImpl{
		DB: dbMock,
	}

	points, ok := receiptService.GetReceipt("1")

	assert.Equal(t, int64(0), points)
	assert.False(t, ok)
}


/*
testing points for retailer name
*/

func TestPointsForRetailerName(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForRetailerName("Target")
	assert.Equal(int64(6), points)

	points = services.PointsForRetailerName("M&M Corner Market")
	assert.Equal(int64(14), points)
}

/*
testing points for receipt total
*/

func TestPointsForReceiptTotal(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForReceiptTotal("35.35")
	assert.Equal(int64(0), points)

	points = services.PointsForReceiptTotal("9.00")
	assert.Equal(int64(75), points)
}

/*
testing points for items
*/

func TestPointsForItems(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForItems([]models.Item{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price: "6.49",
		},
		{
			ShortDescription: "Emils Cheese Pizza",
			Price: "12.25",
		},
		{
			ShortDescription: "Knorr Creamy Chicken",
			Price: "1.26",
		},
		{
			ShortDescription: "Doritos Nacho Cheese",
			Price: "3.35",
		},
		{
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price: "12.00",
		},
	})

	assert.Equal(int64(10), points)
}

/*
testing points for item description
*/

func TestPointsForItemDescription(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForItemDescription([]models.Item{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price: "6.49",
		},
		{
			ShortDescription: "Emils Cheese Pizza",
			Price: "12.25",
		},
		{
			ShortDescription: "Knorr Creamy Chicken",
			Price: "1.26",
		},
		{
			ShortDescription: "Doritos Nacho Cheese",
			Price: "3.35",
		},
		{
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price: "12.00",
		},
	})

	assert.Equal(int64(6), points)
}

/*
testing points for receipt purchase date
*/

func TestPointsForReceiptPurchaseDate(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForReceiptPurchaseDate("2022-01-01")
	assert.Equal(int64(6), points)

	points = services.PointsForReceiptPurchaseDate("2022-03-20")
	assert.Equal(int64(0), points)
}

/*
testing points for receipt purchase time
*/

func TestPointsForReceiptPurchaseTime(t *testing.T) {
	assert := assert.New(t)

	points := services.PointsForReceiptPurchaseTime("13:01")
	assert.Equal(int64(0), points)

	points = services.PointsForReceiptPurchaseTime("14:33")
	assert.Equal(int64(10), points)
}




