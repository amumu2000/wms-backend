package goods

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type deleteReq struct {
	GoodsID *[]int64 `json:"goods_id"`
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

	goodsID := *req.GoodsID
	goodNames := make([]string, 0)

	for _, id := range goodsID {
		good := models.GetGoodByID(id)

		if good == nil {
			utils.BadRequest(c)
			return
		}

		inventories := models.GetInventories(id, -1)

		if len(inventories) > 0 {
			utils.Error400(c, nil, fmt.Sprintf("商品id %d 已有库存，无法删除。", id))
			return
		}

		goodNames = append(goodNames, good.GoodName)
	}

	models.DeleteGood(goodsID)

	for _, goodName := range goodNames {
		models.AddLog("删除商品", fmt.Sprintf("已删除商品%s", goodName), session.UserID)
	}

	utils.Status200(c, nil, "")
}
