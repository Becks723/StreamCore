package service

import (
	"context"
	"errors"
	"testing"

	"StreamCore/internal/pkg/ai/agent"
	"StreamCore/internal/pkg/ai/mcp"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/db"
	dbai "StreamCore/internal/pkg/db/ai"
	"StreamCore/internal/pkg/db/model"
	kitexai "StreamCore/kitex_gen/ai"
	"StreamCore/pkg/util"
	"github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func TestAIService_CreateBot(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.PatchConvey("成功", t, func() {
		_CreateBotUser := mockey.GetMethod(dbai.NewAIDatabase(nil), "CreateBotUser")
		mockey.Mock(_CreateBotUser).To(func(ctx context.Context, bot *model.UserModel) error {
			bot.ID = 1
			return nil
		}).Build()

		mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
		svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

		resp, err := CreateBot(svc, &kitexai.CreateBotReq{
			BotName:      "TestBot",
			SystemPrompt: "test",
			Provider:     "openai-default",
			TriggerMode:  0,
		})
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Base.Code, ShouldEqual, 200)
		So(resp.Bot.BotId, ShouldEqual, "1")
	})
}

func TestAIService_GetBot(t *testing.T) {
	type testCase struct {
		name           string
		expectingError bool
		expectedError  string
		botID          string
		getUserError   error
	}

	botUser := &model.UserModel{
		Model:    gorm.Model{ID: 1},
		Username: "TestBot",
		IsBot:    true,
		BotConfig: func() *string {
			s := `{"system_prompt":"test","provider":"openai-default","trigger_mode":0,"tool_ids":[]}`
			return &s
		}(),
	}

	testCases := []testCase{
		{
			name:           "成功",
			expectingError: false,
			botID:          util.Uint2String(botUser.ID),
		},
		{
			name:           "不存在",
			expectingError: true,
			expectedError:  "bot not found",
			botID:          util.Uint2String(999),
			getUserError:   gorm.ErrRecordNotFound,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			_GetBotUser := mockey.GetMethod(dbai.NewAIDatabase(nil), "GetBotUser")
			mockey.Mock(_GetBotUser).To(func(ctx context.Context, botID uint) (*model.UserModel, error) {
				if tc.getUserError != nil {
					return nil, tc.getUserError
				}
				return botUser, nil
			}).Build()

			mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
			svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

			resp, err := GetBot(svc, &kitexai.GetBotReq{BotId: tc.botID})

			if tc.expectingError {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.expectedError)
			} else {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.Bot.BotId, ShouldEqual, tc.botID)
			}
		})
	}
}

func TestAIService_UpdateBot(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.PatchConvey("数据库失败", t, func() {
		_GetBotUser := mockey.GetMethod(dbai.NewAIDatabase(nil), "GetBotUser")
		mockey.Mock(_GetBotUser).To(func(ctx context.Context, botID uint) (*model.UserModel, error) {
			return nil, errors.New("db error")
		}).Build()

		mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
		svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

		_, err := UpdateBot(svc, &kitexai.UpdateBotReq{BotId: "1"})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "get bot user failed")
	})
}

func TestAIService_DeleteBot(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.PatchConvey("数据库失败", t, func() {
		_DeleteBotUser := mockey.GetMethod(dbai.NewAIDatabase(nil), "DeleteBotUser")
		mockey.Mock(_DeleteBotUser).To(func(ctx context.Context, botID uint) error {
			return errors.New("db error")
		}).Build()

		mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
		svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

		_, err := DeleteBot(svc, &kitexai.DeleteBotReq{BotId: "1"})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "delete user failed")
	})
}

func TestAIService_ListBots(t *testing.T) {
	botUser := &model.UserModel{
		Model:    gorm.Model{ID: 1},
		Username: "TestBot",
		IsBot:    true,
		BotConfig: func() *string {
			s := `{"system_prompt":"test","provider":"openai-default","trigger_mode":0,"tool_ids":[]}`
			return &s
		}(),
	}

	defer mockey.UnPatchAll()
	mockey.PatchConvey("成功", t, func() {
		_ListBotUsers := mockey.GetMethod(dbai.NewAIDatabase(nil), "ListBotUsers")
		mockey.Mock(_ListBotUsers).To(func(ctx context.Context, page, pageSize int) ([]*model.UserModel, int64, error) {
			return []*model.UserModel{botUser}, 3, nil
		}).Build()

		mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
		svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

		pageSize := int32(20)
		page := int32(1)
		resp, err := ListBots(svc, &kitexai.ListBotsReq{
			PageSize: &pageSize,
			Page:     &page,
		})
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(*resp.Total, ShouldEqual, 3)
	})
}
