package panel

import (
	"context"
	"errors"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx/field"
)

//go:embed "templates/password.html"
var passwordHTML []byte
var passwordTemplate = templating.Parse[userFormContext](
	"password.html", passwordHTML, httpx.FormTemplate,

	templating.Title("Change Password"),
	templating.Assets(assets.AssetsUser),
)

var (
	errPasswordsNotIdentical = errors.New("passwords are not identical")
	errCredentialsIncorrect  = errors.New("credentials are not correct")
	errPasswordSetFailure    = errors.New("error saving new password")
	errTOTPSetFailure        = errors.New("unable to enable totp")
	errTOTPUnsetFailure      = errors.New("unable to disable totp")
	errPasswordSet           = errors.New("password was updated")
)

func (panel *UserPanel) routePassword(ctx context.Context) http.Handler {
	tpl := passwordTemplate.Prepare(panel.Dependencies.Templating)

	return &httpx.Form[struct{}]{
		Fields: []field.Field{
			{Name: "old", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Current Passcode (optional)"},
			{Name: "new", Type: field.Password, Autocomplete: field.NewPassword, EmptyOnError: true, Label: "New Password"},
			{Name: "new2", Type: field.Password, Autocomplete: field.NewPassword, EmptyOnError: true, Label: "New Password (again)"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderTemplate:        tpl.Template(),
		RenderTemplateContext: panel.UserFormContext(tpl, menuChangePassword),

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			old, passcode, new, new2 := values["old"], values["otp"], values["new"], values["new2"]

			if new != new2 {
				return struct{}{}, errPasswordsNotIdentical
			}

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckCredentials(r.Context(), []byte(old), passcode)
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}

			{
				err := user.CheckPasswordPolicy(new)
				if err != nil {
					return struct{}{}, err
				}
			}

			{
				err := user.SetPassword(r.Context(), []byte(new))
				if err != nil {
					return struct{}{}, errPasswordSetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			return errPasswordSet
		},
	}
}
