package account

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	login(e)
	register(e)
	updatePWD(e)
	username(e)
	userPage(e)
	pullUserList(e)
	pullUserById(e)
	pullUserByRole(e)
	saveUser(e)
	deleteUserById(e)
	delUserBatch(e)
	echartsMember(e)
	userExport(e)
	userImport(e)
}
