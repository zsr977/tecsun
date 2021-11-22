package table

import "time"

type Booking struct {
	Id         int    `json:"id"`
	Uuid       string `json:"uuid"`
	TenantUuid string `json:"tenant_uuid"`

	IsDelete   uint8     `json:"is_delete"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DeleteTime time.Time `json:"delete_time"`
}
