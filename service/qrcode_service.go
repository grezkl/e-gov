package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type QRCodeService struct {
}

var ctx = context.Background()

func (qrs *QRCodeService) GenerateQRCode(c *gin.Context) *api.JsonResult {
	mTicket := uuid.NewV4().String()
	qrCodeUrl := fmt.Sprintf("http://192.168.1.100:7000/api/v1/qrcode/sanlogin?m_ticket=%s", mTicket)

	// result, err := _redisHelper.SetNX(ctx, constant.MOOSE_TICKET, mTicket, 3*time.Minute).Result()
	redisHelper := engine.GetRedisEngine()
	ticketKey := fmt.Sprintf(constant.MOOSE_SCAN_TICKET, mTicket)
	authInfo := &model.AuthInfo{Token: "", Status: 0}
	// 设置过期时间，三分钟
	_, err := redisHelper.Set(ctx, ticketKey, authInfo, 3*time.Minute).Result()
	if err != nil {
		log.Println(err)
		return api.JsonError(api.QRCodeGetFailErr)
	}

	c.SetCookie("m_ticket", mTicket, 3*60, "/", "localhost", true, true)

	qrCodeInfo := &model.QRCodeInfo{
		CodeUrl: qrCodeUrl,
	}
	return api.JsonData(qrCodeInfo)
}

func (qrs *QRCodeService) AskQRCode(c *gin.Context) *api.JsonResult {
	mTicket, err := c.Cookie("m_ticket")
	log.Printf("the m_ticket %s", mTicket)

	if err != nil {
		log.Println(err)
		return api.JsonError(api.QRCodeRetryErr)
	}

	redisHelper := engine.GetRedisEngine()
	ticketKey := fmt.Sprintf(constant.MOOSE_SCAN_TICKET, mTicket)
	result, err := redisHelper.Get(ctx, ticketKey).Result()

	if err != nil {
		return api.JsonError(api.QRCodeRetryErr)
	}

	// 检查 ticket
	if len(result) <= 0 || result == "" {
		return api.JsonError(api.QRCodeRetryErr)
	}

	var authInfo model.AuthInfo
	json.Unmarshal([]byte(result), &authInfo)

	// 扫描 token 不为空，返回之后，设置过期，在 redis Expire
	if authInfo.Token != "" {
		redisHelper.Expire(ctx, ticketKey, 0)
	}

	log.Printf("ticketKey %s result %s ", ticketKey, result)
	return api.JsonData(authInfo)
}

// 需要登录，从 app 调用
//   - app 必须是登录状态
//   - app 扫描之后，携带 token 校验是否合法
// 判断 m_ticket 是否存在
// 校验 token 是否合法
func (qrs *QRCodeService) ScanLogin(c *gin.Context) *api.JsonResult {
	mTicket := c.Query("m_ticket")
	// 校验 m_ticket
	if mTicket == "" {
		return api.JsonError(api.QRCodeRetryErr)
	}

	redisHelper := engine.GetRedisEngine()
	ticketKey := fmt.Sprintf(constant.MOOSE_SCAN_TICKET, mTicket)
	result, err := redisHelper.Get(ctx, ticketKey).Result()

	if err != nil {
		log.Println(err)
		return api.JsonError(api.QRCodeRetryErr)
	}

	// 检查 ticket
	if len(result) <= 0 || result == "" {
		return api.JsonError(api.QRCodeRetryErr)
	}

	var authInfo model.AuthInfo
	err = json.Unmarshal([]byte(result), &authInfo)
	if err != nil {
		log.Println(err)
		return api.JsonError(api.QRCodeRetryErr)
	}

	authInfo.Status = 1

	mToken := c.Query("token")
	// if token valid success
	if mToken != "" {
		authInfo.Token = mToken
	}

	_, err = redisHelper.Set(ctx, ticketKey, &authInfo, 3*time.Minute).Result()
	if err != nil {
		log.Println(err)
		return api.JsonError(api.QRCodeGetFailErr)
	}
	// 设置过期时间，三分钟
	// redisHelper.Expire(ctx, ticketKey, 3*time.Minute).Result()
	log.Printf("key %s result %s", ticketKey, result)
	return api.JsonData(authInfo)
}
