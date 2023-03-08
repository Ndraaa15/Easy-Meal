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
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func (h *handler) OnlinePayment() {
	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New("YOUR-SERVER-KEY", midtrans.Sandbox)

	// 2. Initiate Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-GO-ID-" + time.Now().UTC().Format("2006010215040105"),
			GrossAmt: 200000,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "John",
			LName: "Doe",
			Email: "john@doe.com",
			Phone: "081234567890",
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Someitem",
			},
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.CreateTransaction(snapReq)
	fmt.Println("Response :", snapResp)
}

func (h *handler) OfflinePayment(c *gin.Context) {
	userClaims, _ := c.Get("user")
	user := userClaims.(model.UserClaims)
	cart, err := h.Repository.GetCart(user.ID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Faied get cart", nil)
		return
	}

	// Membuat UUID secara acak
	id := uuid.New()

	// Mengonversi UUID ke string
	uniqueCode := id.String()

	payment := entities.OfflinePayment{}
	payment.TotalPrice = cart.TotalPrice
	payment.PaymentCode = uniqueCode
	payment.UserID = user.ID

	if err := h.Repository.CreatePayment(&payment); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Failed order", nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Selamat pemesanan anda telah berhasil dilakukan!!!", &payment)
}
