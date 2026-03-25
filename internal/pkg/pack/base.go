package pack

import (
	"StreamCore/kitex_gen/common"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func BuildSuccessResp() *common.BaseResp {
	return BuildBaseResp(nil)
}

func BuildBaseResp(err error) *common.BaseResp {
	if err == nil {
		return &common.BaseResp{
			Code: consts.StatusOK,
			Msg:  "ok",
		}
	}

	return &common.BaseResp{
		Code: consts.StatusInternalServerError,
		Msg:  err.Error(),
	}
}
