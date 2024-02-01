package permissions

// Basic user actions
const (
	UseCommands Permissions = BasicUserActions | 1<<(10+iota)

	// Content
	SelectPersonalBackgroundPreset Permissions = BasicUserActions | 1<<(11+iota)
	UploadPersonalBackground

	// Connections
	CreatePersonalConnection Permissions = BasicUserActions | 1<<(12+iota)
	UpdatePersonalConnection
	RemovePersonalConnection

	// Subscriptions
	RetrievePersonalSubscriptions Permissions = BasicUserActions | 1<<(13+iota)
	CreatePersonalSubscription
	ExtendPersonalSubscription
)

// Moderation actions
const (
	// Background Presets
	UploadBackgroundPreset Permissions = ModerationActions | 1<<(20+iota)
	RemoveBackgroundPreset

	// Manage User Content
	UpdateUserBackground Permissions = ModerationActions | 1<<(21+iota)
	RemoveUserBackground

	// Subscriptions
	RetrieveUserSubscriptions Permissions = ModerationActions | 1<<(22+iota)
	CreateUserSubscription
	ExtendUserSubscription
	TerminateUserSubscription

	// Connections
	RetrieveUserConnections Permissions = ModerationActions | 1<<(23+iota)
	ManageUserConnectionVerification
	RemoveUserConnection

	// Restrictions
	RetrieveUserRestrictions Permissions = ModerationActions | 1<<(24+iota)
	CreateUserRestriction
	RemoveUserRestriction
)

const (
	// Admin
	ManageUserRoles Permissions = AdminActions | 1<<(30+iota)
)

func init() {
	PermissionsMap["actions/useCommands"] = UseCommands

	PermissionsMap["actions/selectPersonalBackgroundPreset"] = SelectPersonalBackgroundPreset
	PermissionsMap["actions/uploadPersonalBackground"] = UploadPersonalBackground

	PermissionsMap["actions/createPersonalConnection"] = CreatePersonalConnection
	PermissionsMap["actions/updatePersonalConnection"] = UpdatePersonalConnection
	PermissionsMap["actions/removePersonalConnection"] = RemovePersonalConnection

	PermissionsMap["actions/retrievePersonalSubscriptions"] = RetrievePersonalSubscriptions
	PermissionsMap["actions/createPersonalSubscription"] = CreatePersonalSubscription
	PermissionsMap["actions/extendPersonalSubscription"] = ExtendPersonalSubscription

	PermissionsMap["actions/uploadBackgroundPreset"] = UploadBackgroundPreset
	PermissionsMap["actions/removeBackgroundPreset"] = RemoveBackgroundPreset

	PermissionsMap["actions/updateUserBackground"] = UpdateUserBackground
	PermissionsMap["actions/removeUserBackground"] = RemoveUserBackground

	PermissionsMap["actions/retrieveUserSubscriptions"] = RetrieveUserSubscriptions
	PermissionsMap["actions/createUserSubscription"] = CreateUserSubscription
	PermissionsMap["actions/extendUserSubscription"] = ExtendUserSubscription
	PermissionsMap["actions/terminateUserSubscription"] = TerminateUserSubscription

	PermissionsMap["actions/retrieveUserConnections"] = RetrieveUserConnections
	PermissionsMap["actions/manageUserConnectionVerification"] = ManageUserConnectionVerification
	PermissionsMap["actions/removeUserConnection"] = RemoveUserConnection

	PermissionsMap["actions/retrieveUserRestrictions"] = RetrieveUserRestrictions
	PermissionsMap["actions/createUserRestriction"] = CreateUserRestriction
	PermissionsMap["actions/removeUserRestriction"] = RemoveUserRestriction

	PermissionsMap["actions/manageUserRoles"] = ManageUserRoles
}
