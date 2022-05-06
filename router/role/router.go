package role

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	rolePage(e)
	pullRoleList(e)
	pullRoleById(e)
	saveRole(e)
	deleteRoleById(e)
	delRoleBatch(e)
	getRoleMenu(e)
	pullMenuIds(e)
}
