package permissions

const (
	SubscriptionEarlySupporter = UploadPersonalBackground
)

func init() {
	PermissionsMap["subscriptions/earlySupporter"] = SubscriptionEarlySupporter
}
