package scopes

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type UserLoggedIn struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*UserLoggedIn)(nil)
)

const (
	ScopeUserValid Scope = "user.valid"
)

func (*UserLoggedIn) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:       ScopeUserValid,
		Description: "session must have a valid user",
		TakesParam:  false,
	}
}

func (iu *UserLoggedIn) HasScope(param string, r *http.Request) (bool, error) {
	_, user, err := iu.dependencies.Auth.SessionOf(r)
	return user != nil, err
}
