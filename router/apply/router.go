package apply

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	applyPage(e)
	pullApplyList(e)
	pullApplyById(e)
	pullMyApply(e)
	saveApply(e)
	deleteApplyById(e)
	delApplyBatch(e)
}
