package dto

type PowerAction struct {
	Action int `json:"action" binding:"required" example:"8"`
}
