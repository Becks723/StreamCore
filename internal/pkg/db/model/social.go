package model

import "time"

type FollowModel struct {
	Id         uint `gorm:"primaryKey"`
	FollowerId uint
	Follower   *UserModel `gorm:"foreignKey:FollowerId;references:ID"`
	FolloweeId uint
	Followee   *UserModel `gorm:"foreignKey:FolloweeId;references:ID"`
	Status     int        // 0-关注 1-取关
	Time       time.Time
}
