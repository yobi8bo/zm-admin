package v1

import (
	"zm-project/backend/internal/service"
	"zm-project/backend/internal/repository"
)

type Handler struct {
	Auth *AuthHandler
	User *UserHandler
	Role *RoleHandler
	Menu *MenuHandler
	Dept *DeptHandler
	Log  *LogHandler
	File *FileHandler
}

func NewHandler(
	authSvc *service.AuthService,
	userSvc *service.UserService,
	roleSvc *service.RoleService,
	menuSvc *service.MenuService,
	deptSvc *service.DeptService,
	logRepo *repository.LogRepo,
) *Handler {
	return &Handler{
		Auth: NewAuthHandler(authSvc),
		User: NewUserHandler(userSvc),
		Role: NewRoleHandler(roleSvc),
		Menu: NewMenuHandler(menuSvc),
		Dept: NewDeptHandler(deptSvc),
		Log:  NewLogHandler(logRepo),
		File: NewFileHandler(),
	}
}
