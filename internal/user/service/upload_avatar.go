package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"StreamCore/config"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
)

func (s *UserService) UploadAvatar(uid uint, data []byte) (*common.UserInfo, error) {
	var err error
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid := uid

	if !util.IsValidImage(data) {
		return nil, errors.New("bad image format")
	}

	// exceeds image limit
	limit := config.Instance().General.ImageSizeLimit
	size := len(data)
	if size > util.ToByte(limit) {
		return nil, fmt.Errorf("exceeds image size limit (current %.2fmb but limits %.2fmb)", util.ToMb(size), limit)
	}

	// save image locally
	dst := fmt.Sprintf(localPrefix+accessPrefix+"/avatars/%d_%d.png", // TODO: match extensions
		curUid, time.Now().Unix())
	err = util.SaveFile(data, dst)
	if err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// update db
	newUrl, _ := strings.CutPrefix(dst, localPrefix)
	u, err := s.db.UpdateAvatar(curUid, newUrl)
	if err != nil {
		return nil, fmt.Errorf("error updating database avatar: %w", err)
	}

	return pack.UserInfo(u), nil
}
