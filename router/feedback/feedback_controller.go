package feedback

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

// @Summary 反馈分页
// @Tags Feedback
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "反馈名称"
// @Router /feedbackPage [get]
func feedbackPage(r *gin.Engine) {
	r.GET("/feedbackPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var feedback_list_query []model.Feedback
		name := page_form.Name

		var total int64
		// global.DB.Model(&feedback_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&feedback_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&feedback_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = feedback_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

func pullMyFeedback(r *gin.Engine) {
	r.GET("/pullMyFeedback", func(c *gin.Context) {

	})
}

// @Summary 获取反馈列表
// @Tags Feedback
// @version 1.0
// @Accept json
// @Router /pullFeedbackList [get]
func pullFeedbackList(r *gin.Engine) {
	r.GET("/pullFeedbackList", func(c *gin.Context) {
		var feedback_list_query []model.Feedback
		global.DB.Find(&feedback_list_query)

		middleware.Success(c, model.CODE_200, "", feedback_list_query)
	})
}

// @Summary 根据用户编号查询反馈
// @Tags Feedback
// @version 1.0
// @Accept json
// @Param id path int true "用户编号"
// @Router /pullFeedbackByUserId/{id} [get]
func pullFeedbackByUserId(r *gin.Engine) {
	r.GET("/pullFeedbackByUserId/:id", func(c *gin.Context) {
		id := c.Param("id")
		var feedback_query []model.Feedback
		global.DB.Where("user_id = ?", id).Find(&feedback_query)

		middleware.Success(c, model.CODE_200, "", feedback_query)
	})
}

// @Summary 反馈提交
// @Tags Feedback
// @version 1.0
// @Accept json
// @Param feedback body model.Feedback true "反馈内容"
// @Router /saveFeedback [post]
func saveFeedback(r *gin.Engine) {
	r.POST("/saveFeedback", middleware.JWTAuth(), func(c *gin.Context) {
		var commit_fb model.Feedback
		if err := c.ShouldBindJSON(&commit_fb); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&commit_fb)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 撤销反馈（单个）
// @Tags Feedback
// @version 1.0
// @Accept json
// @Param id path int true "反馈编号"
// @Router /deleteFeedbackById/{id} [delete]
func deleteFeedbackById(r *gin.Engine) {
	r.DELETE("/deleteFeedbackById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.Feedback{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 撤销反馈（批量）
// @Tags Feedback
// @version 1.0
// @Accept json
// @Param id path int true "反馈编号数组"
// @Router /delFeedbackBatch/ [delete]
func delFeedbackBatch(r *gin.Engine) {
	r.DELETE("/delFeedbackBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Feedback{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

func deleteFeedbackByUserId(r *gin.Engine) {
	r.DELETE("/deleteFeedbackByUserId/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("user_id = ?", id).Delete(&model.Feedback{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}
