package feedback

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	feedbackPage(e)
	pullFeedbackList(e)
	pullFeedbackByUserId(e)
	saveFeedback(e)
	deleteFeedbackById(e)
	deleteFeedbackByUserId(e)
	delFeedbackBatch(e)
}
