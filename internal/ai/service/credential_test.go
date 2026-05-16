package service

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestAIService_Credential(t *testing.T) {
	type testCase struct {
		name           string
		expectingError bool
		expectedError  string
		credentialID   string
		serviceName    string
		username       string
		password       string
		getError       error
		existingCred   *model.CredentialModel
		saveError      error
		deleteError    error
		listError      error
	}

	uid := uint(234)

	testCases := []testCase{
		{
			name:           "SaveCredential-新建成功",
			expectingError: false,
			serviceName:    "video-cdn",
			username:       "admin",
			password:       "secret",
			getError:       gorm.ErrRecordNotFound,
		},
		{
			name:           "SaveCredential-更新已有",
			expectingError: false,
			serviceName:    "video-cdn",
			username:       "admin",
			password:       "new-secret",
			existingCred: &model.CredentialModel{
				ID:          1,
				UserID:      uid,
				ServiceName: "video-cdn",
				Username:    "admin",
			},
		},
		{
			name:           "SaveCredential-GetCredentialDB错误",
			expectingError: true,
			expectedError:  "db error",
			serviceName:    "video-cdn",
			username:       "admin",
			getError:       errors.New("db error"),
		},
		{
			name:           "SaveCredential-Save DB错误",
			expectingError: true,
			expectedError:  "save failed",
			serviceName:    "video-cdn",
			username:       "admin",
			getError:       gorm.ErrRecordNotFound,
			saveError:      errors.New("save failed"),
		},
		{
			name:           "DeleteCredential-成功",
			expectingError: false,
			credentialID:   "1",
		},
		{
			name:           "DeleteCredential-无效ID",
			expectingError: true,
			expectedError:  "bad credential_id",
			credentialID:   "abc",
		},
		{
			name:           "DeleteCredential-DB错误",
			expectingError: true,
			expectedError:  "delete failed",
			credentialID:   "1",
			deleteError:    errors.New("delete failed"),
		},
		{
			name:           "ListCredentials-成功",
			expectingError: false,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			_GetCredential := mockey.GetMethod(dbai.NewAIDatabase(nil), "GetCredential")
			mockey.Mock(_GetCredential).To(func(ctx context.Context, uid uint, serviceName string) (*model.CredentialModel, error) {
				if tc.getError != nil {
					return nil, tc.getError
				}
				return tc.existingCred, nil
			}).Build()

			_SaveCredential := mockey.GetMethod(dbai.NewAIDatabase(nil), "SaveCredential")
			mockey.Mock(_SaveCredential).To(func(ctx context.Context, cred *model.CredentialModel) error {
				if tc.saveError != nil {
					return tc.saveError
				}
				cred.ID = 1
				cred.CreatedAt = time.Now()
				cred.UpdatedAt = time.Now()
				return nil
			}).Build()

			_DeleteCredential := mockey.GetMethod(dbai.NewAIDatabase(nil), "DeleteCredential")
			mockey.Mock(_DeleteCredential).To(func(ctx context.Context, credID uint) error {
				if tc.deleteError != nil {
					return tc.deleteError
				}
				return nil
			}).Build()

			_ListCredentials := mockey.GetMethod(dbai.NewAIDatabase(nil), "ListCredentials")
			mockey.Mock(_ListCredentials).To(func(ctx context.Context, uid uint) ([]*model.CredentialModel, error) {
				if tc.listError != nil {
					return nil, tc.listError
				}
				return []*model.CredentialModel{
					{ID: 1, UserID: uid, ServiceName: "video-cdn", Username: "admin", CreatedAt: time.Now()},
				}, nil
			}).Build()

			mockInfraSet := &base.InfraSet{DB: db.NewDatabaseSet(nil)}
			svc := NewAIService(context.Background(), mockInfraSet, &agent.Agent{}, &mcp.ToolRegistry{})

			switch {
			case tc.name == "SaveCredential-新建成功" || tc.name == "SaveCredential-更新已有" ||
				tc.name == "SaveCredential-GetCredentialDB错误" || tc.name == "SaveCredential-Save DB错误":
				resp, err := SaveCredential(svc, uid, &kitexai.SaveCredentialReq{
					ServiceName: tc.serviceName,
					Username:    tc.username,
					Password:    tc.password,
				})
				if tc.expectingError {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, tc.expectedError)
				} else {
					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(resp.Base.Code, ShouldEqual, 200)
					So(resp.Credential, ShouldNotBeNil)
					if tc.existingCred != nil {
						So(resp.Credential.CredentialId, ShouldEqual, util.Uint2String(tc.existingCred.ID))
					}
				}

			case tc.name == "DeleteCredential-成功" || tc.name == "DeleteCredential-无效ID" ||
				tc.name == "DeleteCredential-DB错误":
				resp, err := DeleteCredential(svc, &kitexai.DeleteCredentialReq{
					CredentialId: tc.credentialID,
				})
				if tc.expectingError {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, tc.expectedError)
				} else {
					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
				}

			case tc.name == "ListCredentials-成功":
				resp, err := ListCredentials(svc, uid, &kitexai.ListCredentialsReq{})
				if tc.expectingError {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(resp.Data, ShouldNotBeNil)
					So(len(resp.Data.Credentials), ShouldBeGreaterThan, 0)
				}
			}
		})
	}
}
