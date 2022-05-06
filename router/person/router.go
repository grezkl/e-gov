package person

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	// personPage(e)
	// pullPersonList(e)
    pullPersonById(e)
	pullPersonByIdentity(e)
	savePerson(e)
	// deletePersonById(e)
	// delPersonBatch(e)
}
