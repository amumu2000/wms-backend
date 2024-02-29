package warehouse

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminDeleteReq struct {
	WarehouseID *[]int64 `json:"warehouse_id"`
}

func AdminDelete(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminDeleteReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	whids := *req.WarehouseID
	names := make([]string, 0)

	for _, id := range whids {
		warehouse := models.GetWarehouseByID(id)

		if warehouse == nil {
			utils.BadRequest(c)
			return
		}

		names = append(names, warehouse.WarehouseName)
	}

	models.DeleteWarehouse(whids)

	for _, name := range names {
		models.AddLog("删除仓库", fmt.Sprintf("已删除仓库%s", name), session.UserID)
	}

	utils.Status200(c, nil, "")
}
