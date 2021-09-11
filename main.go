package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"weathertimemachineinfra/methods"

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

	vars := &methods.Input{
		SshName:            "test-02",
		SshPubKey:          os.Getenv("WTM_PUB_KEY"),
		Region:             "ams",
		NetworkDescription: "Private network for WTM",
		NetworkSubnet:      os.Getenv("WTM_CIDR"),
		NetworkSubnetMask:  24,
		InstanceLabel:      "wtm-01",
		InstanceHostname:   "wtm-01",
		InstancePlan:       "sdfdsfds",
		InstanceOSId:       445,
	}

	//Method include:
	//1 = create
	//2 = list
	//3 = delete
	method := os.Getenv("METHOD")
	if method == "1" {
		_, checkSSH, err := methods.ListSSHkey(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if checkSSH != vars.SshName {
			fmt.Printf("SSH key doesn't exist, creating one named: %s", vars.SshName)
			ssh, err := methods.CreateSSH(vultrClient, *vars)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("SSH key was created at %s, ID: %s, Name: %s.\n", ssh.DateCreated, ssh.ID, ssh.Name)
		} else {
			fmt.Printf("There is already a key named: %s. Please choose another name for new SSH key. \n", vars.SshName)
			fmt.Print("Continue to Private Network.....\n")
		}
		_, checkNetwork, checkNetworkSubnetMask, err := methods.ListPrivateNetwork(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if checkNetwork != vars.NetworkSubnet {
			network, err := methods.CreatePrivateNetwork(vultrClient, *vars)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Private network was created at %s, ID: %s, Network Description: %s. Range: %s/%d. \n", network.DateCreated, network.NetworkID, network.Description, network.V4Subnet, network.V4SubnetMask)
		} else {
			fmt.Printf("There is already a Private Network with the same subnet: %s/%d. Please choose another subnet for new Private network. \n", checkNetwork, checkNetworkSubnetMask)
		}

		/*instance, err := methods.CreateInstance(vultrClient, *vars, ssh.ID, network.NetworkID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Instance was created at %s , Hostname: %s , OS: %s", instance.DateCreated, instance.Label, instance.Os)*/

	}

	if method == "3" {
		sshID, sshName, err := methods.ListSSHkey(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if len(sshName) == 0 {
			fmt.Print("There is no SSH key to delete. \n")
			fmt.Print("Contuine to delete private network.... \n")
		} else if sshName != vars.SshName {
			log.Fatalf("SSH mismatch, existing SSH key: %s, proposed SSH key for deletion: %s.\n", sshName, vars.SshName)
		} else {
			err = methods.DeleteSSH(vultrClient, sshID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("SSH key deleted , ID: %s,Name: %s\n", sshID, sshName)
		}

		networkID, networkSubnet, networkSubnetMask, err := methods.ListPrivateNetwork(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if len(networkSubnet) == 0 {
			log.Fatal("There is no Private Network to delete.")
		} else if networkSubnet != vars.NetworkSubnet && networkSubnetMask != vars.NetworkSubnetMask {
			log.Fatalf("PrivateNetwork mismatch, existing PrivateNetwork: %s,proposed PrivateNetwork for deletion: %s.\n", networkSubnet, vars.NetworkSubnet)
		} else {
			err = methods.DeletePrivateNetwork(vultrClient, networkID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Private Network deleted, ID: %s, Subnet: %s\n", networkID, networkSubnet)
		}
	}
}
