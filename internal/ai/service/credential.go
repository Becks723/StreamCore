package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	"StreamCore/pkg/util"
	"gorm.io/gorm"
)

func SaveCredential(s *AIService, uid uint, req *kitexai.SaveCredentialReq) (*kitexai.SaveCredentialResp, error) {
	cred := &model.CredentialModel{
		UserID:      uid,
		ServiceName: req.ServiceName,
		Username:    req.Username,
		Password:    req.Password, // TODO: AES encrypt before storing
	}

	existing, err := s.db.GetCredential(s.ctx, uid, req.ServiceName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("SaveCredential: %w", err)
	}
	if existing != nil {
		cred.ID = existing.ID
		cred.CreatedAt = existing.CreatedAt
	}

	if err := s.db.SaveCredential(s.ctx, cred); err != nil {
		return nil, fmt.Errorf("SaveCredential: %w", err)
	}

	resp := new(kitexai.SaveCredentialResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Credential = pack.CredentialInfo(cred)
	return resp, nil
}

func DeleteCredential(s *AIService, req *kitexai.DeleteCredentialReq) (*kitexai.DeleteCredentialResp, error) {
	credID, err := util.ParseUint(req.CredentialId)
	if err != nil {
		return nil, fmt.Errorf("DeleteCredential: bad credential_id: %w", err)
	}

	if err := s.db.DeleteCredential(s.ctx, credID); err != nil {
		return nil, fmt.Errorf("DeleteCredential: %w", err)
	}

	resp := new(kitexai.DeleteCredentialResp)
	resp.Base = pack.BuildSuccessResp()
	return resp, nil
}

func ListCredentials(s *AIService, uid uint, req *kitexai.ListCredentialsReq) (*kitexai.ListCredentialsResp, error) {
	creds, err := s.db.ListCredentials(s.ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ListCredentials: %w", err)
	}

	items := make([]*kitexai.CredentialInfo, len(creds))
	for i, c := range creds {
		items[i] = pack.CredentialInfo(c)
	}

	resp := new(kitexai.ListCredentialsResp)
	resp.Base = pack.BuildSuccessResp()
	resp.Data = &kitexai.CredentialListData{Credentials: items}
	return resp, nil
}
