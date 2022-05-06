package file

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	filePage(e)
	deleteFileById(e)
	delFileBatch(e)
	updateFile(e)
	uploadFile(e)
	downloadFile(e)
}
