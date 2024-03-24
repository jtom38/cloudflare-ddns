package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidStatusCode  error = errors.New("did not get an acceptiable status code from the server")
	ErrFailedToDecodeBody error = errors.New("unable to decode the body")
	ErrFailedToDecodeJson error = errors.New("unexpected json format was returned")
	ErrWasNotJson         error = errors.New("response from server was not json")
	ErrDomainNotFound     error = errors.New("unable to find requested domain on cloudflare")
	ErrUnknownIpStack     error = errors.New("unknown ip stack")
)

type IpAddr struct {
	Ipv4 string
	Ipv6 string
}

func GetCurrentIpAddress(ipStack string) (IpAddr, error) {
	switch strings.ToLower(ipStack) {
	case "ipv4":
		ipv4, err := GetIpv4Addr()
		if err != nil {
			return IpAddr{}, err
		}
		return IpAddr{Ipv4: ipv4}, nil
	case "ipv6":
		ipv6, err := GetIpv4Addr()
		if err != nil {
			return IpAddr{}, err
		}
		return IpAddr{Ipv6: ipv6}, nil
	case "dual":
		ipv4, err := GetIpv4Addr()
		if err != nil {
			return IpAddr{}, fmt.Errorf("ipv4: %w", err)
		}
		ipv6, err := GetIpv6Addr()
		if err != nil {
			return IpAddr{}, fmt.Errorf("ipv6: %w", err)
		}
		return IpAddr{Ipv4: ipv4, Ipv6: ipv6}, nil
	default:
		return IpAddr{}, ErrUnknownIpStack
	}
}

func GetIpv4Addr() (string, error) {
	return GetAddr("https://v4.ident.me")
}

func GetIpv6Addr() (string, error) {
	return GetAddr("https://v6.ident.me")
}

func GetAddr(url string) (string, error) {
	resp, err := http.Get(url)
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
		cron.RunCloudflareCheck(cfg)
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
