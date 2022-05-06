package person

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

/*
func personPage(r *gin.Engine) {
	r.GET("/personPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var person_list_query []model.Person
		name := page_form.Name

		var total int64
		// global.DB.Model(&person_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&person_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&person_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = person_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

func pullPersonList(r *gin.Engine) {
	r.GET("/pullPersonList", func(c *gin.Context) {
		var person_list_query []model.Person
		global.DB.Find(&person_list_query)

		middleware.Success(c, model.CODE_200, "", person_list_query)
	})
}

*/

// @Summary 根据用户编号获取实名信息
// @Tags Person
// @version 1.0
// @Accept json
// @Param id path int true "用户编号"
// @Router /pullPersonById/{id} [get]
func pullPersonById(r *gin.Engine) {
	r.GET("/pullPersonById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var person_query model.Person
		global.DB.Where("id = ?", id).First(&person_query)

		if person_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", person_query)
		} else {
			middleware.Err(c, model.CODE_600, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 保存实名信息
// @Tags Person
// @version 1.0
// @Accept json
// @Param person body model.Person true "保存实名信息"
// @Router /savePerson [post]
func savePerson(r *gin.Engine) {
	r.POST("/savePerson", middleware.JWTAuth(), func(c *gin.Context) {
		var inp_person model.Person
		if err := c.ShouldBindJSON(&inp_person); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&inp_person)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

/*
func deletePersonById(r *gin.Engine) {
	r.DELETE("/deletePersonById/:id", func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.Person{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

func delPersonBatch(r *gin.Engine) {
	r.POST("/delPersonBatch", func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Person{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}
*/

// @Summary 根据身份证号获取真实信息
// @Tags Person
// @version 1.0
// @Accept json
// @Param idt path string true "身份证号"
// @Router /pullPersonByIdentity/{idt} [get]
func pullPersonByIdentity(r *gin.Engine) {
	r.GET("/pullPersonByIdentity/:idt", func(c *gin.Context) {
		idt := c.Param("idt")
		var person model.Person
		global.DB.Debug().Where("identity = ?", idt).First(&person)

		middleware.Success(c, model.CODE_200, "", person)
	})
}
