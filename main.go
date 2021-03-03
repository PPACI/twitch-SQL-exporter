package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"sync"
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

	minViewers int
}

func getAppConfig() *appConfig {
	config := appConfig{}

	flag.StringVar(&config.clientID, "clientID", "", "Twitch client ID.")
	flag.StringVar(&config.clientSecret, "clientSecret", "", "Twitch client Secret.")
	flag.StringVar(&config.dbType, "dbType", "sqlite", "Type of DB. Support sqlite,postgres.")
	flag.StringVar(&config.dsn, "dataSourceName", "", "DB Data Source Name.")
	flag.IntVar(&config.minViewers, "minViewer", 500, "Minimum viewer for record")
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
	pollData(twitchClient, db, config.minViewers)
	for range time.Tick(time.Minute) {
		pollData(twitchClient, db, config.minViewers)
	}

}

func pollData(twitchClient *helix.Client, db *gorm.DB, minViewer int) {
	streams, err := twitchClient.GetStreams(&helix.GetStreamsOpts{Language: "fr", First: 100})
	if err != nil {
		log.Fatalln(err)
	}
	var streamRecords []sql.StreamRecord
	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	for _, stream := range streams.Data {
		stream := stream
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Debugf("%+v\n", stream)
			if stream.ViewerCount < minViewer {
				log.WithField("name", stream.UserName).
					WithField("viewer", stream.ViewerCount).
					Debugln("Skipped due to not enough viewer")
				return
			}
			follow, err := twitchClient.GetFollows(&helix.GetFollowsOpts{UserId: stream.UserId})
			if err != nil {
				log.WithField("name", stream.UserName).Errorln(err)
			}
			streamRecord := sql.StreamRecord{
				StreamId:      stream.Id,
				GameId:        stream.GameId,
				GameName:      stream.GameName,
				Language:      stream.Language,
				StartedAt:     stream.StartedAt,
				Title:         stream.Title,
				UserName:      stream.UserName,
				UserLogin:     stream.UserLogin,
				UserId:        stream.UserId,
				ViewerCount:   stream.ViewerCount,
				FollowerCount: follow.Total,
			}
			m.Lock()
			streamRecords = append(streamRecords, streamRecord)
			m.Unlock()
		}()
	}
	wg.Wait()
	log.Infof("Fetched stream records")
	results := db.Create(streamRecords)
	log.Debugln("affected", results.RowsAffected, "rows")
}
