package search

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

var client *elastic.Client

func GetElasticClient() *elastic.Client {
	if client == nil {
		var err error
		conf, err := config.Parse(os.Getenv("ELASTICURL"))
		if err != nil {
			panic("invalid elastic search connection string")
		}
		options, err := configToOptions(conf)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err.Error())
		}
		options = append(options, elastic.SetHealthcheckTimeoutStartup(60*time.Second))
		client, err = elastic.DialContext(context.Background(), options...)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err.Error())
		}
	}
	return client
}


func configToOptions(cfg *config.Config) ([]elastic.ClientOptionFunc, error) {
	var options []elastic.ClientOptionFunc
	if cfg != nil {
		if cfg.URL != "" {
			options = append(options, elastic.SetURL(cfg.URL))
		}
		if cfg.Errorlog != "" {
			f, err := os.OpenFile(cfg.Errorlog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize error log")
			}
			l := log.New(f, "", 0)
			options = append(options, elastic.SetErrorLog(l))
		}
		if cfg.Tracelog != "" {
			f, err := os.OpenFile(cfg.Tracelog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize trace log")
			}
			l := log.New(f, "", 0)
			options = append(options, elastic.SetTraceLog(l))
		}
		if cfg.Infolog != "" {
			f, err := os.OpenFile(cfg.Infolog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, errors.Wrap(err, "unable to initialize info log")
			}
			l := log.New(f, "", 0)
			options = append(options, elastic.SetInfoLog(l))
		}
		if cfg.Username != "" || cfg.Password != "" {
			options = append(options, elastic.SetBasicAuth(cfg.Username, cfg.Password))
		}
		if cfg.Sniff != nil {
			options = append(options, elastic.SetSniff(*cfg.Sniff))
		}
		/*
			if cfg.Healthcheck != nil {
				options = append(options, SetHealthcheck(*cfg.Healthcheck))
			}
		*/
	}
	return options, nil
}
