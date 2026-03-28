package domain

import "time"

type Follow struct {
	TargetUid uint
	StartedAt time.Time

	FollowerId uint
	FolloweeId uint
	Status     int // 0-关注 1-取关
	Time       time.Time
}
