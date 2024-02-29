package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type rejectReq struct {
	DispatchID    *int64  `json:"dispatch_id"`
	RejectComment *string `json:"reject_comment"`
}

func Reject(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 2)
	if !ok {
		return
	}

	req := rejectReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	dispatchID := *req.DispatchID
	rejectComment := *req.RejectComment

	if rejectComment == "" {
		utils.Error400(c, nil, "必需填写无法完成的理由。")
		return
	}

	dispatch := models.GetDispatchByID(dispatchID)
	if dispatch == nil {
		utils.BadRequest(c)
		return
	}

	if session.Role == 2 {
		if dispatch.ExecutorID != session.UserID {
			utils.BadRequest(c)
			return
		}
	}

	if dispatch.Status != 0 {
		utils.Error400(c, nil, fmt.Sprintf("派遣id %d 的状态不是待完成！", dispatchID))
		return
	}

	dispatch.Status = 2 //无法完成
	dispatch.RejectComment = rejectComment
	dispatch.UpdateTime = time.Now().Unix()
	dispatch.EndTime = time.Now().Unix()

	models.EditDispatch(*dispatch)

	good := models.GetGoodByID(dispatch.GoodsID)

	models.AddLog("无法完成派遣", fmt.Sprintf("无法完成派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", dispatch.Src, dispatch.Dest, good.GoodName, dispatch.Count), session.UserID)

	utils.Status200(c, nil, "")
}
