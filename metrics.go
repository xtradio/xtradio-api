package main

import "github.com/prometheus/client_golang/prometheus"

var (
	songsPlayed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xtradio_songs_played",
			Help: "Number of songs played.",
		},
		[]string{"artist", "title", "show"},
	)
	tuneinSubmission = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xtradio_tunein_submission",
			Help: "Successfull TuneIn.com submissions.",
		},
		[]string{"artist", "title"},
	)
	dbConnectionFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "xtradio_db_connection_fail",
			Help: "Failed connections to the DB.",
		},
	)
	liqConnectionFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xtradio_liquidsoap_connection_fail",
			Help: "Failed number of connections to Liquidsoap",
		},
		[]string{"host", "port"},
	)
	// hdFailures = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "hd_errors_total",
	// 		Help: "Number of hard-disk errors.",
	// 	},
	// 	[]string{"device"},
	// )
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(songsPlayed)
	prometheus.MustRegister(tuneinSubmission)
	prometheus.MustRegister(dbConnectionFailure)
	prometheus.MustRegister(liqConnectionFailure)
	// prometheus.MustRegister(hdFailures)
}
