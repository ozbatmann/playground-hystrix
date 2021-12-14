package models

type LimitCheckResponse struct {
	collectOperationAvailable bool     `json:"collectOperationAvailable"`
	active                    bool      `json:"active"`
}
