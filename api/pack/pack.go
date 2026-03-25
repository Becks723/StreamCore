package pack

import (
	"StreamCore/kitex_gen/common"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type base struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type respWithData struct {
	Base *base `json:"base"`
	Data any   `json:"data,omitempty"`
}

func RespRPCError(c *app.RequestContext, err error) {
	c.JSON(consts.StatusInternalServerError, respWithData{
		Base: &base{
			Code: consts.StatusInternalServerError,
			Msg:  err.Error(),
		},
	})
}

func RespBizError(c *app.RequestContext, resp *common.BaseResp) bool {
	if resp.Code != consts.StatusOK {
		c.JSON(int(resp.Code), respWithData{
			Base: &base{
				Code: int64(resp.Code),
				Msg:  resp.Msg,
			},
		})
		return true
	}
	return false
}

func RespParamError(c *app.RequestContext, err error) {
	c.JSON(consts.StatusBadRequest, respWithData{
		Base: &base{
			Code: consts.StatusBadRequest,
			Msg:  "Invalid parameter: " + err.Error(),
		},
	})
}

func RespUnauthorizedError(c *app.RequestContext, err error) {
	c.JSON(consts.StatusUnauthorized, respWithData{
		Base: &base{
			Code: consts.StatusUnauthorized,
			Msg:  err.Error(),
		},
	})
}

func RespSuccess(c *app.RequestContext) {
	c.JSON(consts.StatusOK, respWithData{
		Base: &base{
			Code: consts.StatusOK,
			Msg:  "",
		},
	})
}

func RespWithData(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, respWithData{
		Base: &base{
			Code: consts.StatusOK,
			Msg:  "",
		},
		Data: data,
	})
}
