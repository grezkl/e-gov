package apply

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

// @Summary 政务申请分页
// @Tags Apply
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "政务申请名称"
// @Router /applyPage [get]
func applyPage(r *gin.Engine) {
	r.GET("/applyPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var apply_list_query []model.Apply
		name := page_form.Name

		var total int64
		// global.DB.Model(&apply_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&apply_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&apply_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = apply_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

/*
type PageForm struct {
	PageNum  int    `form:"pageNum" json:"pageNum" binding:"required"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"required"`
	Name     string `form:"name" json:"name"`
	UserId   int    `form:"userId" json:"userId"`
}

type ApplyPageResp struct {
	Records []model.User `form:"records" json:"records"`
	Total   int64        `form:"total" json:"total"`
}
*/

// @Summary 获取我的政务分页
// @Tags Apply
// @version 1.0
// @Accept json
// @Param page_form body model.PageForm true "页面规格及关键词"
// @Router /pullMyApply [get]
func pullMyApply(r *gin.Engine) {
	r.GET("/pullMyApply", middleware.JWTAuth(), func(c *gin.Context) {
		var page_form model.PageForm

		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}
		var apply_list_query []model.Apply
		name := page_form.Name

		userId, _ := c.Get("userId")
		user_id := uint(userId.(uint))

		var total int64
		// global.DB.Model(&apply_list_query).Count(&total)

		if name != "" {
			global.DB.Debug().Order("id desc").Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("user_id = ? and name like ?", user_id, name).Find(&apply_list_query).Count(&total)
		} else {
			global.DB.Debug().Order("id desc").Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("user_id = ?", user_id).Find(&apply_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = apply_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)

	})
}

// @Summary 获取政务申请列表
// @Tags Apply
// @version 1.0
// @Accept json
// @Router /pullApplyList [get]
func pullApplyList(r *gin.Engine) {
	r.GET("/pullApplyList", func(c *gin.Context) {
		var apply_list_query []model.Apply
		global.DB.Find(&apply_list_query)

		middleware.Success(c, model.CODE_200, "", apply_list_query)
	})
}

// @Summary 根据编号查找政务申请
// @Tags Account
// @version 1.0
// @Accept json
// @Param id path int true "政务申请编号"
// @Router /pullUserById/{id} [get]
func pullApplyById(r *gin.Engine) {
	r.GET("/pullApplyById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var apply_query model.Apply
		global.DB.Where("id = ?", id).First(&apply_query)

		if apply_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", apply_query)
		} else {
			middleware.Err(c, model.CODE_600, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 政务申请办理
// @Tags Apply
// @version 1.0
// @Accept json
// @Param apply body model.Apply true "政务申请"
// @Router /saveApply [post]
func saveApply(r *gin.Engine) {
	r.POST("/saveApply", middleware.JWTAuth(), func(c *gin.Context) {
		var inp_apply model.Apply
		if err := c.ShouldBindJSON(&inp_apply); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&inp_apply)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 撤销政务申请（单个）
// @Tags Apply
// @version 1.0
// @Accept json
// @Param id path int true "政务申请编号"
// @Router /deleteApplyById/{id} [delete]
func deleteApplyById(r *gin.Engine) {
	r.DELETE("/deleteApplyById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")

		var total int64
		global.DB.Where("id = ?", id).Delete(&model.Apply{}).Count(&total)

		if total == 0 {
			middleware.Err(c, model.CODE_600, "删除失败", nil)
		} else {
			middleware.Success(c, model.CODE_200, "删除成功", nil)
		}
	})
}

// @Summary 撤销政务申请（批量）
// @Tags Apply
// @version 1.0
// @Accept json
// @Param ids path []int true "政务申请编号数组"
// @Router /delApplyBatch [delete]
func delApplyBatch(r *gin.Engine) {
	r.DELETE("/delApplyBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Apply{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}
