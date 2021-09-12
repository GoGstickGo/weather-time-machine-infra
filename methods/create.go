package methods

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

type Input struct {
	SshName            string
	SshPubKey          string
	Region             string
	NetworkDescription string
	NetworkSubnet      string
	NetworkSubnetMask  int
	InstanceLabel      string
	InstanceHostname   string
	InstancePlan       string
	InstanceOSId       int
}

func CreateSSH(vultrClient *govultr.Client, i Input) (resSSH *govultr.SSHKey, err error) {
	sshkey := &govultr.SSHKeyReq{
		Name:   i.SshName,
		SSHKey: i.SshPubKey,
	}

	resSSH, err = vultrClient.SSHKey.Create(context.Background(), sshkey)
	if err != nil {
		return nil, fmt.Errorf("error with vultr API(CreateSSH): %v", err)
	}
	return resSSH, nil
}

func CreatePrivateNetwork(vultrClient *govultr.Client, i Input) (resNetwork *govultr.Network, err error) {

	privateNetwork := &govultr.NetworkReq{
		Region:       i.Region,
		Description:  i.NetworkDescription,
		V4Subnet:     i.NetworkSubnet,
		V4SubnetMask: i.NetworkSubnetMask,
	}

	resNetwork, err = vultrClient.Network.Create(context.Background(), privateNetwork)
	if err != nil {
		return nil, fmt.Errorf("error with vultr API(CreateNetwork): %v", err)
	}

	return resNetwork, nil
}

func CreateInstance(vultrClient *govultr.Client, i Input, networkID, sshID string) (resNetwork *govultr.Instance, err error) {
	instanceOptions := &govultr.InstanceCreateReq{
		Label:                i.InstanceLabel,
		Hostname:             i.InstanceHostname,
		Region:               i.Region,
		Plan:                 i.InstancePlan,
		OsID:                 i.InstanceOSId,
		AttachPrivateNetwork: []string{networkID},
		SSHKeys:              []string{sshID},
	}
	resInstance, err := vultrClient.Instance.Create(context.Background(), instanceOptions)

	if err != nil {
		return nil, fmt.Errorf("error with vultr API: %v", err)
	}
	return resInstance, nil
}
