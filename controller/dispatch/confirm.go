package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type confirmReq struct {
	DispatchID *int64 `json:"dispatch_id"`
}

func Confirm(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := confirmReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	dispatchID := *req.DispatchID

	dispatch := models.GetDispatchByID(dispatchID)

	if dispatch == nil {
		utils.BadRequest(c)
		return
	}

	if session.Role == 1 {
		if dispatch.ManagerID != session.UserID {
			utils.BadRequest(c)
			return
		}
	}

	if dispatch.Status != 1 { //已完成
		utils.Error400(c, nil, fmt.Sprintf("派遣id %d 的状态不是已完成！", dispatchID))
		return
	}

	if dispatch.SrcID > 0 {
		inventories := models.GetInventories(dispatch.GoodsID, dispatch.SrcID)

		if len(inventories) == 0 || inventories[0].Count < dispatch.Count {
			utils.Error400(c, nil, fmt.Sprintf("派遣id %d 无法确认，源库库存不足！", dispatchID))
			return
		}
	}

	if dispatch.SrcID > 0 {
		inventories := models.GetInventories(dispatch.GoodsID, dispatch.SrcID)
		inventory := inventories[0]

		inventory.Count -= dispatch.Count
		if inventory.Count > 0 {
			models.EditInventory(inventory)
		} else {
			models.DeleteInventory([]int64{inventory.ID})
		}
	}

	if dispatch.DestID > 0 {
		inventories := models.GetInventories(dispatch.GoodsID, dispatch.DestID)

		if len(inventories) > 0 {
			inventory := inventories[0]
			inventory.Count += dispatch.Count

			models.EditInventory(inventory)
		} else {
			models.AddInventory(dispatch.GoodsID, dispatch.DestID, dispatch.Count)
		}
	}

	dispatch.Status = 3 //已完成
	dispatch.UpdateTime = time.Now().Unix()

	models.EditDispatch(*dispatch)

	good := models.GetGoodByID(dispatch.GoodsID)

	models.AddLog("已确认派遣", fmt.Sprintf("已确认派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", dispatch.Src, dispatch.Dest, good.GoodName, dispatch.Count), session.UserID)

	utils.Status200(c, nil, "")
}
