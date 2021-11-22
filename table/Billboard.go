package table

import "time"

type Billboard struct {
	Id         int       `json:"id"`
	Uuid       string    `json:"uuid"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	IsDelete   uint8     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
