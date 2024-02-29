package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type receipt struct {
	ID string `json:"id"`
	Points int `json:"points"`
    Retailer string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items []receiptItem `json:"items"`
	Total string `json:"total"`
}

type receiptItem struct {
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}

var receipts []receipt
var numReceipts int

func main() {
    router := gin.Default()

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)

    router.Run("localhost:8080")
}

// Searches the receipts slice for a receipt matching the ID and returns that receipt.
func getPoints(c *gin.Context) {
    id := c.Param("id")

    for _, a := range receipts {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, gin.H{ "points": a.Points })
            return
		}
    }
	
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Receipt not found"})
}

// Takes in a receipt, calculates the points, adds it to the receipts slice, and returns the ID.
func processReceipt(c *gin.Context) {
	var newReceipt receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid receipt"})
		return
	}
	newReceipt.ID = strconv.Itoa(numReceipts)
	numReceipts++
	calculatePoints(&newReceipt)
	receipts = append(receipts, newReceipt)
	c.IndentedJSON(http.StatusAccepted, gin.H{"id": newReceipt.ID})
}

// Calculates how many points a receipt is worth
func calculatePoints(receipt *receipt) {
	receipt.Points += calculateRetailerPoints(receipt.Retailer)
	receipt.Points += calculateTotalPoints(receipt.Total)
	receipt.Points += len(receipt.Items) / 2 * 5
	receipt.Points += calculateDescriptionPoints(receipt.Items)
	receipt.Points += calculateDatePoints(receipt.PurchaseDate)
	receipt.Points += calculateTimePoints(receipt.PurchaseTime)
}

// Calculates how many points the retailer field is worth
func calculateRetailerPoints(retailer string) int {
	var returnValue int
	for _, r := range retailer {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			returnValue++
		}
	}

	return returnValue
}

// Calculates how many points the total field is worth
func calculateTotalPoints(receiptTotal string) int {
	var returnValue int
	total, _ := strconv.ParseFloat(receiptTotal, 64)

	// Adds 50 points if the total is a whole number
	if total == float64(int(total)) {
		returnValue += 50
	}
	
	// Adds 25 points if the total is a multiple of 0.25
	if total == float64(int(total * 4) / 4) {
		returnValue += 25
	}

	return returnValue
}

// Calculates how many points the description field is worth
func calculateDescriptionPoints(items []receiptItem) int { 
	var returnValue int
	for _, i := range items {
		price, _ := strconv.ParseFloat(i.Price, 64)
		trimmedDescription := strings.TrimSpace(i.ShortDescription)

		if len(trimmedDescription) % 3 == 0 {
			returnValue += int(math.Ceil(price * 0.2))
		}
	}
	return returnValue
 }

// Calculates how many points the date field is worth
func calculateDatePoints(date string) int {
	dateComponents := strings.Split(date, "-")
	day, _ := strconv.Atoi(dateComponents[2])

	if day % 2 != 0 { // If the day is odd
		return 6
	} else {
		return 0
	}
}

// Calculates how many points the time field is worth
func calculateTimePoints(time string) int {
	timeComponents := strings.Split(time, ":")
	hour, _ := strconv.Atoi(timeComponents[0])
	minutes, _ := strconv.Atoi(timeComponents[1])

	if hour < 14 || hour > 16 { // If the time is not between 2:00 PM and 4:00 PM
		return 0
	} else if hour == 14 && minutes == 0 { // If the time is exactly 2:00 PM
		return 0 
	} else { 
		return 10 
	}
}