package table

import "time"

type WorkflowForm struct {
	Id         int       `json:"id"`
	Uuid       string    `json:"uuid"`
	WfUuid     string    `json:"wf_uuid"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Name       string    `json:"name"`
	Default    string    `json:"default"`
	Multiple   uint8     `json:"multiple"`
	Required   uint8     `json:"required"`
	List       string    `json:"list"`
	Sort       uint8     `json:"sort"`
	IsDelete   uint8     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
