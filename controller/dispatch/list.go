package dispatch

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type listReq struct {
	ManagerID  *int64 `json:"manager_id"`
	ExecutorID *int64 `json:"executor_id"`
	Status     *int   `json:"status"`
}

func List(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 2)
	if !ok {
		return
	}

	req := listReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	managerID, executorID, status := int64(-1), int64(-1), -1

	if req.ManagerID != nil {
		managerID = *req.ManagerID
	}

	if req.ExecutorID != nil {
		executorID = *req.ExecutorID
	}

	if req.Status != nil {
		status = *req.Status
	}

	if session.Role == 1 && managerID >= 0 {
		utils.Error400(c, nil, "不可查询其他管理员仓库的派遣！")
		return
	}

	if session.Role == 2 {
		if managerID >= 0 || executorID >= 0 {
			utils.Error400(c, nil, "不可查询其他人的派遣！")
			return
		}
	}

	if session.Role == 1 {
		managerID = session.UserID
	} else if session.Role == 2 {
		executorID = session.UserID
	}

	dispatches := models.GetDispatches(managerID, executorID, status)
	data := make([]gin.H, 0)

	for _, dispatch := range dispatches {
		data = append(data, gin.H{
			"id":              dispatch.ID,
			"manager_id":      dispatch.ManagerID,
			"executor_id":     dispatch.ExecutorID,
			"goods_id":        dispatch.GoodsID,
			"status":          dispatch.Status,
			"comment":         dispatch.Comment,
			"create_time":     dispatch.CreateTime,
			"start_time":      dispatch.StartTime,
			"expect_end_time": dispatch.ExpectEndTime,
			"end_time":        dispatch.EndTime,
			"update_time":     dispatch.UpdateTime,
			"src":             dispatch.Src,
			"src_id":          dispatch.SrcID,
			"dest":            dispatch.Dest,
			"dest_id":         dispatch.DestID,
			"type":            dispatch.Type,
			"count":           dispatch.Count,
			"reject_comment":  dispatch.RejectComment,
		})
	}

	utils.Status200(c, gin.H{
		"dispatches": data,
	}, "")
}
