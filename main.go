// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"emqx-exporter/client"
	"emqx-exporter/config"
	"emqx-exporter/prober"

	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"gopkg.in/yaml.v3"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var sc = config.NewSafeConfig(prometheus.DefaultRegisterer)

func main() {
	os.Exit(run(kingpin.CommandLine, os.Args[1:], &http.Server{}))
}

func run(app *kingpin.Application, args []string, srv *http.Server) (exitCode int) {
	var (
		configFile             = app.Flag("config.file", "EMQX exporter configuration file.").Default(filepath.Join(filepath.Dir(os.Args[0]), "config.yaml")).String()
		maxProcs               = app.Flag("runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)").Envar("GOMAXPROCS").Default("4").Int()
		maxRequests            = app.Flag("web.max-requests", "Maximum number of parallel scrape requests. Use 0 to disable.").Default("40").Int()
		disableExporterMetrics = app.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Bool()
		toolkitFlags           = kingpinflag.AddFlags(app, ":8085")
	)
	app.Version(version.Print("emqx-exporter"))
	app.UsageWriter(os.Stdout)
	app.HelpFlag.Short('h')

	promlogConfig := &promlog.Config{}
	flag.AddFlags(app, promlogConfig)

	kingpin.MustParse(app.Parse(args))

	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting emqx-exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	// GOMAXPROCS returns the previous setting. If n < 1, it does not change the current setting.
	runtime.GOMAXPROCS(*maxProcs)
	level.Debug(logger).Log("msg", "Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	if err := sc.ReloadConfig(*configFile); err != nil {
		level.Error(logger).Log("msg", "Error loading config", "err", err)
		return 1
	}
	level.Info(logger).Log("msg", "Loaded config file")

	mux := http.NewServeMux()
	mux.Handle("/metrics", client.NewHandler(*disableExporterMetrics, *maxRequests, sc.C.Metrics, logger))

	mux.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		sc.Lock()
		probes := sc.C.Probes
		sc.Unlock()
		prober.Handler(w, r, probes, logger, nil)
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		sc.RLock()
		c, err := yaml.Marshal(sc.C)
		sc.RUnlock()
		if err != nil {
			level.Warn(logger).Log("msg", "Error marshalling configuration", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(c)
	})

	landingConfig := web.LandingConfig{
		Name:        "EMQX Exporter",
		Description: "EMQX Exporter",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{
				Address: "/metrics",
				Text:    "Metrics",
			},
			{
				Address: "/probe",
				Text:    "Probe",
			},
			{
				Address: "/config",
				Text:    "Config",
			},
		},
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		level.Error(logger).Log("err", err)
		return 1
	}
	mux.Handle("/", landingPage)

	srv.Handler = mux
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		return 1
	}

	return 0
}
