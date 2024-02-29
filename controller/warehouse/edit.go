package warehouse

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminEditReq struct {
	WarehouseID   *int64  `json:"warehouse_id"`
	WarehouseName *string `json:"warehouse_name"`
	Location      *string `json:"location"`
	ManagerID     *int64  `json:"manager_id"`
	Comment       *string `json:"comment"`
}

func AdminEdit(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminEditReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	whid := *req.WarehouseID

	name, location, managerID, comment := "", "", int64(-1), ""

	if req.WarehouseName != nil {
		name = *req.WarehouseName
	}

	if req.Location != nil {
		location = *req.Location
	}

	if req.ManagerID != nil {
		managerID = *req.ManagerID
	}

	if req.Comment != nil {
		comment = *req.Comment
	}

	if managerID >= 0 {
		manager := models.GetUserByID(managerID)

		if manager == nil || manager.Role != 1 {
			utils.Error400(c, nil, "仓库管理员不存在！")
			return
		}
	}

	warehouse := models.GetWarehouseByID(whid)

	if name != "" {
		warehouse.WarehouseName = name
	}

	if location != "" {
		warehouse.Location = location
	}

	if managerID >= 0 {
		warehouse.ManagerID = managerID
	}

	if comment != "" {
		warehouse.Comment = comment
	}

	models.EditWarehouse(*warehouse)

	models.AddLog("修改仓库", fmt.Sprintf("已修改仓库%s", warehouse.WarehouseName), session.UserID)

	utils.Status200(c, nil, "")
}
