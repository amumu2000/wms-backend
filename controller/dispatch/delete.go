package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type deleteReq struct {
	DispatchID *[]int64 `json:"dispatch_id"`
}

func Delete(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := deleteReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	dispatchID := *req.DispatchID
	dispatches := make([]map[string]interface{}, 0)

	for _, id := range dispatchID {
		dispatch := models.GetDispatchByID(id)

		if dispatch == nil {
			utils.BadRequest(c)
			return
		}

		if session.Role == 1 {
			if dispatch.ManagerID != session.UserID {
				utils.Error400(c, nil, "不能删除其他仓库管理员的派遣！")
				return
			}
		}

		good := models.GetGoodByID(dispatch.GoodsID)

		dispatches = append(dispatches, map[string]interface{}{
			"src":       dispatch.Src,
			"dest":      dispatch.Dest,
			"count":     dispatch.Count,
			"good_name": good.GoodName,
		})
	}

	models.DeleteDispatch(dispatchID)

	for _, dispatch := range dispatches {
		src := dispatch["src"]
		dest := dispatch["dest"]
		count := dispatch["count"]
		goodName := dispatch["good_name"]
		models.AddLog("删除派遣", fmt.Sprintf("已删除派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", src, dest, goodName, count), session.UserID)
	}

	utils.Status200(c, nil, "")
}
