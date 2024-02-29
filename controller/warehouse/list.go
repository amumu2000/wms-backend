package warehouse

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type listRequest struct {
	WarehouseID   *int64  `json:"warehouse_id"`
	WarehouseName *string `json:"warehouse_name"`
	Location      *string `json:"location"`
	ManagerID     *int64  `json:"manager_id"`
}

func List(c *gin.Context) {
	if !utils.CheckToken(c, 99) {
		return
	}

	req := listRequest{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	warehouseID, name, location, managerID := int64(-1), "", "", int64(-1)

	if req.WarehouseID != nil {
		warehouseID = *req.WarehouseID
	}

	if req.WarehouseName != nil {
		name = *req.WarehouseName
	}

	if req.Location != nil {
		location = *req.Location
	}

	if req.ManagerID != nil {
		managerID = *req.ManagerID
	}

	warehouses := models.FindWarehouses(warehouseID, name, location, managerID)

	data := make([]gin.H, 0)
	for _, wh := range warehouses {
		data = append(data, gin.H{
			"id":             wh.ID,
			"warehouse_name": wh.WarehouseName,
			"location":       wh.Location,
			"manager_id":     wh.ManagerID,
			"comment":        wh.Comment,
		})
	}

	utils.Status200(c, gin.H{
		"warehouses": data,
	}, "")
}
