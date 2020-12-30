package models

import "time"

type (
	Command struct {
		Id         int64     `json:"id"`
		AgentId    string    `json:"agent_id"`
		Content    string    `json:"content"`
		Status     int       `json:"status"`
		CreateTime time.Time `xorm:"created"`
		UpdateTime time.Time `xorm:"updated"`
	}
)
