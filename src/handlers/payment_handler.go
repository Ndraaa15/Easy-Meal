package handlers

import (
	"fmt"
	"time"

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

func (h *handler) OfflinePayment() {

}
