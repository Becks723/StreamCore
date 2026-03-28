package model

import "time"

type LikeEvent struct {
	TarType int       `json:"tar_type"`
	TarId   uint      `json:"tar_id"`
	Uid     uint      `json:"uid"`
	Action  int       `json:"action"`
	Time    time.Time `json:"time"`
}
