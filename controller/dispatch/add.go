package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type addReq struct {
	ExecutorID    *int64  `json:"executor_id"`
	GoodsID       *int64  `json:"goods_id"`
	Status        *int    `json:"status"`
	Comment       *string `json:"comment"`
	StartTime     *int64  `json:"start_time"`
	ExpectEndTime *int64  `json:"expect_end_time"`
	Src           *string `json:"src"`
	SrcID         *int64  `json:"src_id"`
	Dest          *string `json:"dest"`
	DestID        *int64  `json:"dest_id"`
	Type          *int    `json:"type"`
	Count         *int    `json:"count"`
}

func Add(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := addReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	if session.Role == 0 {
		utils.Error400(c, nil, "管理员不能添加派遣！")
		return
	}

	managerID := session.UserID
	executorID := *req.ExecutorID
	goodsID := *req.GoodsID
	status := *req.Status
	comment := *req.Comment
	startTime := *req.StartTime
	expectEndTime := *req.ExpectEndTime
	src := *req.Src
	dest := *req.Dest
	srcID := *req.SrcID
	destID := *req.DestID
	ttype := *req.Type
	count := *req.Count

	if ttype != 0 && ttype != 1 {
		utils.BadRequest(c)
		return
	}

	if !(status >= 0 && status <= 3) {
		utils.BadRequest(c)
		return
	}

	if count <= 0 {
		utils.Error400(c, nil, "商品数量不能为0！")
		return
	}

	executor := models.GetUserByID(executorID)
	if executor == nil {
		utils.Error400(c, nil, "派遣人员不存在！")
		return
	}

	good := models.GetGoodByID(goodsID)
	if good == nil {
		utils.Error400(c, nil, "商品不存在！")
		return
	}

	if ttype == 0 {
		//入库
		if destID < 0 {
			utils.Error400(c, nil, "必需入库到自己管理的仓库！")
			return
		}

		warehouseDest := models.GetWarehouseByID(destID)
		if warehouseDest == nil || warehouseDest.ManagerID != managerID {
			utils.Error400(c, nil, "必需入库到自己管理的仓库！")
			return
		}

		dest = warehouseDest.WarehouseName

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
		if srcID < 0 {
			utils.Error400(c, nil, "必需从自己管理的仓库出库！")
			return
		}

		warehouseSrc := models.GetWarehouseByID(srcID)
		if warehouseSrc == nil || warehouseSrc.ManagerID != managerID {
			utils.Error400(c, nil, "必需从自己管理的仓库出库！")
			return
		}

		src = warehouseSrc.WarehouseName

		inventorySrc := models.GetInventories(goodsID, srcID)
		if len(inventorySrc) == 0 || inventorySrc[0].Count < count {
			utils.Error400(c, nil, "库存不足！")
			return
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

	dispatch := models.Dispatch{
		ManagerID:     managerID,
		ExecutorID:    executorID,
		GoodsID:       goodsID,
		Status:        status,
		Comment:       comment,
		CreateTime:    time.Now().Unix(),
		StartTime:     startTime,
		ExpectEndTime: expectEndTime,
		UpdateTime:    time.Now().Unix(),
		Src:           src,
		SrcID:         srcID,
		Dest:          dest,
		DestID:        destID,
		Type:          ttype,
		Count:         count,
	}

	dispatchID := models.AddDispatch(dispatch)

	models.AddLog("添加派遣", fmt.Sprintf("已添加派遣，源仓库%s，目标仓库%s，商品为%s，数量%d", src, dest, good.GoodName, count), session.UserID)

	utils.Status200(c, gin.H{
		"dispatch_id": dispatchID,
	}, "")
}
