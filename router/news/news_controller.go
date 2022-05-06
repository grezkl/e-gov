package news

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

// @Summary 新闻动态分页
// @Tags News
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "新闻名称"
// @Router /newsPage [get]
func newsPage(r *gin.Engine) {
	r.GET("/newsPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var news_list_query []model.News
		name := page_form.Name

		var total int64
		// global.DB.Model(&news_list_query).Count(&total)

		if name != "" {
			global.DB.Order("id desc").Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&news_list_query).Count(&total)
		} else {
			global.DB.Order("id desc").Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&news_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = news_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 获取新闻列表
// @Tags News
// @version 1.0
// @Accept json
// @Router /pullNewsList [get]
func pullNewsList(r *gin.Engine) {
	r.GET("/pullNewsList", func(c *gin.Context) {
		var news_list_query []model.News
		global.DB.Find(&news_list_query)

		middleware.Success(c, model.CODE_200, "", news_list_query)
	})
}

// @Summary 获取近期新闻
// @Tags News
// @version 1.0
// @Accept json
// @Router /pullNewsList [get]
func pullNewsRecent(r *gin.Engine) {
	r.GET("/pullNewsRecent", func(c *gin.Context) {
		var news_list_query []model.News
		global.DB.Limit(3).Order("id desc").Find(&news_list_query)

		middleware.Success(c, model.CODE_200, "", news_list_query)
	})

}

// @Summary 根据编号获取新闻
// @Tags News
// @version 1.0
// @Accept json
// @Param id path int true "新闻编号"
// @Router /pullNewsById/{id} [get]
func pullNewsById(r *gin.Engine) {
	r.GET("/pullNewsById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var news_query model.News
		global.DB.Where("id = ?", id).First(&news_query)

		if news_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", news_query)
		} else {
			middleware.Err(c, model.CODE_500, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 保存新闻内容
// @Tags News
// @version 1.0
// @Accept json
// @Param news body model.News true "新闻内容"
// @Router /updateNews [post]
func updateNews(r *gin.Engine) {
	r.POST("/saveNews", middleware.JWTAuth(), func(c *gin.Context) {
		var news model.News
		if err := c.ShouldBindJSON(&news); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&news)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除新闻（单个）
// @Tags News
// @version 1.0
// @Accept json
// @Param id path int true "新闻编号"
// @Router /deleteNewsById/{id} [delete]
func deleteNewsById(r *gin.Engine) {
	r.DELETE("/deleteNewsById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.News{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除新闻（批量）
// @Tags News
// @version 1.0
// @Accept json
// @Param id path []int true "新闻编号数组"
// @Router /delNewsBatch/ [delete]
func delNewsBatch(r *gin.Engine) {
	r.DELETE("/delNewsBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.News{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}
