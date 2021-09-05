package main

import (
	"context"
	"fmt"
	"os"

	"weathertimemachineinfra/create"

	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

func main() {
	apiKey := os.Getenv("VULTR_APIKEY")

	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: apiKey})
	vultrClient := govultr.NewClient(oauth2.NewClient(ctx, ts))

	// Optional changes
	_ = vultrClient.SetBaseURL("https://api.vultr.com")
	vultrClient.SetUserAgent("mycool-app")
	vultrClient.SetRateLimit(500)

	sshkeyID := create.CreateSSH(vultrClient)
	networkID := create.CreatePrivateNetwork(vultrClient)
	instance := create.CreateInstance(vultrClient, networkID, sshkeyID)

	fmt.Println(instance)

}
