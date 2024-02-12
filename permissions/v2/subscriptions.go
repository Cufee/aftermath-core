package permissions

const (
	SubscriptionAftermathPlus = UploadPersonalBackground | UseLiveSessions
	SubscriptionAftermathPro  = SubscriptionAftermathPlus
)

func init() {
	PermissionsMap["subscriptions/aftermathPlus"] = SubscriptionAftermathPlus
}
