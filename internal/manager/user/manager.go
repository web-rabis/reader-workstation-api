package user

import (
	"context"

	ssoClient "github.com/web-rabis/sso-client/client"
	ssoModel "github.com/web-rabis/sso-client/model"
)

type Manager struct {
	userCl ssoClient.UserService
}

func NewUserManager(userCl ssoClient.UserService) *Manager {
	return &Manager{
		userCl: userCl,
	}
}

func (m *Manager) ById(ctx context.Context, id int) (*ssoModel.User, error) {
	return m.userCl.ById(ctx, int64(id))
}
