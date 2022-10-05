package control

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	component.ComponentBase

	Instances *instances.Instances
}

func (Info) Name() string { return "control-info" }

func (*Info) Routes() []string { return []string{"/dis/"} }

func (info *Info) Handler(route string, io stream.IOStream) (http.Handler, error) {
	mux := http.NewServeMux()

	// handle everything under /dis/!
	mux.HandleFunc("/dis/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dis/" {
			http.Redirect(w, r, "/dis/index", http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	})

	// static stuff
	static, err := info.disStatic()
	if err != nil {
		return nil, err
	}
	mux.Handle("/dis/static/", static)

	// render everything
	mux.Handle("/dis/index", httpx.HTMLHandler[disIndex]{
		Handler:  info.disIndex,
		Template: indexTemplate,
	})

	mux.Handle("/dis/instance/", httpx.HTMLHandler[disInstance]{
		Handler:  info.disInstance,
		Template: instanceTemplate,
	})

	// api -- for future usage
	mux.Handle("/dis/api/v1/instance/get/", httpx.JSON(info.getinstance))
	mux.Handle("/dis/api/v1/instance/all", httpx.JSON(info.allinstances))

	// ensure that everyone is logged in!
	return httpx.BasicAuth(mux, "WissKI Distillery Admin", func(user, pass string) bool {
		return user == info.Config.DisAdminUser && pass == info.Config.DisAdminPassword
	}), nil
}

// disIndex is the context of the "/dis/index" page
type disIndex struct {
	Time time.Time

	Config *config.Config

	Instances []instances.WissKIInfo

	TotalCount   int
	RunningCount int
	StoppedCount int

	Backups []models.Snapshot
}

func (info *Info) disIndex(r *http.Request) (idx disIndex, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// load instances
		idx.Instances, err = info.allinstances(r)
		if err != nil {
			return err
		}

		// count how many are running and how many are stopped
		for _, i := range idx.Instances {
			if i.Running {
				idx.RunningCount++
			} else {
				idx.StoppedCount++
			}
		}
		idx.TotalCount = len(idx.Instances)

		return nil
	})

	// get the log entries
	group.Go(func() (err error) {
		idx.Backups, err = info.Instances.SnapshotLogFor("")
		return
	})

	// get the static properties
	idx.Config = info.Config

	// current time
	idx.Time = time.Now().UTC()

	// wait for everything!
	group.Wait()

	return
}

// disInstance is the context of the "/dis/instance/*" page
type disInstance struct {
	Time time.Time

	Instance models.Instance
	Info     instances.WissKIInfo
}

func (info *Info) disInstance(r *http.Request) (is disInstance, err error) {
	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// find the instance itself!
	instance, err := info.Instances.WissKI(slug)
	if err == instances.ErrWissKINotFound {
		return is, httpx.ErrNotFound
	}
	if err != nil {
		return is, err
	}
	is.Instance = instance.Instance

	// get some more info about the wisski
	is.Info, err = instance.Info(false)
	if err != nil {
		return is, err
	}

	// current time
	is.Time = time.Now().UTC()

	return
}

//go:embed html/static
var htmlStaticFS embed.FS

func (*Info) disStatic() (http.Handler, error) {
	fs, err := fs.Sub(htmlStaticFS, "html/static")
	if err != nil {
		return nil, err
	}

	return http.StripPrefix("/dis/static/", http.FileServer(http.FS(fs))), nil
}

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = template.Must(template.New("index.html").Parse(indexTemplateStr))

//go:embed "html/instance.html"
var instanceTemplateString string
var instanceTemplate = template.Must(template.New("instance.html").Parse(instanceTemplateString))

func (info *Info) getinstance(r *http.Request) (iinfo instances.WissKIInfo, err error) {
	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// load the wisski instance!
	wisski, err := info.Instances.WissKI(strings.TrimSuffix(slug, "/"))
	if err == instances.ErrWissKINotFound {
		return iinfo, httpx.ErrNotFound
	}
	if err != nil {
		return iinfo, err
	}

	// get info about it!
	return wisski.Info(false)
}

func (info *Info) allinstances(*http.Request) (infos []instances.WissKIInfo, err error) {
	var errgroup errgroup.Group

	// list all the instances
	all, err := info.Instances.All()
	if err != nil {
		return nil, err
	}

	// get all of their info!
	infos = make([]instances.WissKIInfo, len(all))
	for i, instance := range all {
		{
			i := i
			instance := instance

			errgroup.Go(func() (err error) {
				infos[i], err = instance.Info(true)
				return err
			})
		}
	}

	// wait for the results, and return
	err = errgroup.Wait()
	return
}