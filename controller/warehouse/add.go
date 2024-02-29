package warehouse

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminAddReq struct {
	WarehouseName *string `json:"warehouse_name"`
	Location      *string `json:"location"`
	ManagerID     *int64  `json:"manager_id"`
	Comment       *string `json:"comment"`
}

func AdminAdd(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminAddReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	name := *req.WarehouseName
	location := *req.Location
	managerID := *req.ManagerID
	comment := *req.Comment

	manager := models.GetUserByID(managerID)

	if manager == nil || manager.Role != 1 {
		utils.Error400(c, nil, "仓库管理员不存在！")
		return
	}

	warehouseID := models.AddWarehouse(name, location, managerID, comment)

	models.AddLog("添加仓库", fmt.Sprintf("已添加仓库%s", name), session.UserID)

	utils.Status200(c, gin.H{
		"warehouse_id": warehouseID,
	}, "")
}
