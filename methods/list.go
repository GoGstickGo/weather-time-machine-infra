package methods

import (
	"context"
	"fmt"
	"log"

	"github.com/vultr/govultr/v2"
)

type vultReturn struct {
	sshId             string
	sshName           string
	networkID         string
	networkSubnet     string
	networkSubnetMask int
	instanceID        string
	instanceLabel     string
}

func ListSSHkey(vultrClient *govultr.Client) (id, name string, err error) {
	keys, _, err := vultrClient.SSHKey.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		return "", "", fmt.Errorf("error with vultr API(List SSH): %v", err)
	}
	var sshList vultReturn
	for _, v := range keys {
		sshList = vultReturn{
			sshId:   v.ID,
			sshName: v.Name,
		}
	}
	return sshList.sshId, sshList.sshName, nil
}

func ListPrivateNetwork(vultrClient *govultr.Client) (networkID, networkSubnet string, networkSubnetMask int, err error) {
	privateNetworks, _, err := vultrClient.Network.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		return "", "", 0, fmt.Errorf("error with vultr API(List PrivateNetwork): %v", err)
	}
	var networkList vultReturn
	for _, v := range privateNetworks {
		networkList = vultReturn{
			networkID:         v.NetworkID,
			networkSubnet:     v.V4Subnet,
			networkSubnetMask: v.V4SubnetMask,
		}
	}
	return networkList.networkID, networkList.networkSubnet, networkList.networkSubnetMask, nil
}

func ListInstance(vultrClient *govultr.Client) (instanceID, instanceLabel string) {
	instances, _, err := vultrClient.Instance.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		log.Fatal(err)
	}

	var instanceList vultReturn
	for _, v := range instances {
		instanceList = vultReturn{
			instanceID:    v.ID,
			instanceLabel: v.Label,
		}
	}
	return instanceList.instanceID, instanceList.instanceLabel
}
