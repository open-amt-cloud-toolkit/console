package dto

type PowerState struct {
	PowerState int `json:"powerstate" binding:"required" example:"0"`
}
