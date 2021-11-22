package table

import "time"

type Answer struct {
	Id         int       `json:"id"`
	Uuid       string    `json:"uuid"`
	UserUuid   string    `json:"user_uuid"`
	NaireUuid  string    `json:"naire_uuid"`
	IsDelete   uint8     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
