package info

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/lifetime"
	"github.com/tkw1536/pkglib/sema"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	ingredient.Base
	Dependencies struct {
		PHP      *php.PHP
		Fetchers []ingredient.WissKIFetcher
	}

	Analytics *lifetime.Analytics
}

var (
	_ ingredient.WissKIFetcher = (*Info)(nil)
)

// Information fetches information about this WissKI.
// TODO: Rework this to be able to determine what kind of information is available.
func (wisski *Info) Information(ctx context.Context, quick bool) (info status.WissKI, err error) {
	// setup flags
	flags := ingredient.FetcherFlags{
		Quick:   quick,
		Context: ctx,
	}

	var serversUsed uint64
	pool := sema.Pool[*phpx.Server]{
		// limit the number of processes running in this container
		// to avoid long overheads
		Limit: 5,
		New: func() *phpx.Server {
			atomic.AddUint64(&serversUsed, 1)
			return wisski.Dependencies.PHP.NewServer()
		},
		Discard: func(s *phpx.Server) {
			s.Close()
		},
	}
	defer pool.Close()

	// setup a dictionary to record data about how long each operation took.
	// we use a slice as opposed to a map to avoid having to mutex!
	fetcherTimes := make([]time.Duration, len(wisski.Dependencies.Fetchers))
	recordTime := func(i int) func() {
		start := time.Now()
		return func() {
			fetcherTimes[i] = time.Since(start)
		}
	}

	start := time.Now()
	{
		var group errgroup.Group
		for i, fetcher := range wisski.Dependencies.Fetchers {
			fetcher, flags, i := fetcher, flags, i
			group.Go(func() error {
				// quick: don't need to create servers
				if flags.Quick {
					defer recordTime(i)()
					return fetcher.Fetch(flags, &info)
				}

				// complete: need to use a server from the pool
				return pool.Use(func(s *phpx.Server) error {
					defer recordTime(i)()
					flags.Server = s
					return fetcher.Fetch(flags, &info)
				})
			})
		}

		// wait for all the results
		err = group.Wait()
	}
	took := time.Since(start)

	var tookSum time.Duration

	// get a map of how long each fetcher took
	times := zerolog.Dict()
	for i, fetcher := range wisski.Dependencies.Fetchers {
		tookSum += fetcherTimes[i]
		times = times.Dur(fetcher.Name(), fetcherTimes[i])
	}

	// compute the ratio taken
	tookRatio := float64(took) / float64(tookSum)

	// and send it to debugging output
	zerolog.Ctx(ctx).Debug().Uint64("servers", serversUsed).Dict("fetchers_took_ms", times).Dur("took_ms", took).Dur("took_sum_ms", tookSum).Float64("took_ratio", tookRatio).Bool("quick", quick).Msg("ran information fetchers")

	return
}

func (wisski *Info) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) error {
	info.Time = time.Now().UTC()
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()
	return nil
}
