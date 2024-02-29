package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type completeReq struct {
	DispatchID *[]int64 `json:"dispatch_id"`
}

func Complete(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 2)
	if !ok {
		return
	}

	req := completeReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	dispatchID := *req.DispatchID
	for _, id := range dispatchID {
		dispatch := models.GetDispatchByID(id)

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

		if dispatch.Status != 0 { //待完成
			utils.Error400(c, nil, fmt.Sprintf("派遣id %d 的状态不是待完成！", id))
			return
		}
	}

	dispatches := make([]map[string]interface{}, 0)

	for _, id := range dispatchID {
		dispatch := models.GetDispatchByID(id)

		dispatch.Status = 1 //已完成
		dispatch.UpdateTime = time.Now().Unix()
		dispatch.EndTime = time.Now().Unix()
		models.EditDispatch(*dispatch)

		good := models.GetGoodByID(dispatch.GoodsID)

		dispatches = append(dispatches, map[string]interface{}{
			"src":       dispatch.Src,
			"dest":      dispatch.Dest,
			"count":     dispatch.Count,
			"good_name": good.GoodName,
		})
	}

	for _, dispatch := range dispatches {
		src := dispatch["src"]
		dest := dispatch["dest"]
		count := dispatch["count"]
		goodName := dispatch["good_name"]
		models.AddLog("完成派遣", fmt.Sprintf("已完成派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", src, dest, goodName, count), session.UserID)
	}

	utils.Status200(c, nil, "")
}
