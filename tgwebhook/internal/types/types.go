// Code generated by goctl. DO NOT EDIT.
package types

type TgSendReq struct {
	Chatid int64  `json:"chatid"`
	Token  string `json:"token"`
	Text   string `json:"text"`
}

type TgSendResp struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}
