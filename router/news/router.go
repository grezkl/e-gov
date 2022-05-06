package news

import (
	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	newsPage(e)
	pullNewsList(e)
	pullNewsRecent(e)
	pullNewsById(e)
	updateNews(e)
	deleteNewsById(e)
	delNewsBatch(e)
}
