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

func (c cronClient) RunCloudflareCheck(ApiToken string, Email string, Domain string, Hosts []string) {
	log.Println("Starting check...")
	log.Println("Checking the current IP Address")
	currentIp, err := GetCurrentIpAddress()
	if err != nil {
		log.Println(err)
		return
	}

	cf := NewCloudFlareClient(ApiToken, Email)
	log.Println("Checking domain information on CloudFlare")
	domainDetails, err := cf.GetDomainByName(Domain)
	if err != nil {
		log.Println("Unable to get information from CloudFlare.")
		log.Println("Double check the API Token to make sure it's valid.")
		log.Println(err)
		return
	}

	for _, host := range Hosts {
		hostname := fmt.Sprintf("%v.%v", host, Domain)
		log.Printf("Reviewing '%v'", hostname)
		dns, err := cf.GetDnsEntriesByDomain(domainDetails.Result[0].ID, host, Domain)
		if err != nil {
			log.Println("failed to collect dns entry")
			return
		}

		var result = dns.Result[0]
		if result.Content != currentIp {
			log.Println("IP Address no longer matches, sending an update")
			err = cf.UpdateDnsEntry(domainDetails.Result[0].ID, dns, currentIp)
			if err != nil {
				log.Println("Failed to update the DNS record!")
			}
		}
	}
	log.Println("Done!")
}

func (c cronClient) HelloWorldJob() {
	log.Print("Hello World")
}
