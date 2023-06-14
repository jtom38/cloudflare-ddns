package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrInvalidStatusCode  error = errors.New("did not get an acceptiable status code from the server")
	ErrFailedToDecodeBody error = errors.New("unable to decode the body")
	ErrFailedToDecodeJson error = errors.New("unexpected json format was returned")
	ErrWasNotJson         error = errors.New("response from server was not json")
	ErrDomainNotFound     error = errors.New("unable to find requested domain on cloudflare")
)

func GetCurrentIpAddress() (string, error) {
	resp, err := http.Get("https://v4.ident.me")
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", ErrInvalidStatusCode
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrFailedToDecodeBody
	}

	return string(body), nil
}

func main() {
	config := NewConfigClient()
	cfg := config.LoadConfig()

	if cfg.Email == "" {
		log.Println("Unable to find 'EMAIL' env value.")
		return
	}

	if cfg.Token == "" {
		log.Println("Unable to find 'API_TOKEN' env value.")
		return
	}

	if cfg.Domain == "" {
		log.Println("Unable to find 'DOMAIN' env value.")
		return
	}

	if len(cfg.Hosts) == 0 {
		log.Println("Unable to find 'HOSTS' env value.")
	}

	log.Println("Config Check: OK")

	cron := NewCron()
	log.Println("Cloudflare Check will run every 15 minutes.")
	cron.scheduler.AddFunc("0/5 * * * *", func() {
		cron.RunCloudflareCheck(cfg.Token, cfg.Email, cfg.Domain, cfg.Hosts)
	})
	cron.scheduler.AddFunc("0/1 * * * *", func() {
		cron.HelloWorldJob()
	})
	cron.scheduler.Start()

	log.Println("Application has started!")
	for {
		time.Sleep(30 * time.Second)
	}

}
