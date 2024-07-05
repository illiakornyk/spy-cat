package common

type Target struct {
	ID        int64  `json:"id,omitempty"`
	MissionID int64  `json:"mission_id,omitempty"`
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Country   string `json:"country" validate:"required,min=1,max=100"`
	Notes     string `json:"notes" validate:"max=500"`
	Complete  bool   `json:"complete"`
}
