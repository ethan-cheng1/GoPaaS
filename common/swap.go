package common

import (
	"encoding/json"
)

// Assign struct through json tags
func SwapTo(request, target interface{}) (err error)  {
	dataByte, err := json.Marshal(request)
	if err != nil {
		return
	}
	err = json.Unmarshal(dataByte, target)
	return
}