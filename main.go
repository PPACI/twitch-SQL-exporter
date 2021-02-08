package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"time"
	"twitch-SQL-exporter/pkg/helix"
	"twitch-SQL-exporter/pkg/sql"
)

type appConfig struct {
	clientID     string
	clientSecret string
	debug        bool

	dbType string
	dsn    string
}

func getAppConfig() *appConfig {
	config := appConfig{}

	flag.StringVar(&config.clientID, "clientID", "", "Twitch client ID.")
	flag.StringVar(&config.clientSecret, "clientSecret", "", "Twitch client Secret.")
	flag.StringVar(&config.dbType, "dbType", "sqlite", "Type of DB. Support sqlite,postgres.")
	flag.StringVar(&config.dsn, "dataSourceName", "", "DB Data Source Name.")
	flag.BoolVar(&config.debug, "debug", false, "Enable debug mode.")
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

	log.Infoln("Connect to DB")
	db, err := sql.NewDb(&sql.DbOpts{
		Type: config.dbType,
		Dsn:  config.dsn,
		Opts: &gorm.Config{},
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Update DB schema")
	err = db.AutoMigrate(&sql.StreamRecord{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Initiate twitch client")
	clientOpts := helix.ClientOpts{
		ClientID:     config.clientID,
		ClientSecret: config.clientSecret,
	}

	twitchClient := helix.NewClient(&clientOpts, context.Background())

	log.Infoln("Start polling loop")
	pollData(twitchClient, db)
	for range time.Tick(time.Minute) {
		pollData(twitchClient, db)
	}

}

func pollData(twitchClient *helix.Client, db *gorm.DB) {
	streams, err := twitchClient.GetStreams(&helix.GetStreamsOpts{Language: "fr", First: 100})
	if err != nil {
		log.Fatalln(err)
	}
	var streamRecords []sql.StreamRecord
	for _, stream := range streams.Data {
		log.Debugf("%+v\n", stream)
		streamRecord := sql.StreamRecord{
			StreamId:    stream.Id,
			GameId:      stream.GameId,
			GameName:    stream.GameName,
			Language:    stream.Language,
			StartedAt:   stream.StartedAt,
			Title:       stream.Title,
			UserName:    stream.UserName,
			UserLogin:   stream.UserLogin,
			UserId:      stream.UserId,
			ViewerCount: stream.ViewerCount,
		}
		streamRecords = append(streamRecords, streamRecord)
	}
	log.Infoln("Fetched stream records")
	db.Create(streamRecords)
}
