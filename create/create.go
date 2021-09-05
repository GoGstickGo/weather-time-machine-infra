package create

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vultr/govultr/v2"
)

func CreateSSH(vultrClient *govultr.Client) string {
	sshkey := &govultr.SSHKeyReq{
		Name:   "wtm-key-01",
		SSHKey: os.Getenv("VUTLR_PUB_KEY"),
	}

	resSSH, err := vultrClient.SSHKey.Create(context.Background(), sshkey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SSH key id: %v\n", resSSH.ID)
	return resSSH.ID
}

func CreatePrivateNetwork(vultrClient *govultr.Client) string {
	privateNetwork := &govultr.NetworkReq{
		Region:       "ams",
		Description:  "WTM private network ",
		V4Subnet:     os.Getenv("VULTR_PRIVATE_NETWORK_RANGE"),
		V4SubnetMask: 23,
	}

	resNetwork, err := vultrClient.Network.Create(context.Background(), privateNetwork)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PrivateNetwork ID: %s\n", resNetwork.NetworkID)
	return resNetwork.NetworkID
}

func CreateInstance(vultrClient *govultr.Client, networkID, sshID string) *govultr.Instance {
	instanceOptions := &govultr.InstanceCreateReq{
		Label:                "wtm-01",
		Hostname:             "wtm-01",
		Region:               "ams",
		Plan:                 "vc2-1c-1gb",
		OsID:                 445,
		AttachPrivateNetwork: []string{networkID},
		SSHKeys:              []string{sshID},
	}
	resInstance, err := vultrClient.Instance.Create(context.Background(), instanceOptions)

	if err != nil {
		log.Fatal(err)
	}
	return resInstance
}
