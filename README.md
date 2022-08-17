# cloudflare-ddns

Golang based tool to maintain Cloudflare Dynamic DNS for your sites.  Run this locally or as a Docker image.  

When the application starts, it will run every 15 minutes and check your defined hosts and compare the IP Address listed.  If they dont match, it will be updated to refect the public IP address of the server its running from.

This does not have any UI elements so you need to check the logs to se how its going.


## Usage

```yaml
# docker-compose.yaml
version: "3"
    
services:
  app:
    image: ghcr.io/jtom38/cloudflare-ddns:master
    container_name: cfddns
    environment:
      EMAIL: "yourcloudflareemailaddress"
      API_TOKEN: "cloudflare-api-key"
      DOMAIN: "exampledomain.com"
      HOSTS: "example1,example2,www"

```

## Credit

The original post that gave me this idea can be found [here](https://adamtheautomator.com/cloudflare-dynamic-dns/).  This was written in PowerShell and thought I could improve on it with Go and Docker.