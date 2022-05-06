package role

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"

	"github.com/gin-gonic/gin"
)

// @Summary 角色管理分页
// @Tags Role
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "角色名称"
// @Router /rolePage [get]
func rolePage(r *gin.Engine) {
	r.GET("/rolePage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var role_list_query []model.Role
		name := page_form.Name

		var total int64
		// global.DB.Model(&role_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name like ?", name).Find(&role_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&role_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = role_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 获取角色列表
// @Tags Role
// @version 1.0
// @Accept json
// @Router /pullRoleList [get]
func pullRoleList(r *gin.Engine) {
	r.GET("/pullRoleList", func(c *gin.Context) {
		var role_list_query []model.Role
		global.DB.Find(&role_list_query)

		middleware.Success(c, model.CODE_200, "", role_list_query)
	})
}

func pullRoleById(r *gin.Engine) {
	r.GET("/pullRoleById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var role_query model.Role
		global.DB.Where("id = ?", id).First(&role_query)

		if role_query.ID != 0 {
			middleware.Success(c, model.CODE_200, "", role_query)
		} else {
			// middleware.Err(c, model.CODE_500, "该 id 查询无结果", nil)
		}
	})
}

// @Summary 保存角色内容
// @Tags Role
// @version 1.0
// @Accept json
// @Param role body model.Role true "角色内容"
// @Router /saveRole [post]
func saveRole(r *gin.Engine) {
	r.POST("/saveRole", middleware.JWTAuth(), func(c *gin.Context) {
		var role model.Role
		if err := c.ShouldBindJSON(&role); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&role)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除角色（单个）
// @Tags Role
// @version 1.0
// @Accept json
// @Param id path int true "角色编号"
// @Router /deleteRoleById/{id} [delete]
func deleteRoleById(r *gin.Engine) {
	r.DELETE("/deleteRoleById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.Role{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除角色（批量）
// @Tags Role
// @version 1.0
// @Accept json
// @Param id path []int true "角色编号数组"
// @Router /delRoleBatch/ [delete]
func delRoleBatch(r *gin.Engine) {
	r.DELETE("/delRoleBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Role{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

func roleMenu(r *gin.Engine) {
	r.POST("/roleMenu/:roleId", func(c *gin.Context) {

	})
}

// @Summary 获取角色菜单
// @Tags Role
// @version 1.0
// @Accept json
// @Param id path int true "角色编号"
// @Router /getRoleMenu/{roleId} [get]
func getRoleMenu(r *gin.Engine) {
	r.GET("/getRoleMenu/:roleId", func(c *gin.Context) {
		roleId := c.Param("roleId")
		var role_menu_query []model.RoleMenu
		// var menu_query model.Menu
		global.DB.Where("role_id = ?", roleId).Find(&role_menu_query)

		var list []uint
		for i := 0; i < len(role_menu_query); i++ {
			list = append(list, role_menu_query[i].MenuId)
		}

		middleware.Success(c, model.CODE_200, "", list)
	})

}

func pullMenuIds(r *gin.Engine) {
	r.GET("/pullMenuIds", func(c *gin.Context) {
		var menu_query []model.Menu
		global.DB.Find(&menu_query)

		var list []uint
		for i := 0; i < len(list); i++ {
			list = append(list, menu_query[i].ID)
		}
		middleware.Success(c, model.CODE_200, "", list)
	})
}
