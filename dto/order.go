package dto

type NewOrderReq struct {
	SellerIdentity string `json:"seller_identity" bson:"seller_identity"`
	BuyerIdentity  string `json:"buyer_identity" bson:"buyer_identity"`
	Type           string `json:"type" bson:"type"`
	Amount         int64  `json:"amount" bson:"amount" `
}
