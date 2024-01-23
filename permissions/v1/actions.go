package permissions

// Basic user actions
const (
	UseCommands Permissions = 10 << iota

	// Content
	SelectPersonalBackgroundPreset Permissions = 11 << iota
	UploadPersonalBackground

	// Connections
	CreatePersonalConnection Permissions = 11 << iota
	UpdatePersonalConnection
	RemovePersonalConnection

	// Subscriptions
	RetrievePersonalSubscriptions Permissions = 12 << iota
	CreatePersonalSubscription
	ExtendPersonalSubscription
)

// Moderation actions
const (
	// Background Presets
	UploadBackgroundPreset Permissions = 15 << iota
	RemoveBackgroundPreset

	// Manage User Content
	UpdateUserBackground Permissions = 16 << iota
	RemoveUserBackground

	// Subscriptions
	RetrieveUserSubscriptions Permissions = 20 << iota
	CreateUserSubscription
	ExtendUserSubscription
	TerminateUserSubscription

	// Connections
	RetrieveUserConnections Permissions = 21 << iota
	ManageUserConnectionVerification
	RemoveUserConnection

	// Restrictions
	RetrieveUserRestrictions Permissions = 22 << iota
	CreateUserRestriction
	RemoveUserRestriction
)

const (
	// Admin
	ManageUserRoles Permissions = 30 << iota
)
