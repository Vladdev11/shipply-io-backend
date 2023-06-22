package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomBarcode(length int) string {
	//TODO MAKE SURE UNIQUE
	rand.Seed(time.Now().UnixNano())
	barcode := ""
	for i := 0; i < length; i++ {
		digit := rand.Intn(10) // Generate a random digit between 0 and 9
		barcode += fmt.Sprintf("%d", digit)
	}
	return barcode
}
