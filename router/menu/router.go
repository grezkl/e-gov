package menu

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	menuPage(e)
	pullMenuList(e)
	pullMenuById(e)
	saveMenu(e)
	deleteMenuById(e)
	delMenuBatch(e)
	getMenuIcons(e)
}
