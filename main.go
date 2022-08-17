package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidStatusCode error = errors.New("did not get an acceptiable status code from the server")
	ErrFailedToDecodeBody error = errors.New("unable to decode the body")
	ErrFailedToDecodeJson error = errors.New("unexpected json format was returned")
	ErrWasNotJson error = errors.New("response from server was not json")
	ErrDomainNotFound error = errors.New("unable to find requested domain on cloudflare")
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
	email := config.GetConfig(ConfigEmail)
	if email == "" {
		log.Println("Unable to find 'EMAIL' env value.")
		return
	}

	token := config.GetConfig(ConfigToken)
	if token == "" {
		log.Println("Unable to find 'API_TOKEN' env value.")
	}

	domain := config.GetConfig(ConfigDomain)
	if token == "" {
		log.Println("Unable to find 'DOMAIN' env value.")
	}

	hosts := config.GetConfig(ConfigHosts)
	if token == "" {
		log.Println("Unable to find 'HOSTS' env value.")
	}
	hostsArray := strings.Split(hosts, ",")
	log.Println("Env Check: OK")

	cron := NewCron()
	log.Println("Cloudflare Check will run every 15 minutes.")
	cron.scheduler.AddFunc("0,15,30,45 * * * *", func() { 
		cron.RunCloudflareCheck(token, email, domain, hostsArray)
	})
	cron.scheduler.Start()

	log.Println("Application has started!")
	for {
		time.Sleep(30 * time.Second)
	}

}

