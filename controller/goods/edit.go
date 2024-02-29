package goods

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type editReq struct {
	GoodsID  *int64  `json:"goods_id"`
	GoodName *string `json:"good_name"`
	Category *string `json:"category"`
	Price    *int    `json:"price"`
	Comment  *string `json:"comment"`
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

	goodsID := *req.GoodsID
	goodName, category, price, comment := "", "", -1, ""

	if req.GoodName != nil {
		goodName = *req.GoodName
	}

	if req.Category != nil {
		category = *req.Category
	}

	if req.Price != nil {
		price = *req.Price
	}

	if req.Comment != nil {
		comment = *req.Comment
	}

	good := models.GetGoodByID(goodsID)

	if good == nil {
		utils.BadRequest(c)
		return
	}

	if goodName != "" {
		good.GoodName = goodName
	}

	if category != "" {
		good.Category = category
	}

	if price >= 0 {
		good.Price = price
	}

	if comment != "" {
		good.Comment = comment
	}

	models.EditGood(*good)

	models.AddLog("编辑商品", fmt.Sprintf("已编辑商品%s", good.GoodName), session.UserID)

	utils.Status200(c, nil, "")
}
