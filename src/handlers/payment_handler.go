package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	midtrans "github.com/veritrans/go-midtrans"
)

func (h *handler) OnlinePayment(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)

	cart, err := h.Repository.GetCartForPayment(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var totalPayment float64
	for _, p := range cart.CartProducts {
		fmt.Println(p.Quantity)
		product, _ := h.Repository.GetProductByID(p.ProductID)
		fmt.Println(product.Price)
		totalPayment += (float64(p.Quantity) * product.Price)
	}

	userFound, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get user", nil)
		return
	}

	// 1. Initiate Snap client
	midclient := midtrans.NewClient()
	midclient.ServerKey = h.config.GetEnv("SERVER_KEY")
	midclient.ClientKey = h.config.GetEnv("CLIENT_KEY")
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{}
	snapGateway = midtrans.SnapGateway{
		Client: midclient,
	}

	// 2. Initiate Snap request
	custAddress := &midtrans.CustAddress{
		FName:   userFound.FName,
		Phone:   userFound.Contact,
		Address: userFound.Address,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-PAY-EM-" + time.Now().UTC().Format("2006010215040105"),
			GrossAmt: int64(totalPayment),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName:    userFound.FName,
			Email:    userFound.Email,
			Phone:    userFound.Contact,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment token", err.Error())
		return
	}

	//create payment for database
	dataBuyer := model.DataPayment{}
	if err := c.ShouldBindJSON(&dataBuyer); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "The data you entered is in an invalid format. Please check and try again!", err.Error())
		return
	}

	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 1); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", err.Error())
		return
	}

	payment := entities.Payment{
		TotalPrice:  totalPayment,
		UserID:      user.ID,
		CartID:      cart.ID,
		Type:        "Online Payment",
		StatusID:    1,
		PaymentCode: snapResp.Token,
		Status:      status,
		FName:       dataBuyer.FName,
		Address:     dataBuyer.Address,
		Contact:     dataBuyer.Contact,
		City:        dataBuyer.City,
		Email:       dataBuyer.Email,
	}

	if err := h.Repository.CreatePayment(&payment); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to create order", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": snapResp.Token,
		"URL":   snapResp.RedirectURL,
	})

	helper.SuccessResponse(c, http.StatusOK, "Selamat pemesanan anda telah berhasil dilakukan!!!", &payment)
}

func (h *handler) OfflinePayment(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	cart, err := h.Repository.GetCartForPayment(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get cart", err.Error())
	}

	// Membuat UUID secara acak
	id := uuid.New()

	// Mengonversi UUID ke string
	uniqueCode := id.String()
	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 1); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", err.Error())
		return
	}

	payment := entities.Payment{}
	// payment.TotalPrice = cart.TotalPrice
	payment.PaymentCode = uniqueCode
	payment.UserID = user.ID
	payment.CartID = cart.ID
	payment.Type = "Offline Payment"
	payment.StatusID = 1
	payment.Status = status

	if err := h.Repository.CreatePayment(&payment); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed order", err.Error())
		payment.StatusID = 3
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Selamat pemesanan anda telah berhasil dilakukan!!!", &payment)
}
