package admin

import (
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance.html"
var instanceHTML []byte
var instanceTemplate = templating.Parse[instanceContext](
	"instance.html", instanceHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceContext struct {
	templating.RuntimeFlags

	Instance models.Instance
	Info     status.WissKI
}

func (admin *Admin) instance(ctx context.Context) http.Handler {
	tpl := instanceTemplate.Prepare(
		admin.Dependencies.Templating,
		templating.Crumbs(
			component.MenuItem{Title: "Admin", Path: "/admin/"},
			component.DummyMenuItem,
		),
		templating.Actions(
			component.DummyMenuItem,
			component.DummyMenuItem,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (ic instanceContext, funcs []templating.FlagFunc, err error) {
		slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

		// find the instance itself!
		instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
		if err == instances.ErrWissKINotFound {
			return ic, nil, httpx.ErrNotFound
		}
		if err != nil {
			return ic, nil, err
		}
		ic.Instance = instance.Instance

		// get some more info about the wisski
		ic.Info, err = instance.Info().Information(r.Context(), false)
		if err != nil {
			return ic, nil, err
		}

		funcs = []templating.FlagFunc{
			templating.ReplaceCrumb(1, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}),
			templating.ReplaceAction(0, component.MenuItem{Title: "Grants", Path: template.URL("/admin/grants/" + slug)}),
			templating.ReplaceAction(1, component.MenuItem{Title: "Ingredients", Path: template.URL("/admin/ingredients/" + slug), Priority: component.SmallButton}),

			templating.Title(instance.Slug),
		}

		return
	})
}
