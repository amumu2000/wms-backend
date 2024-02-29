package inventory

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type listReq struct {
	GoodsID           *int64 `json:"goods_id"`
	WarehouseID       *int64 `json:"warehouse_id"`
	CountGreater      *int   `json:"count_greater"`
	CountGreaterEqual *bool  `json:"count_greater_equal"`
	CountLess         *int   `json:"count_less"`
	CountLessEqual    *bool  `json:"count_less_equal"`
}

func List(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := listReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	goodsID, warehouseID := int64(-1), int64(-1)
	countGreater, countGreaterEqual, countLess, countLessEqual := -1, false, -1, false

	if req.GoodsID != nil {
		goodsID = *req.GoodsID
	}

	if req.WarehouseID != nil {
		warehouseID = *req.WarehouseID
	}

	if req.CountGreater != nil {
		countGreater = *req.CountGreater
	}

	if req.CountGreaterEqual != nil {
		countGreaterEqual = *req.CountGreaterEqual
	}

	if req.CountLess != nil {
		countLess = *req.CountLess
	}

	if req.CountLessEqual != nil {
		countLessEqual = *req.CountLessEqual
	}

	if session.Role == 1 {
		if warehouseID < 0 {
			utils.Error400(c, nil, "请指定需要查询的仓库！")
			return
		}

		warehouse := models.GetWarehouseByID(warehouseID)
		if warehouse == nil || warehouse.ManagerID != session.UserID {
			utils.Error400(c, nil, "不可查询其他管理员仓库的库存。")
			return
		}
	}

	inventories := models.FindInventories(goodsID, warehouseID, countGreater, countLess, countGreaterEqual, countLessEqual)
	data := make([]gin.H, 0)

	for _, inventory := range inventories {
		good := models.GetGoodByID(inventory.GoodsID)
		warehouse := models.GetWarehouseByID(inventory.WarehouseID)
		manager := models.GetUserByID(warehouse.ManagerID)

		data = append(data, gin.H{
			"id":           inventory.ID,
			"goods_id":     inventory.GoodsID,
			"warehouse_id": inventory.WarehouseID,
			"count":        inventory.Count,
			"good": gin.H{
				"id":        good.ID,
				"good_name": good.GoodName,
				"category":  good.Category,
				"comment":   good.Comment,
				"price":     good.Price,
			},
			"warehouse": gin.H{
				"id":             warehouse.ID,
				"warehouse_name": warehouse.WarehouseName,
				"location":       warehouse.Location,
				"manager_id":     warehouse.ManagerID,
				"manager_name":   manager.Username,
				"comment":        warehouse.Comment,
			},
		})
	}

	utils.Status200(c, gin.H{
		"inventories": data,
	}, "")
}
