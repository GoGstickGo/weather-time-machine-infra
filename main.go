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

type returnID struct {
	ssh     string
	network string
}

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

	i, err := methods.ReadYaml("infra.yml")
	if err != nil {
		log.Fatal(err)
	}

	var id returnID
	if i.Method == "1" {
		_, checkSSH, err := methods.ListSSHkey(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if checkSSH != i.SSH.Name {
			fmt.Printf("SSH key doesn't exist, creating one named: %s", i.SSH.Name)
			ssh, err := methods.CreateSSH(vultrClient, i)
			if err != nil {
				log.Fatal(err)
			}
			id = returnID{
				ssh: ssh.ID,
			}
			fmt.Printf("SSH key was created at %s, ID: %s, Name: %s.\n", ssh.DateCreated, ssh.ID, ssh.Name)
		} else {
			fmt.Printf("There is already a key named: %s. Please choose another name for new SSH key. \n", i.SSH.Name)
			fmt.Print("Continue to Private Network.....\n")
		}
		_, checkNetwork, checkNetworkSubnetMask, err := methods.ListPrivateNetwork(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if checkNetwork != os.Getenv("WTM_SUBNET") {
			network, err := methods.CreatePrivateNetwork(vultrClient, i)
			if err != nil {
				log.Fatal(err)
			}
			id = returnID{
				network: network.NetworkID,
			}
			fmt.Printf("Private network was created at %s, ID: %s, Network Description: %s. Range: %s/%d. \n", network.DateCreated, network.NetworkID, network.Description, network.V4Subnet, network.V4SubnetMask)
		} else {
			fmt.Printf("There is already a Private Network with the same subnet: %s/%d. Please choose another subnet for new Private network. \n", checkNetwork, checkNetworkSubnetMask)
		}

		instance, err := methods.CreateInstance(vultrClient, i, id.ssh, id.network)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Instance was created at %s , Hostname: %s , OS: %s", instance.DateCreated, instance.Label, instance.Os)

	}
	if i.Method == "2" {
		sshID, sSHName, err := methods.ListSSHkey(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if len(sSHName) == 0 {
			fmt.Print("There is no SSH key to delete. \n")
			fmt.Print("Continue to delete private network.... \n")
		} else if sSHName != i.SSH.Name {
			log.Fatalf("SSH mismatch, existing SSH key: %s, proposed SSH key for deletion: %s.\n", sSHName, i.SSH.Name)
		} else {
			err = methods.DeleteSSH(vultrClient, sshID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("SSH key deleted , ID: %s,Name: %s\n", sshID, sSHName)
		}

		networkID, networkSubnet, networkSubnetMask, err := methods.ListPrivateNetwork(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		subnetMask := fmt.Sprintf("%d", networkSubnetMask)
		if len(networkSubnet) == 0 {
			fmt.Print("There is no Private Network to delete.\n")
		} else if networkSubnet != os.Getenv("WTM_SUBNET") && subnetMask != os.Getenv("WTM_SUBNET_MASK") {
			log.Fatalf("PrivateNetwork mismatch, existing PrivateNetwork: %s,proposed PrivateNetwork for deletion: %s.\n", networkSubnet, os.Getenv("WTM_SUBNET"))
		} else {
			err = methods.DeletePrivateNetwork(vultrClient, networkID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Private Network deleted, ID: %s, Subnet: %s\n", networkID, networkSubnet)
		}
		instanceID, instanceTag, err := methods.ListInstance(vultrClient)
		if err != nil {
			log.Fatal(err)
		}
		if len(instanceTag) == 0 {
			log.Fatalf("There is no Instance with the define tag: %s", instanceTag)
		} else if instanceTag != i.Instance.Tag {
			log.Fatalf("Instance mismatch, existing Instance: %s,proposed Instance for deletion :%s.\n", i.Instance.Tag, instanceTag)
		} else {
			err = methods.DeleteInstance(vultrClient, instanceID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Instance deleted, ID: %s ,Tag %s\n", instanceID, instanceTag)
		}

	}
}
