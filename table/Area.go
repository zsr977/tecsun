package table

import "time"

type Area struct {
	Id         int       `json:"id"`
	Uuid       string    `json:"uuid"`
	Company    string    `json:"company"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Directions string    `json:"directions"`
	Remark     string    `json:"remark"`
	IsDelete   uint8     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
