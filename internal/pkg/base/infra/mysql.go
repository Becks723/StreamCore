package infra

import (
	"fmt"
	"net/url"

	"StreamCore/config"
	"StreamCore/internal/pkg/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL() (*gorm.DB, error) {
	// connect to mysql database
	c := config.Instance().MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName)

	params := url.Values{}
	params.Set("charset", c.Charset)
	params.Set("parseTime", "true")

	dsn += fmt.Sprintf("?%s", params.Encode())
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening mysql: %w", err)
	}

	// auto migrate
	err = db.AutoMigrate(
		&model.UserModel{},
		&model.VideoModel{},
		&model.VisitCountModel{},
		&model.LikeRelationModel{},
		&model.LikeCountModel{},
		&model.CommentModel{},
		&model.FollowModel{},
		&model.ChatMsgModel{},
		&model.ChatGroupModel{},
		&model.ChatGroupMemberModel{},
		&model.ChatGroupApplyModel{},
	)
	if err != nil {
		return nil, fmt.Errorf("error auto migrating: %w", err)
	}

	return db, nil
}
