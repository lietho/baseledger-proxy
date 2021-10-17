package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	uuid "github.com/kthomas/go.uuid"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/models"
	"github.com/unibrightio/proxy-api/restutil"
)

type createTxDto struct {
	Payload string `json:"payload"`
}

func CreateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		user := &models.User{}
		err = json.Unmarshal(buf, &user)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		if !user.Create() {
			logger.Error("error when creating new user")
			restutil.RenderError("error when creating new user", 500, c)
			return
		}

		restutil.Render(user.Email, 200, c)
	}
}

func LoginUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		user := &models.User{}
		err = json.Unmarshal(buf, &user)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		token, err := user.Login()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		restutil.Render(token, 200, c)
	}
}

func CreateTransactionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !restutil.HasEnoughBalance() {
			restutil.RenderError("not enough token balance", 400, c)
			return
		}

		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &createTxDto{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		// TODO: set proper size
		maximumPayloadSize := 20
		if len(req.Payload) > maximumPayloadSize {
			restutil.RenderError("payload maximum size exceeded", 400, c)
			return
		}

		transactionId := uuid.NewV4()
		signAndBroadcastPayload := restutil.SignAndBroadcastPayload{
			TransactionId: transactionId.String(),
			Payload:       req.Payload,
			OpCode:        9,
		}

		txHash := restutil.SignAndBroadcast(signAndBroadcastPayload)

		logger.Infof("Transaction hash of custom payload %v\n", txHash)
		restutil.Render(txHash, 200, c)
	}
}
