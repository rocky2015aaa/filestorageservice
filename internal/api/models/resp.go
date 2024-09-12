package models

type Resp struct {
	Success     bool        `json:"success" example:"true"`
	Data        interface{} `json:"data"`
	Error       string      `json:"error" example:""`
	Description string      `json:"description" example:"proper message for the situation"`
}
