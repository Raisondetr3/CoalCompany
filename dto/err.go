package dto

import (
	"encoding/json"
	"time"
)

type Err struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func (e Err) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}

	return string(b)
}
