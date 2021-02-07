package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"twitch-SQL-exporter/pkg/helix"
)

type appConfig struct {
	clientID     string
	clientSecret string
	debug        bool
}

func getAppConfig() *appConfig {
	config := appConfig{}
	flag.StringVar(&config.clientID, "clientID", "", "Twitch client ID")
	flag.StringVar(&config.clientSecret, "clientSecret", "", "Twitch client Secret")
	flag.BoolVar(&config.debug, "debug", false, "Enable debug mode")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.DefValue == "" && f.Value.String() == "" {
			flag.Usage()
			os.Exit(1)
		}
	})
	return &config
}

func main() {
	config := getAppConfig()
	if config.debug {
		log.SetLevel(log.DebugLevel)
	}

	clientOpts := helix.ClientOpts{
		ClientID:     config.clientID,
		ClientSecret: config.clientSecret,
	}

	twitchClient := helix.NewClient(&clientOpts, context.Background())

	streams, err := twitchClient.GetStreams(&helix.GetStreamsOpts{First: 100})
	if err != nil {
		log.Fatalln(err)
	}

	for _, stream := range streams.Data {
		log.Debugf("%+v\n", stream)
	}
	log.Debugln(len(streams.Data))

}
