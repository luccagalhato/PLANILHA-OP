package models

import "io"

// Login ...
type Login struct {
	Username     string `json:"username,omitempty"`
	Userpassword string `json:"userpassword,omitempty"`
}
type Logout struct {
	Token string `json:"token,omitempty"`
}

//Response ...
type Response struct {
	Status string
	Error  string
	Data   interface{}
}

//Excel ...
type Excel struct {
	Name  string
	Value io.Reader
}

//ResponseExcel ...
type ResponseExcel struct {
	Status string
	Error  string
	Data   interface{}
	Excel  interface{}
}
