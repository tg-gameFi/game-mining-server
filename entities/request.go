package entities

type ProxyGetParam struct {
	Url string `form:"url" binding:"required,url"`
}

type CoinPriceParam struct {
	CoinSymbols string `form:"coinSymbols" binding:"required,min=0,max=2000"` // a string of list: BTC,ETH,BNB
	FiatSymbol  string `form:"fiatSymbol" binding:"required,alpha,min=1,max=10"`
}

type UserLoginParam struct {
	InitDataRaw string `json:"initDataRaw" binding:"required"`
	Referral    string `json:"referral" binding:"omitempty,alphanum,min=1,max=20"`
	RandPoint   int64  `json:"randPoint" binding:"omitempty,number,min=150,max=300"`
}

type UserCheckinClaimParam struct {
	CheckinId string `json:"checkinId" binding:"required,uuid"`
}

type InvitedUserListParam struct {
	Offset int `form:"offset" binding:"number,min=0,max=10000"`
	Limit  int `form:"limit" binding:"required,number,min=0,max=100"`
}

type UserTaskClaimParam struct {
	TaskGroup string `json:"taskGroup" binding:"required,oneof=social wallet invite"`
	ClaimKey  string `json:"claimKey" binding:"required"` // for social claim, key is taskId, for wallet claim, key is transaction hash, for invite claim, key is level
}

type LeaderBoardParam struct {
	Offset int `form:"offset" binding:"number,min=0,max=10000"`
	Limit  int `form:"limit" binding:"required,number,min=0,max=100"`
}

type CreateMomentParam struct {
	UserId   int64  `json:"user_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url"`
}

type DeleteMomentParam struct {
	UserId   int64 `json:"user_id" binding:"required"`
	MomentId int64 `json:"moment_id" binding:"required"`
}

type GetLatestMomentsParam struct {
	Offset int `form:"offset" binding:"number,min=0,max=100"`
	Limit  int `form:"limit" binding:"required,number,min=0,max=20"`
}

type AddCommentParam struct {
	UserId        int64  `json:"user_id" binding:"required"`
	MomentId      int64  `json:"moment_id" binding:"required"`
	Content       string `json:"content" binding:"required"`
	ReplyToUserId *int64 `json:"reply_to_user_id"`
}

type DeleteCommentParam struct {
	UserId    int64 `json:"user_id" binding:"required"`
	CommentId int64 `json:"comment_id" binding:"required"`
}

type GetCommentsForMomentParam struct {
	MomentId int64 `form:"moment_id" binding:"required"`
	Offset   int   `form:"offset" binding:"number,min=0"`
	Limit    int   `form:"limit" binding:"required,number,min=0,max=20"`
}

type LikeMomentParam struct {
	UserId   int64 `json:"user_id" binding:"required"`
	MomentId int64 `json:"moment_id" binding:"required"`
}

type RewardMomentParam struct {
	MomentId   int64 `json:"moment_id" binding:"required"`
	FromUserId int64 `json:"from_user_id" binding:"required"`
	ToUserId   int64 `json:"to_user_id" binding:"required"`
	Amount     int   `json:"amount" binding:"required"`
}
