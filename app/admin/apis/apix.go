package apis

import "context"

type Apis struct {
	Auth *AuthApis
	Menu *MenuApis
	Role *RoleApis
	User *UserApis
}

func NewApis(ctx context.Context) *Apis {
	return &Apis{
		Auth: NewAuth(ctx),
		Menu: NewMenu(ctx),
		Role: NewRole(ctx),
		User: NewUser(ctx),
	}
}
