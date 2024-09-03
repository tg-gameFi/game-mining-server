package configs

// Env constant
const (
	EnvDEV  = "dev"
	EnvTEST = "test"
	EnvPROD = "prod"
)

const (
	FiatUSD = "USD"
)

const (
	CurUser = "c_u" // key for user data in middleware
)

const (
	CheckinStatusUnclaimed = 0
	CheckinStatusClaimed   = 1

	CheckinBaseRewardPoint = int64(10)
)

const (
	TaskGroupSocial = "social"
	TaskGroupWallet = "wallet"
	TaskGroupInvite = "invite"

	TaskTypeSocialSubscribeTgChannel = "socialSubscribeTgChannel"
	TaskTypeSocialFollowCfOnX        = "socialFollowCfOnX"
	TaskTypeSocialRtAnn              = "socialRtAnn"
	TaskTypeWalletSendTx             = "walletSendTx"
	TaskTypeInviteFriends            = "inviteFriends"

	TaskStatusCreated   = 0
	TaskStatusClaimable = 1
	TaskStatusClaimed   = 2

	TaskSocialBaseRewardPoint = int64(10)
	TaskWalletBaseRewardPoint = int64(50)
)
