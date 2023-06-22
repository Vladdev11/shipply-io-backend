package util

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(stuff interface{}) {
	orderPrint1, _ := json.MarshalIndent(stuff, "", "    ")
	fmt.Printf("%s\n", string(orderPrint1))
}
