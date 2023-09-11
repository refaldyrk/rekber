package dto

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginV2Req struct {
	Email string `json:"email"`
}
