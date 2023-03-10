package handlers

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/helper"
	"bcc-project-v/src/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func (h *handler) OnlinePayment(c *gin.Context) {
	//GetCart
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	cart, err := h.Repository.GetCartForPayment(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userPayment, err := h.Repository.FindUserByID(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed to get user", nil)
		return
	}

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_MIDTRANS_KEY"), midtrans.Sandbox)

	// 2. Initiate Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-PAY-EM-" + time.Now().UTC().Format("2006010215040105"),
			GrossAmt: int64(cart.TotalPrice),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: userPayment.FName,
			Email: userPayment.Email,
			Phone: userPayment.Contact,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.CreateTransaction(snapReq)

	//create payment for database
	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 1); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", nil)
		return
	}
	payment := entities.Payment{}
	payment.TotalPrice = cart.TotalPrice
	payment.UserID = user.ID
	payment.CartID = cart.ID
	payment.Type = "Online Payment"
	payment.StatusID = 1
	payment.PaymentCode = snapResp.Token
	payment.Status = status

	if err := h.Repository.CreatePayment(&payment); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed order", nil)
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
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get cart", nil)
	}

	// Membuat UUID secara acak
	id := uuid.New()

	// Mengonversi UUID ke string
	uniqueCode := id.String()
	status := entities.Status{}
	if err := h.Repository.FindStatus(&status, 1); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed get status", nil)
		return
	}

	payment := entities.Payment{}
	payment.TotalPrice = cart.TotalPrice
	payment.PaymentCode = uniqueCode
	payment.UserID = user.ID
	payment.CartID = cart.ID
	payment.Type = "Offline Payment"
	payment.StatusID = 1
	payment.Status = status

	if err := h.Repository.CreatePayment(&payment); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed order", nil)
		payment.StatusID = 3
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Selamat pemesanan anda telah berhasil dilakukan!!!", &payment)
}
