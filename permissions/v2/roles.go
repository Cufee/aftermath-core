package permissions

const (
	User = UseCommands | SelectPersonalBackgroundPreset | CreatePersonalConnection | UpdatePersonalConnection | RemovePersonalConnection | RetrievePersonalSubscriptions | CreatePersonalSubscription | ExtendPersonalSubscription

	ContentModerator = User | UpdateUserBackground | RemoveUserBackground | CreateUserRestriction
	GlobalModerator  = ContentModerator | RetrieveUserSubscriptions | CreateUserSubscription | ExtendUserSubscription | TerminateUserSubscription | UploadBackgroundPreset | RemoveBackgroundPreset | RetrieveUserConnections | ManageUserConnectionVerification | RemoveUserConnection | RetrieveUserRestrictions | RemoveUserRestriction

	Admin = GlobalModerator | ManageUserRoles
)

func init() {
	PermissionsMap["roles/user"] = User
	PermissionsMap["roles/contentModerator"] = ContentModerator
	PermissionsMap["roles/globalModerator"] = GlobalModerator
	PermissionsMap["roles/admin"] = Admin
}
