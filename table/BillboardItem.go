package table

import "time"

type BillboardItem struct {
	Id            int       `json:"id"`
	Uuid          string    `json:"uuid"`
	BillboardUuid string    `json:"billboard_uuid"`
	TenantUuid    string    `json:"tenant_uuid"`
	Status        uint8     `json:"status"`
	IsDelete      uint8     `json:"is_delete"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	DeleteTime    time.Time `json:"delete_time"`
}
