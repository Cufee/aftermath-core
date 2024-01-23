package permissions

const (
	User = UseCommands | SelectPersonalBackgroundPreset | CreatePersonalConnection | UpdatePersonalConnection | RemovePersonalConnection | RetrievePersonalSubscriptions | CreatePersonalSubscription | ExtendPersonalSubscription

	ContentModerator = User | UpdateUserBackground | RemoveUserBackground | CreateUserRestriction
	GlobalModerator  = ContentModerator | RetrieveUserSubscriptions | CreateUserSubscription | ExtendUserSubscription | TerminateUserSubscription | UploadBackgroundPreset | RemoveBackgroundPreset | RetrieveUserConnections | ManageUserConnectionVerification | RemoveUserConnection | RetrieveUserRestrictions | RemoveUserRestriction

	Admin = GlobalModerator | ManageUserRoles
)
