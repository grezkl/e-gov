package affair

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	affairPage(e)
	pullAffairList(e)
	pullAffairById(e)
	pullAffairByAuditId(e)
	pullAffairPage(e)
	saveAffair(e)
	deleteAffairById(e)
	delAffairBatch(e)
}
