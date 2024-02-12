package permissions

const (
	// Basic user actions
	UseBasicCommands Permissions = 1 << (iota)
	UseLiveSessions

	// Content
	SelectPersonalBackgroundPreset Permissions = 1 << (5 + iota)
	UploadPersonalBackground

	// Connections
	CreatePersonalConnection Permissions = 1 << (10 + iota)
	UpdatePersonalConnection
	RemovePersonalConnection

	// Subscriptions
	RetrievePersonalSubscriptions Permissions = 1 << (15 + iota)
	CreatePersonalSubscription
	ExtendPersonalSubscription
)

// Moderation actions
const (
	// Background Presets
	UploadBackgroundPreset Permissions = 1 << (20 + iota)
	RemoveBackgroundPreset

	// Manage User Content
	UpdateUserBackground Permissions = 1 << (25 + iota)
	RemoveUserBackground

	// Subscriptions
	RetrieveUserSubscriptions Permissions = 1 << (30 + iota)
	CreateUserSubscription
	ExtendUserSubscription
	TerminateUserSubscription

	// Connections
	RetrieveUserConnections Permissions = 1 << (35 + iota)
	ManageUserConnectionVerification
	RemoveUserConnection

	// Restrictions
	RetrieveUserRestrictions Permissions = 1 << (40 + iota)
	CreateUserRestriction
	RemoveUserRestriction
)

const (
	// Admin
	ManageUserRoles Permissions = 1 << (45 + iota)
)

func init() {
	PermissionsMap["actions/useBasicCommands"] = UseBasicCommands
	PermissionsMap["actions/useLiveSessions"] = UseLiveSessions

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
