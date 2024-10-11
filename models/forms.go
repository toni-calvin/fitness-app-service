package models

type CreateMesocycleForm struct {
	NumberMicrocycles string `json:"numberMicrocycles"`
	StartDate         string `json:"startDate"`
	Objectives        string `json:"objectives"`
}
