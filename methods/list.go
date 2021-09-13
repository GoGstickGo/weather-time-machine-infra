package methods

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

type vultReturn struct {
	sshId             string
	sshName           string
	networkID         string
	networkSubnet     string
	networkSubnetMask int
	instanceID        string
	instanceTag       string
}

func ListSSHkey(vultrClient *govultr.Client) (id, name string, err error) {
	keys, _, err := vultrClient.SSHKey.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		return "", "", fmt.Errorf("error with vultr API(LS): %v", err)
	}
	var sshL vultReturn
	for _, v := range keys {
		sshL = vultReturn{
			sshId:   v.ID,
			sshName: v.Name,
		}
	}
	return sshL.sshId, sshL.sshName, nil
}

func ListPrivateNetwork(vultrClient *govultr.Client) (networkID, networkSubnet string, networkSubnetMask int, err error) {
	privateNetworks, _, err := vultrClient.Network.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		return "", "", 0, fmt.Errorf("error with vultr API(LPN): %v", err)
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

func ListInstance(vultrClient *govultr.Client) (instanceID, instanceTag string, err error) {
	instances, _, err := vultrClient.Instance.List(context.Background(), &govultr.ListOptions{PerPage: 1})
	if err != nil {
		return "", "", fmt.Errorf("error with vultr API(LI): %v", err)
	}

	var instanceList vultReturn
	for _, v := range instances {
		instanceList = vultReturn{
			instanceID:  v.ID,
			instanceTag: v.Tag,
		}
	}
	return instanceList.instanceID, instanceList.instanceTag, nil
}
