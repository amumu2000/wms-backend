package models

type Dispatch struct {
	ID            int64  `gorm:"column:id"`
	ManagerID     int64  `gorm:"column:manager_id"`
	ExecutorID    int64  `gorm:"column:executor_id"`
	GoodsID       int64  `gorm:"column:goods_id"`
	Status        int    `gorm:"column:status"`
	Comment       string `gorm:"column:comment"`
	CreateTime    int64  `gorm:"column:create_time"`
	StartTime     int64  `gorm:"column:start_time"`
	ExpectEndTime int64  `gorm:"column:expect_end_time"`
	EndTime       int64  `gorm:"column:end_time"`
	UpdateTime    int64  `gorm:"column:update_time"`
	Src           string `gorm:"column:src"`
	SrcID         int64  `gorm:"column:src_id"`
	Dest          string `gorm:"column:dest"`
	DestID        int64  `gorm:"column:dest_id"`
	Type          int    `gorm:"column:type"`
	Count         int    `gorm:"column:count"`
	RejectComment string `gorm:"column:reject_comment"`
}

func GetDispatchByID(id int64) *Dispatch {
	var count int64
	dispatch := Dispatch{}
	db.Table("dispatch").Where("id = ?", id).Count(&count).First(&dispatch)

	if count == 0 {
		return nil
	} else {
		return &dispatch
	}
}

func GetDispatches(managerID, executorID int64, status int) []Dispatch {
	query := db.Table("dispatch")

	if managerID >= 0 {
		query = query.Where("manager_id = ?", managerID)
	}

	if executorID >= 0 {
		query = query.Where("executor_id = ?", executorID)
	}

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	var dispatches []Dispatch
	query.Find(&dispatches)

	return dispatches
}

func AddDispatch(dispatch Dispatch) int64 {
	db.Table("dispatch").Create(&dispatch)

	return dispatch.ID
}

func DeleteDispatch(id []int64) {
	db.Table("dispatch").Delete(&Dispatch{}, id)
}

func EditDispatch(dispatch Dispatch) {
	db.Table("dispatch").Save(&dispatch)
}
