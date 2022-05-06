package affair

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

type PageForm struct {
	PageNum    int    `form:"pageNum" json:"pageNum" binding:"required"`
	PageSize   int    `form:"pageSize" json:"pageSize" binding:"required"`
	Department string `form:"department" json:"department"`
	Theme      string `form:"theme" json:"theme"`
}

// @Summary 政务办理分页
// @Tags Affair
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param department query string false "政务部门"
// @Param theme query string false "政务主题"
// @Router /pullAffairPage [get]
func pullAffairPage(r *gin.Engine) {
	r.GET("/pullAffairPage", func(c *gin.Context) {

		var page_form PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var affair_list_query []model.Affair
		department := page_form.Department
		theme := page_form.Theme

		var total int64

		global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("department = ? AND theme = ?", department, theme).Find(&affair_list_query).Count(&total)

		var page_resp model.PageResp
		page_resp.Records = affair_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 政务管理分页
// @Tags Affair
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "政务名称"
// @Router /affairPage [get]
func affairPage(r *gin.Engine) {
	r.GET("/affairPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_401, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		// var affair_list_query []model.Affair
		name := page_form.Name

		var total int64
		// global.DB.Model(&affair_list_query).Count(&total)
		var affair_res []model.AffairRes

		if name != "" {
			// global.DB.Debug().Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Raw("select affair.* , user.nickname as audit from affair left join user on affair.audit_id = user.id where name like ?", name).Find(&affair_list_query).Count(&total)
			// global.DB.Debug().Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Model(&model.Affair{}).Select("affair.id as id,affair.name as name, affair.department as department, affair.theme as theme, affair.cost as cost, affair.state as state, affair.description as description, affair.audit_id as audit_id, user.nickname as audit").Joins("left join user on affair.audit_id = user.id where name like ?", name).Scan(&affair_res).Count(&total)
			global.DB.Debug().Limit(page_form.PageSize).
				Offset((page_form.PageNum-1)*page_form.PageSize).
				Where("name like ?", name).Model(&model.Affair{}).
				Select("affair.id as id,affair.name as name, affair.department as department, affair.theme as theme, affair.cost as cost, affair.state as state, affair.description as description, affair.audit_id as audit_id, user.nickname as audit").Joins("left join user on affair.audit_id = user.id").
				Scan(&affair_res).Count(&total)
		} else {
			// global.DB.Debug().Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Raw("select affair.* , user.nickname as audit from affair left join user on affair.audit_id = user.id").Find(&affair_list_query).Count(&total)
			// global.DB.Debug().Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Model(&model.Affair{}).Select("affair.id,affair.name, affair.department, affair.theme, affair.cost, affair.state, affair.description, affair.audit_id,user.nickname as audit").Joins("left join user on affair.audit_id = user.id").Scan(&affair_list_query).Count(&total)
			global.DB.Debug().Limit(page_form.PageSize).
				Offset((page_form.PageNum - 1) * page_form.PageSize).
				Model(&model.Affair{}).Select("affair.id as id,affair.name as name, affair.department as department, affair.theme as theme, affair.cost as cost, affair.state as state, affair.description as description, affair.audit_id as audit_id, user.nickname as audit").Joins("left join user on affair.audit_id = user.id").
				Scan(&affair_res).Count(&total)
		}

		var page_resp model.PageResp
		// page_resp.Records = affair_list_query
		page_resp.Records = affair_res
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 获取政务列表
// @Tags Affair
// @version 1.0
// @Accept json
// @Router /pullAffairList [get]
func pullAffairList(r *gin.Engine) {
	r.GET("/pullAffairList", func(c *gin.Context) {
		var affair_list_query []model.Affair
		global.DB.Find(&affair_list_query)

		middleware.Success(c, model.CODE_200, "", affair_list_query)
	})
}

// @Summary 根据政务编号获取信息
// @Tags Affair
// @version 1.0
// @Accept json
// @Param id path int true "政务编号"
// @Router /pullAffairById/{id} [get]
func pullAffairById(r *gin.Engine) {
	r.GET("/pullAffairById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var affair_query model.Affair
		global.DB.Where("id = ?", id).First(&affair_query)

		if affair_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", affair_query)
		} else {
			// middleware.Err(c, model.CODE_500, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 获取审批人员负责的政务
// @Tags Affair
// @version 1.0
// @Accept json
// @Param id path int true "审批人员编号"
// @Router /pullAffairByAuditId/{id} [get]
func pullAffairByAuditId(r *gin.Engine) {
	r.GET("/pullAffairByAuditId/:id", func(c *gin.Context) {
		id := c.Param("id")
		var affair_query []model.Affair
		global.DB.Where("audit_id = ?", id).Find(&affair_query)

		middleware.Success(c, model.CODE_200, "", affair_query)
	})
}

// @Summary 政务内容保存
// @Tags Affair
// @version 1.0
// @Accept json
// @Param affair body model.Affair true "政务内容"
// @Router /saveAffair [post]
func saveAffair(r *gin.Engine) {
	r.POST("/saveAffair", middleware.JWTAuth(), func(c *gin.Context) {
		var edited_affair model.Affair
		if err := c.ShouldBindJSON(&edited_affair); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&edited_affair)
		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除政务（单选）
// @Tags Affair
// @version 1.0
// @Accept json
// @Param id path int true "政务编号"
// @Router /deleteAffairById/{id} [delete]
func deleteAffairById(r *gin.Engine) {
	r.DELETE("/deleteAffairById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.Affair{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除政务（批量）
// @Tags Affair
// @version 1.0
// @Accept json
// @Param id body []int true "政务编号数组"
// @Router /delAffairBatch/ [delete]
func delAffairBatch(r *gin.Engine) {
	r.POST("/delAffairBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Affair{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

// @Summary 政务内容更新
// @Tags Affair
// @version 1.0
// @Accept json
// @Param affair body model.Affair true "政务内容"
// @Router /updateAffair [post]
func updateAffair(r *gin.Engine) {
	r.POST("/updateAffair", middleware.JWTAuth(), func(c *gin.Context) {
		var edited_affair model.Affair
		if err := c.ShouldBindJSON(&edited_affair); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Model(&model.Affair{}).Updates(model.Affair{
			Name:        edited_affair.Name,
			Department:  edited_affair.Department,
			Theme:       edited_affair.Theme,
			Cost:        edited_affair.Cost,
			State:       edited_affair.State,
			Description: edited_affair.Description,
			AuditId:     edited_affair.AuditId,
		})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}
