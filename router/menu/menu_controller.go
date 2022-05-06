package menu

import (
	"e-gov/global"
	"e-gov/middleware"
	"e-gov/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 菜单管理分页
// @Tags Menu
// @version 1.0
// @Accept json
// @Param pageNum query int true "页数"
// @Param pageSize query int true "页大小"
// @Param name query string false "政务申请名称"
// @Router /menuPage [get]
func menuPage(r *gin.Engine) {
	r.GET("/menuPage", func(c *gin.Context) {

		var page_form model.PageForm
		if c.ShouldBindQuery(&page_form) != nil {
			middleware.Err(c, model.CODE_500, "参数错误", nil)
			return
		}
		if page_form.PageNum <= 0 {
			page_form.PageNum = 1
		}

		var menu_list_query []model.Menu
		name := page_form.Name

		var total int64
		// global.DB.Model(&menu_list_query).Count(&total)

		if name != "" {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum-1)*page_form.PageSize).Where("name = ?", name).Find(&menu_list_query).Count(&total)
		} else {
			global.DB.Limit(page_form.PageSize).Offset((page_form.PageNum - 1) * page_form.PageSize).Find(&menu_list_query).Count(&total)
		}

		var page_resp model.PageResp
		page_resp.Records = menu_list_query
		page_resp.Total = total

		middleware.Success(c, model.CODE_200, "", page_resp)
	})
}

// @Summary 获取菜单列表
// @Tags Menu
// @version 1.0
// @Accept json
// @Router /pullMenuList [get]
func pullMenuList(r *gin.Engine) {
	r.GET("/pullMenuList", func(c *gin.Context) {
		var menu_list_query []model.Menu
		global.DB.Find(&menu_list_query)

		var menu_list []model.Menu
		var menu_map map[uint]int // map[menu's id]menu's index
		menu_map = make(map[uint]int)
		i := 0

		for count := 0; count < len(menu_list_query); count++ {
			cur_pid := menu_list_query[count].Pid

			if cur_pid == "" {
				menu_list = append(menu_list, menu_list_query[count])
				menu_map[menu_list_query[count].ID] = i
				i++
			} else {
				cur_pid_num, _ := strconv.Atoi(cur_pid)
				u_cur_pid := uint(cur_pid_num)
				last := menu_map[u_cur_pid]

				menu_list[last].Children = append(menu_list[last].Children, menu_list_query[count])
			}
		}

		middleware.Success(c, model.CODE_200, "", menu_list)
	})
}

func pullMenuById(r *gin.Engine) {
	r.GET("/pullMenuById/:id", func(c *gin.Context) {
		id := c.Param("id")
		var menu_query model.Menu
		global.DB.Where("id = ?", id).First(&menu_query)

		middleware.Success(c, model.CODE_200, "", menu_query)
	})
}

// @Summary 保存菜单内容
// @Tags Menu
// @version 1.0
// @Accept json
// @Param menu body model.Menu true "菜单内容"
// @Router /saveMenu [post]
func saveMenu(r *gin.Engine) {
	r.POST("/saveMenu", middleware.JWTAuth(), func(c *gin.Context) {
		var menu model.Menu
		if err := c.ShouldBindJSON(&menu); err != nil {
			middleware.Err(c, model.CODE_401, "请求出现异常，请稍后重试", nil)
			return
		}

		global.DB.Save(&menu)

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除菜单（单个）
// @Tags Menu
// @version 1.0
// @Accept json
// @Param id path int true "菜单编号"
// @Router /deleteMenuById/{id} [delete]
func deleteMenuById(r *gin.Engine) {
	r.DELETE("/deleteMenuById/:id", middleware.JWTAuth(), func(c *gin.Context) {
		id := c.Param("id")
		global.DB.Where("id = ?", id).Delete(&model.Menu{})

		middleware.Success(c, model.CODE_200, "", nil)
	})
}

// @Summary 删除菜单（批量）
// @Tags Menu
// @version 1.0
// @Accept json
// @Param id path int true "菜单编号数组"
// @Router /delMenuBatch/ [delete]
func delMenuBatch(r *gin.Engine) {
	r.DELETE("/delMenuBatch", middleware.JWTAuth(), func(c *gin.Context) {
		var ids []int
		c.ShouldBindJSON(&ids)
		global.DB.Where("id in ?", ids).Delete(&model.Menu{})

		middleware.Success(c, model.CODE_200, "", ids)
	})
}

// @Summary 获取菜单图标
// @Tags Menu
// @version 1.0
// @Accept json
// @Router /getMenuIcons/ [delete]
func getMenuIcons(r *gin.Engine) {
	r.GET("/getMenuIcons", func(c *gin.Context) {
		var icon_list_query []model.Dict
		global.DB.Find(&icon_list_query)

		middleware.Success(c, model.CODE_200, "", icon_list_query)
	})
}
