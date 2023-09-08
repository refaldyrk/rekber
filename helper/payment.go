package helper

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/veritrans/go-midtrans"
	"rekber/model"
	"strconv"
	"time"
)

func generatePaymentID() string {
	unixTime := strconv.Itoa(int(time.Now().Unix()))
	id := fmt.Sprintf("#%s|%d-%d[%d]", unixTime, time.Now().Second(), time.Now().Minute(), time.Now().Hour())
	return id
}

func GenerateLinkPayment(user model.User, amount int) (string, string, error) {
	midtransClient := midtrans.NewClient()
	midtransClient.ClientKey = viper.GetString("CLIENT_KEY")
	midtransClient.ServerKey = viper.GetString("SERVER_KEY")

	midtransClient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		midtransClient,
	}

	snapReq := midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Username,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  generatePaymentID(),
			GrossAmt: int64(amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(&snapReq)
	if err != nil {
		return "", "", err
	}

	return snapTokenResp.RedirectURL, snapReq.TransactionDetails.OrderID, nil
}
