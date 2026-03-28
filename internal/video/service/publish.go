package service

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"StreamCore/config"
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
	"github.com/google/uuid"
)

func (s *VideoService) Publish(uid uint, req *video.PublishReq) error {
	var err error
	var (
		localPrefix  = "./uploads"
		accessPrefix = "/static"
	)
	curUid := uid

	if !util.IsValidVideo(req.Data) {
		return errors.New("bad video format")
	}

	// exceeds video limit
	limit := config.Instance().General.VideoSizeLimit
	size := len(req.Data)
	if size > util.ToByte(limit) {
		return fmt.Errorf("exceeds video size limit (current %.2fmb but limits %.2fmb)", util.ToMb(size), limit)
	}

	// save video locally
	dir := localPrefix + accessPrefix + "/videos/"
	name := uuid.New().String()
	vdst := dir + name + ".mp4"
	if err = util.SaveFile(req.Data, vdst); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	// get cover
	var cdst string
	if req.CoverData != nil {
		if !util.IsValidImage(req.CoverData) {
			return errors.New("bad image format")
		}

		// exceeds image limit
		limit := config.Instance().General.ImageSizeLimit
		size := len(req.CoverData)
		if size > util.ToByte(limit) {
			return fmt.Errorf("exceeds image size limit (current %.2fmb but limits %.2fmb)", util.ToMb(size), limit)
		}

		cdst = dir + name + ".png"
		if err = util.SaveFile(req.CoverData, cdst); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	} else {
		cdst, err = s.randCover(vdst, dir)
		if err != nil {
			return fmt.Errorf("error getting cover: %w", err)
		}
	}

	var title, desc string
	if req.Title == nil {
		title = time.Now().Format("2006-01-02T15:04:05 -070000")
	} else {
		title = *req.Title
	}
	if req.Description != nil {
		desc = *req.Description
	}

	// update db
	vurl, _ := strings.CutPrefix(vdst, localPrefix)
	curl, _ := strings.CutPrefix(cdst, localPrefix)
	now := time.Now()
	v := domain.Video{
		AuthorId:    curUid,
		VideoUrl:    vurl,
		CoverUrl:    curl,
		Title:       title,
		Description: desc,
		PublishedAt: now,
		EditedAt:    now,
	}
	if err = s.db.Create(&v); err != nil {
		return fmt.Errorf("db.Create failed: %w", err)
	}

	// add to es
	err = s.es.AddVideo(s.ctx, &v)
	if err != nil {
		return err
	}
	return nil
}

func (s *VideoService) randCover(videoPath, coverDir string) (coverPath string, err error) {
	rand.New(rand.NewSource(time.Now().Unix()))

	// get video duration
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath)
	output, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("error retrieving video duration: %s", err.Error())
		return
	}

	// random a timepoint between 20% - 80%
	var duration float64
	_, _ = fmt.Sscanf(string(output), "%f", &duration)
	if duration <= 0 {
		err = errors.New("error reading video duration")
		return
	}
	sec := duration * (0.2 + 0.6*rand.Float64())

	// build cover path
	ext := filepath.Ext(videoPath)
	coverPath = coverDir +
		filepath.Base(videoPath[:len(videoPath)-len(ext)]) + ".jpg"

	// extract frame
	cmd = exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%.2f", sec),
		"-i", videoPath,
		"-frames:v", "1",
		"-q:v", "2",
		coverPath)
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("error writing to cover: %s", err.Error())
		return
	}
	return coverPath, nil
}
