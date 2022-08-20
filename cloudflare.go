package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type CloudFlareClient struct {
	apiToken string
	email    string

	httpClient http.Client
}

func NewCloudFlareClient(ApiToken string, UserEmail string) *CloudFlareClient {
	c := CloudFlareClient{
		apiToken: ApiToken,
		email:    UserEmail,
	}
	c.httpClient = http.Client{}

	return &c
}

type listDomainZones struct {
	Result []struct {
		ID                  string      `json:"id,omitempty"`
		Name                string      `json:"name,omitempty"`
		Status              string      `json:"status,omitempty"`
		Paused              bool        `json:"paused,omitempty"`
		Type                string      `json:"type,omitempty"`
		DevelopmentMode     int         `json:"development_mode,omitempty"`
		NameServers         []string    `json:"name_servers,omitempty"`
		OriginalNameServers []string    `json:"original_name_servers,omitempty"`
		OriginalRegistrar   string      `json:"original_registrar,omitempty"`
		OriginalDnshost     interface{} `json:"original_dnshost,omitempty"`
		ModifiedOn          time.Time   `json:"modified_on,omitempty"`
		CreatedOn           time.Time   `json:"created_on,omitempty"`
		ActivatedOn         time.Time   `json:"activated_on,omitempty"`
		Meta                struct {
			Step                    int  `json:"step,omitempty"`
			CustomCertificateQuota  int  `json:"custom_certificate_quota,omitempty"`
			PageRuleQuota           int  `json:"page_rule_quota,omitempty"`
			PhishingDetected        bool `json:"phishing_detected,omitempty"`
			MultipleRailgunsAllowed bool `json:"multiple_railguns_allowed,omitempty"`
		} `json:"meta,omitempty"`
		Owner struct {
			ID    string `json:"id,omitempty"`
			Type  string `json:"type,omitempty"`
			Email string `json:"email,omitempty"`
		} `json:"owner,omitempty"`
		Account struct {
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"account,omitempty"`
		Tenant struct {
			ID   interface{} `json:"id,omitempty"`
			Name interface{} `json:"name,omitempty"`
		} `json:"tenant,omitempty"`
		TenantUnit struct {
			ID interface{} `json:"id,omitempty"`
		} `json:"tenant_unit,omitempty"`
		Permissions []string `json:"permissions,omitempty"`
		Plan        struct {
			ID                string `json:"id,omitempty"`
			Name              string `json:"name,omitempty"`
			Price             int    `json:"price,omitempty"`
			Currency          string `json:"currency,omitempty"`
			Frequency         string `json:"frequency,omitempty"`
			IsSubscribed      bool   `json:"is_subscribed,omitempty"`
			CanSubscribe      bool   `json:"can_subscribe,omitempty"`
			LegacyID          string `json:"legacy_id,omitempty"`
			LegacyDiscount    bool   `json:"legacy_discount,omitempty"`
			ExternallyManaged bool   `json:"externally_managed,omitempty"`
		} `json:"plan,omitempty"`
	} `json:"result,omitempty"`
	ResultInfo struct {
		Page       int `json:"page,omitempty"`
		PerPage    int `json:"per_page,omitempty"`
		TotalPages int `json:"total_pages,omitempty"`
		Count      int `json:"count,omitempty"`
		TotalCount int `json:"total_count,omitempty"`
	} `json:"result_info,omitempty"`
	Success  bool          `json:"success,omitempty"`
	Errors   []interface{} `json:"errors,omitempty"`
	Messages []interface{} `json:"messages,omitempty"`
}

// Lists out all the zones bound to an ac
//
// https://api.cloudflare.com/#zone-list-zones
func (c *CloudFlareClient) GetDomainByName(domain string) (*listDomainZones, error) {
	var items listDomainZones
	uri := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%v", domain)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return &items, err
	}
	req.Header.Set("X-Auth-Email", c.email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &items, err
	}

	if resp.StatusCode != 200 {
		return &items, ErrInvalidStatusCode
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &items, ErrWasNotJson
	}
	//log.Print(string(body))
	err = json.Unmarshal(body, &items)
	if err != nil {
		return &items, ErrFailedToDecodeJson
	}

	if !items.Success {
		log.Println("Failed to find the requested domain on Cloudflare.")
		return &items, ErrDomainNotFound
	}

	return &items, nil
}

type DnsDetails struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		ID         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxiable  bool      `json:"proxiable"`
		Proxied    bool      `json:"proxied"`
		TTL        int       `json:"ttl"`
		Locked     bool      `json:"locked"`
		ZoneID     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
		Data       struct {
		} `json:"data"`
		Meta struct {
			AutoAdded bool   `json:"auto_added"`
			Source    string `json:"source"`
		} `json:"meta"`
	} `json:"result"`
}

func (c *CloudFlareClient) GetDnsEntriesByDomain(DomainId string, Host string, Domain string) (*DnsDetails, error) {
	var items DnsDetails
	name := fmt.Sprintf("%v.%v", Host, Domain)
	uri := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records?name=%v", DomainId, name)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return &items, err
	}
	req.Header.Set("X-Auth-Email", c.email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &items, err
	}

	if resp.StatusCode != 200 {
		return &items, ErrInvalidStatusCode
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &items, ErrWasNotJson
	}
	err = json.Unmarshal(body, &items)
	if err != nil {
		return &items, ErrFailedToDecodeJson
	}

	if !items.Success {
		log.Println("Failed to find the requested domain on Cloudflare.")
		return &items, ErrDomainNotFound
	}

	return &items, nil
}

type dnsUpdate struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func (c *CloudFlareClient) UpdateDnsEntry(DomainId string, DnsDetails *DnsDetails, NewIpAddress string) error {
	param := dnsUpdate{
		Type: DnsDetails.Result[0].Type,
		Name: DnsDetails.Result[0].Name,
		Content: NewIpAddress,
		Ttl: DnsDetails.Result[0].TTL,
		Proxied: DnsDetails.Result[0].Proxied,
	}

	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", DomainId, DnsDetails.Result[0].ID)

	body, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Email", c.email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return errors.New("failed to update the IP address")
	}

	log.Println("IP Address request was sent and no errors reported.")
	return nil
}
