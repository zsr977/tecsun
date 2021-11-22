package table

import "time"

type AnswerDetail struct {
	Id           int       `json:"id"`
	Uuid         string    `json:"uuid"`
	AnswerUuid   string    `json:"answer_uuid"`
	QuestionUuid string    `json:"question_uuid"`
	Content      string    `json:"content"`
	IsDelete     uint8     `json:"is_delete"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	DeleteTime   time.Time `json:"delete_time"`
}
