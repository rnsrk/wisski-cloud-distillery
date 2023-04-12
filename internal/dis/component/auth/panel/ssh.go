package panel

import (
	"context"
	"errors"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/gliderlabs/ssh"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/field"

	gossh "golang.org/x/crypto/ssh"

	_ "embed"
)

//go:embed "templates/ssh.html"
var sshHTML []byte
var sshTemplate = templating.Parse[SSHTemplateContext](
	"ssh.html", sshHTML, nil,

	templating.Title("SSH Keys"),
	templating.Assets(assets.AssetsUser),
)

type SSHTemplateContext struct {
	templating.RuntimeFlags

	Keys []models.Keys

	Domain string // domain name of the distillery
	Port   uint16 // public port of the distillery ssh servers

	Slug     string // slug of the wisski
	Hostname string // hostname of an example wisski

	// Services are the special services reachable via ssh
	Services []ssh2.Intercept
}

func (panel *UserPanel) sshRoute(ctx context.Context) http.Handler {
	tpl := sshTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuSSH,
		),
		templating.Actions(
			menuSSHAdd,
		),
	)

	return tpl.HTMLHandler(func(r *http.Request) (sc SSHTemplateContext, err error) {
		user, err := panel.Dependencies.Auth.UserOf(r)
		if err != nil {
			return sc, err
		}

		sc.Domain = panel.Config.HTTP.PrimaryDomain
		sc.Port = panel.Config.Listen.SSHPort

		// pick the first domain that the user has access to as an example
		grants, err := panel.Dependencies.Policy.User(r.Context(), user.User.User)
		if err != nil && len(grants) > 0 {
			sc.Slug = grants[0].Slug
		} else {
			sc.Slug = "example"
		}
		sc.Hostname = panel.Config.HTTP.HostFromSlug(sc.Slug)

		sc.Keys, err = panel.Dependencies.Keys.Keys(r.Context(), user.User.User)
		if err != nil {
			return sc, err
		}

		sc.Services = panel.Dependencies.SSH2.Intercepts()

		return sc, nil
	})
}

var (
	errInvalidUser = errors.New("invalid user")
	errKeyParse    = errors.New("unable to parse ssh key")
	errAddKey      = errors.New("unable to add key")
)

func (panel *UserPanel) sshDeleteRoute(ctx context.Context) http.Handler {
	logger := zerolog.Ctx(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Err(err).Str("action", "delete ssh key").Msg("failed to parse form")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}
		user, err := panel.Dependencies.Auth.UserOf(r)
		if err != nil {
			logger.Err(err).Str("action", "delete ssh key").Msg("failed to get current user")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		key, _ := parseKey(r.PostFormValue("signature"))
		if key == nil {
			logger.Err(err).Str("action", "delete ssh key").Msg("failed to parse signature")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		if err := panel.Dependencies.Keys.Remove(r.Context(), user.User.User, key); err != nil {
			logger.Err(err).Str("action", "delete ssh key").Msg("failed to delete key")
			httpx.HTMLInterceptor.Fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/user/ssh/", http.StatusSeeOther)
	})
}

//go:embed "templates/ssh_add.html"
var sshAddHTML []byte
var sshAddTemplate = templating.ParseForm(
	"ssh_add.html", sshAddHTML, httpx.FormTemplate,
	templating.Title("Add SSH Key"),
	templating.Assets(assets.AssetsUser),
)

type addKeyResult struct {
	User    *auth.AuthUser
	Comment string
	Key     ssh.PublicKey
}

func (panel *UserPanel) sshAddRoute(ctx context.Context) http.Handler {
	tpl := sshAddTemplate.Prepare(
		panel.Dependencies.Templating,
		templating.Crumbs(
			menuUser,
			menuSSH,
			menuSSHAdd,
		),
	)

	return &httpx.Form[addKeyResult]{
		Fields: []field.Field{
			{Name: "comment", Type: field.Text, Label: "Comment"},
			{Name: "key", Type: field.Textarea, Label: "Key in authorized_keys format"}, // has hacked css!
		},
		FieldTemplate: field.PureCSSFieldTemplate,

		RenderTemplate:        tpl.Template(),
		RenderTemplateContext: templating.FormTemplateContext(tpl),

		Validate: func(r *http.Request, values map[string]string) (ak addKeyResult, err error) {
			ak.User, err = panel.Dependencies.Auth.UserOf(r)
			if err != nil || ak.User == nil {
				return ak, errInvalidUser
			}

			// parse key and comment
			var key, comment string
			ak.Comment, key = values["comment"], values["key"]
			ak.Key, comment = parseKey(key)
			if ak.Key == nil {
				return ak, errKeyParse
			}

			// set the comment if the user didn't provide one!
			if ak.Comment == "" && comment != "" {
				ak.Comment = comment
			}
			return ak, nil
		},

		RenderSuccess: func(ak addKeyResult, values map[string]string, w http.ResponseWriter, r *http.Request) error {
			// add the key to the user
			if err := panel.Dependencies.Keys.Add(r.Context(), ak.User.User.User, ak.Comment, ak.Key); err != nil {
				return errAddKey
			}
			// everything went fine, redirect the user back to the user page!
			http.Redirect(w, r, "/user/ssh/", http.StatusSeeOther)
			return nil
		},
	}
}

func parseKey(authorized_keys string) (out gossh.PublicKey, comment string) {
	var err error
	out, comment, _, _, err = gossh.ParseAuthorizedKey([]byte(authorized_keys))
	if err != nil || out == nil {
		return nil, ""
	}
	return out, comment
}
