package panel

import (
	"context"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx/field"

	_ "embed"
)

//go:embed "templates/totp_enable.html"
var totpEnableHTML []byte
var totpEnable = custom.Parse[userFormContext]("totp_enable.html", totpEnableHTML, static.AssetsUser)

func (panel *UserPanel) routeTOTPEnable(ctx context.Context) http.Handler {
	tpl := totpEnable.Prepare(panel.Dependencies.Custom)

	return &httpx.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},

		RenderTemplate:        tpl.Template(),
		RenderTemplateContext: panel.UserFormContext2(tpl, component.MenuItem{Title: "Enable TOTP", Path: "/user/totp/enable/"}),

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password := values["password"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckPassword(r.Context(), []byte(password))
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				_, err := user.NewTOTP(r.Context())
				if err != nil {
					return struct{}{}, errTOTPSetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/totp/enroll", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_enroll.html"
var totpEnrollHTML []byte
var totpEnrollTemplate = custom.Parse[totpEnrollContext]("totp_enroll.html", totpEnrollHTML, static.AssetsUser)

type totpEnrollContext struct {
	userFormContext

	TOTPSecret string
	TOTPImage  template.URL
	TOTPURL    template.URL
}

func (panel *UserPanel) routeTOTPEnroll(ctx context.Context) http.Handler {
	tpl := totpEnrollTemplate.Prepare(panel.Dependencies.Custom, custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "User", Path: "/user/"},
			{Title: "Enable TOTP", Path: "/user/totp/enable/"},
		},
	})

	return &httpx.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Passcode"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && user.IsTOTPEnabled()
		},
		RenderForm: func(context httpx.FormContext, w http.ResponseWriter, r *http.Request) {
			// TODO: Do we want to reuse the same function here?

			user, err := panel.Dependencies.Auth.UserOf(r)

			ctx := totpEnrollContext{
				userFormContext: userFormContext{
					FormContext: context,
				},
			}

			if err == nil && user != nil {
				ctx.userFormContext.User = &user.User
				secret, err := user.TOTP()
				if err == nil {
					img, _ := auth.TOTPLink(secret, 500, 500)

					ctx.TOTPSecret = secret.Secret()
					ctx.TOTPImage = template.URL(img)
					ctx.TOTPURL = template.URL(secret.URL())
				}
			}
			tpl.Execute(w, r, ctx)
		},

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckPassword(r.Context(), []byte(password))
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				err := user.EnableTOTP(r.Context(), otp)
				if err != nil {
					return struct{}{}, errTOTPSetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}

//go:embed "templates/totp_disable.html"
var totpDisableHTML []byte
var totpDisableTemplate = custom.Parse[userFormContext]("totp_disable.html", totpDisableHTML, static.AssetsUser)

func (panel *UserPanel) routeTOTPDisable(ctx context.Context) http.Handler {
	tpl := totpDisableTemplate.Prepare(panel.Dependencies.Custom)

	return &httpx.Form[struct{}]{
		Fields: []field.Field{
			{Name: "password", Type: field.Password, Autocomplete: field.CurrentPassword, EmptyOnError: true, Label: "Current Password"},
			{Name: "otp", Type: field.Text, Autocomplete: field.OneTimeCode, EmptyOnError: true, Label: "Current Passcode"},
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		SkipForm: func(r *http.Request) (data struct{}, skip bool) {
			user, err := panel.Dependencies.Auth.UserOf(r)
			return struct{}{}, err == nil && user != nil && !user.IsTOTPEnabled()
		},
		RenderTemplate:        tpl.Template(),
		RenderTemplateContext: panel.UserFormContext2(tpl, component.MenuItem{Title: "Disable TOTP", Path: "/user/totp/disable/"}),

		Validate: func(r *http.Request, values map[string]string) (struct{}, error) {
			password, otp := values["password"], values["otp"]

			user, err := panel.Dependencies.Auth.UserOf(r)
			if err != nil {
				return struct{}{}, err
			}

			{
				err := user.CheckCredentials(r.Context(), []byte(password), otp)
				if err != nil {
					return struct{}{}, errCredentialsIncorrect
				}
			}
			{
				err := user.DisableTOTP(r.Context())
				if err != nil {
					return struct{}{}, errTOTPUnsetFailure
				}
			}

			return struct{}{}, nil
		},

		RenderSuccess: func(_ struct{}, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/user/", http.StatusSeeOther)
			return nil
		},
	}
}
