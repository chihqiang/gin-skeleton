package service

type Service struct {
	Auth AuthService
	Menu MenuService
	Role RoleService
	User UserService
}
