package codekit

import "time"

type MockModel struct {
	ID       string `json:"_id"`
	Name     string
	Age      int
	Salary   float64
	IsActive bool
	JoinDate time.Time
	Tags     []string
	Children map[string]MockChildModel
	Pointer  *MockChildModel
}

type MockChildModel struct {
	Key   string
	Value string
}
