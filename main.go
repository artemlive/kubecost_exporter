package main

import (
	"fmt"
	"github.com/artemlive/kubecost_exporter/collector"
	"github.com/artemlive/kubecost_exporter/version"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"strings"
)

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default(":9150").String()
	metricPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	tlsInsecureSkipVerify = kingpin.Flag(
		"tls.insecure-skip-verify",
		"Ignore certificate and server verification when using a tls connection.",
	).Bool()
	kubecostUrl = kingpin.Flag("kubecost.baseUrl", "KubeCost base URL with schema: https://kubecost.example.com").Required().Envar("KUBECOST_URL").URL()
)

// scrapers lists all possible collection methods and if they should be enabled by default.
// Reserved for future use cases, if there will be other endpoints
var scrapers = map[collector.Scraper]bool{
	collector.ScrapeAssets{}: true,
	collector.ScrapeAllocation{}: true,
}

func newHandler(metrics collector.Metrics, scrapers []collector.Scraper, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filteredScrapers := scrapers
		scrapersFilterQuery := r.URL.Query()["collect[]"]
		scrapersParams := make(map[string][]string)
		// Use request context for cancellation when connection gets closed.
		ctx := r.Context()
		level.Debug(logger).Log("msg", "collect[] scrapersFilterQuery", "scrapersFilterQuery", strings.Join(scrapersFilterQuery, ","))

		// Check if we have some "collect[]" query parameters.
		if len(scrapersFilterQuery) > 0 {
			filters := make(map[string]bool)
			for _, param := range scrapersFilterQuery{
				filters[param] = true
			}

			filteredScrapers = nil
			for _, scraper := range scrapers {
				if filters[scraper.Name()] {
					filteredScrapers = append(filteredScrapers, scraper)
					// You can define additional HTTP query parameters for each scraper
					// scraper_name[] = window=1d,api_key=SECRET_KEY
					scrapersParamsKey := fmt.Sprintf("%s[]", scraper.Name())
					scrapersParamsQuery := r.URL.Query()[scrapersParamsKey]
					if len(scrapersParamsQuery) > 0 {
						scrapersParams[scraper.Name()] = strings.Split(scrapersParamsQuery[0], ",")
						level.Debug(logger).Log("msg", fmt.Sprintf("%s params query", scrapersParamsKey), "scrapersParamsQuery", strings.Join(scrapersParamsQuery, ","))
					}
				}
			}

		}

		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.New(ctx, kubecostUrl, metrics, filteredScrapers, scrapersParams, logger, *tlsInsecureSkipVerify))

		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func main() {
	// Generate ON/OFF flags for all scrapers.
	scraperFlags := map[collector.Scraper]*bool{}
	for scraper, enabledByDefault := range scrapers {
		defaultOn := "false"
		if enabledByDefault {
			defaultOn = "true"
		}

		f := kingpin.Flag(
			"collect."+scraper.Name(),
			scraper.Help(),
		).Default(defaultOn).Bool()

		scraperFlags[scraper] = f
	}

	// Parse flags.
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("kubecost_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	// landingPage contains the HTML served at '/'.
	// TODO: Make this nicer and more informative.
	var landingPage = []byte(`<html>
<head><title>KubeCost exporter</title></head>
<body>
<h1>KubeCost Assets API exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)
	level.Info(logger).Log("msg", "Starting kubecost_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	// Register only scrapers enabled by flag.
	// As for now we have only one scraper that gets the info about assets
	var enabledScrapers []collector.Scraper
	for scraper, enabled := range scraperFlags {
		if *enabled {
			level.Info(logger).Log("msg", "Scraper enabled", "scraper", scraper.Name())
			enabledScrapers = append(enabledScrapers, scraper)
		}
	}
	handlerFunc := newHandler(collector.NewMetrics(), enabledScrapers, logger)
	http.Handle(*metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
