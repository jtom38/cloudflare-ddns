package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

type cronClient struct {
	scheduler *cron.Cron
}

func NewCron() cronClient {
	c := cronClient{
		scheduler: cron.New(),
	}

	return c
}

func (c cronClient) RunCloudflareCheck(cfg ConfigModel) {
	log.Println("Starting check...")
	log.Println("Checking the current IP Address")
	currentIp, err := GetCurrentIpAddress(cfg.IpStack)
	if err != nil {
		log.Println(err)
		return
	}

	cf := NewCloudFlareClient(cfg.Token, cfg.Email)
	log.Println("Checking domain information on CloudFlare")
	domainDetails, err := cf.GetDomainByName(cfg.Domain)
	if err != nil {
		log.Println("Unable to get information from CloudFlare.")
		log.Println("Double check the API Token to make sure it's valid.")
		log.Println(err)
		return
	}

	for _, host := range cfg.Hosts {
		hostname := fmt.Sprintf("%v.%v", host, cfg.Domain)
		log.Printf("Reviewing '%v'", hostname)
		dns, err := cf.GetDnsEntriesByDomain(domainDetails.Result[0].ID, host, cfg.Domain)
		if err != nil {
			log.Println("failed to collect dns entry")
			return
		}
		if currentIp.Ipv4 != "" {
			update(cf, currentIp.Ipv4, "A", dns)
		}
		if currentIp.Ipv6 != "" {
			update(cf, currentIp.Ipv6, "AAAA", dns)
		}

	}
	log.Println("Done!")
}

func update(cf *CloudFlareClient, ip, t string, dns *DnsDetails) {
	for _, item := range dns.Result {
		if item.Type == t && item.Content != ip {
			log.Printf("IP Address no longer matches, sending an update, from %s to %s\n", item.Content, ip)
			err := cf.UpdateDnsEntry(item, ip)
			if err != nil {
				log.Printf("Failed to update the DNS record to %s!\n", ip)
			}
		}
	}
}

func (c cronClient) HelloWorldJob() {
	log.Print("Hello World")
}
