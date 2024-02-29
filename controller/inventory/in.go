package inventory

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type inReq struct {
	GoodsID     *int64 `json:"goods_id"`
	WarehouseID *int64 `json:"warehouse_id"`
	Count       *int   `json:"count"`
}

func In(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := inReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	goodsID := *req.GoodsID
	warehouseID := *req.WarehouseID
	count := *req.Count

	if count <= 0 {
		utils.BadRequest(c)
		return
	}

	good := models.GetGoodByID(goodsID)
	if good == nil {
		utils.BadRequest(c)
		return
	}

	warehouse := models.GetWarehouseByID(warehouseID)
	if warehouse == nil {
		utils.BadRequest(c)
		return
	}

	if session.Role == 1 {
		if warehouse.ManagerID != session.UserID {
			utils.Error400(c, nil, "无法入库其他管理员的仓库！")
			return
		}
	}

	inventories := models.GetInventories(goodsID, warehouseID)
	if len(inventories) > 0 {
		inventory := inventories[0]

		inventory.Count += count
		models.EditInventory(inventory)
	} else {
		models.AddInventory(goodsID, warehouseID, count)
	}

	models.AddLog("入库", fmt.Sprintf("商品%s已入仓库%s", good.GoodName, warehouse.WarehouseName), session.UserID)

	utils.Status200(c, nil, "")
}
