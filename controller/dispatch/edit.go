package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type editReq struct {
	DispatchID    *int64  `json:"dispatch_id"`
	GoodsID       *int64  `json:"goods_id"`
	Status        *int    `json:"status"`
	Comment       *string `json:"comment"`
	StartTime     *int64  `json:"start_time"`
	ExpectEndTime *int64  `json:"expect_end_time"`
	EndTime       *int64  `json:"end_time"`
	Src           *string `json:"src"`
	SrcID         *int64  `json:"src_id"`
	Dest          *string `json:"dest"`
	DestID        *int64  `json:"dest_id"`
	Type          *int    `json:"type"`
	Count         *int    `json:"count"`
	RejectComment *string `json:"reject_comment"`
}

func Edit(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := editReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	if session.Role == 0 {
		utils.Error400(c, nil, "管理员不能修改派遣！")
		return
	}

	dispatchID := *req.DispatchID
	managerID := session.UserID

	goodsID, status := int64(-1), -1
	comment, startTime, expectEndTime, endTime := "", int64(-1), int64(-1), int64(-1)
	src, dest, srcID, destID := "", "", int64(-1), int64(-1)
	ttype, count, rejectComment := -1, -1, ""
	hasComment, hasRejectComment := false, false

	if req.GoodsID != nil {
		goodsID = *req.GoodsID
	}

	if req.Status != nil {
		status = *req.Status
	}

	if req.Comment != nil {
		comment = *req.Comment
		hasComment = true
	}

	if req.StartTime != nil {
		startTime = *req.StartTime
	}

	if req.ExpectEndTime != nil {
		expectEndTime = *req.ExpectEndTime
	}

	if req.EndTime != nil {
		endTime = *req.EndTime
	}

	if req.Src != nil {
		src = *req.Src
	}

	if req.SrcID != nil {
		srcID = *req.SrcID
	}

	if req.Dest != nil {
		dest = *req.Dest
	}

	if req.DestID != nil {
		destID = *req.DestID
	}

	if req.Type != nil {
		ttype = *req.Type
	}

	if req.Count != nil {
		count = *req.Count
	}

	if req.RejectComment != nil {
		rejectComment = *req.RejectComment
		hasRejectComment = true
	}

	if ttype >= 0 {
		if ttype != 0 && ttype != 1 {
			utils.BadRequest(c)
			return
		}
	}

	if status >= 0 {
		if !(status >= 0 && status <= 3) {
			utils.BadRequest(c)
			return
		}
	}

	if count == 0 {
		utils.Error400(c, nil, "商品数量不能为0！")
	}

	var good *models.Good

	if goodsID >= 0 {
		good = models.GetGoodByID(goodsID)
		if good == nil {
			utils.Error400(c, nil, "商品不存在！")
			return
		}
	}

	dispatch := models.GetDispatchByID(dispatchID)
	if dispatch == nil {
		utils.BadRequest(c)
		return
	}

	if dispatch.ManagerID != managerID {
		utils.Error400(c, nil, "不能修改其他管理员的派遣！")
		return
	}

	if ttype == 0 {
		//入库
		if destID > 0 {
			warehouseDest := models.GetWarehouseByID(destID)
			if warehouseDest == nil || warehouseDest.ManagerID != managerID {
				utils.Error400(c, nil, "必需入库到自己管理的仓库！")
				return
			}

			dest = warehouseDest.WarehouseName
		}

		if srcID > 0 {
			warehouseSrc := models.GetWarehouseByID(srcID)
			if warehouseSrc == nil {
				utils.Error400(c, nil, "仓库不存在！")
				return
			}

			src = warehouseSrc.WarehouseName
		}
	} else if ttype == 1 {
		//出库
		if srcID > 0 {
			warehouseSrc := models.GetWarehouseByID(srcID)
			if warehouseSrc == nil || warehouseSrc.ManagerID != managerID {
				utils.Error400(c, nil, "必需从自己管理的仓库出库！")
				return
			}

			src = warehouseSrc.WarehouseName

			goodsIDTmp := goodsID
			if goodsIDTmp < 0 {
				goodsIDTmp = dispatch.GoodsID
			}

			inventorySrc := models.GetInventories(goodsIDTmp, srcID)
			if len(inventorySrc) == 0 || inventorySrc[0].Count < count {
				utils.Error400(c, nil, "库存不足！")
				return
			}
		}

		if destID > 0 {
			warehouseDest := models.GetWarehouseByID(destID)
			if warehouseDest == nil {
				utils.Error400(c, nil, "仓库不存在！")
				return
			}

			dest = warehouseDest.WarehouseName
		}
	}

	if goodsID >= 0 {
		dispatch.GoodsID = goodsID
	}

	if status >= 0 {
		dispatch.Status = status
	}

	if hasComment {
		dispatch.Comment = comment
	}

	if startTime >= 0 {
		dispatch.StartTime = startTime
	}

	if expectEndTime >= 0 {
		dispatch.ExpectEndTime = expectEndTime
	}

	if endTime >= 0 {
		dispatch.EndTime = endTime
	}

	if src != "" {
		dispatch.Src = src
	}

	if srcID >= 0 {
		dispatch.SrcID = srcID
	}

	if dest != "" {
		dispatch.Dest = dest
	}

	if destID >= 0 {
		dispatch.DestID = destID
	}

	if ttype >= 0 {
		dispatch.Type = ttype
	}

	if count >= 0 {
		dispatch.Count = count
	}

	if hasRejectComment {
		dispatch.RejectComment = rejectComment
	}

	dispatch.UpdateTime = time.Now().Unix()

	models.EditDispatch(*dispatch)

	models.AddLog("修改派遣", fmt.Sprintf("已修改派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", dispatch.Src, dispatch.Dest, good.GoodName, count), session.UserID)

	utils.Status200(c, nil, "")
}
