package methods

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/vultr/govultr/v2"
)

func CreateSSH(vultrClient *govultr.Client, i Infra) (resSSH *govultr.SSHKey, err error) {
	sshkey := &govultr.SSHKeyReq{
		Name:   i.SSH.Name,
		SSHKey: os.Getenv("WTM_PUB_KEY"),
	}

	resSSH, err = vultrClient.SSHKey.Create(context.Background(), sshkey)
	if err != nil {
		return nil, fmt.Errorf("error with vultr API(CreateSSH): %v", err)
	}
	return resSSH, nil
}

func CreatePrivateNetwork(vultrClient *govultr.Client, i Infra) (resNetwork *govultr.Network, err error) {
	//convert env variable to int
	num, _ := strconv.Atoi(os.Getenv("WTM_SUBNET_MASK"))

	privateNetwork := &govultr.NetworkReq{
		Region:       i.Region,
		Description:  i.Network.Description,
		V4Subnet:     os.Getenv("WTM_SUBNET"),
		V4SubnetMask: num,
	}

	resNetwork, err = vultrClient.Network.Create(context.Background(), privateNetwork)
	if err != nil {
		return nil, fmt.Errorf("error with vultr API(CreateNetwork): %v", err)
	}

	return resNetwork, nil
}

func CreateInstance(vultrClient *govultr.Client, i Infra, networkID, sshID string) (resInstance *govultr.Instance, err error) {
	instanceOptions := &govultr.InstanceCreateReq{
		Label:                i.Instance.Label,
		Hostname:             i.Instance.Hostname,
		Region:               i.Region,
		Plan:                 i.Instance.Plan,
		OsID:                 i.Instance.OSId,
		AttachPrivateNetwork: []string{networkID},
		SSHKeys:              []string{sshID},
		Tag:                  i.Instance.Tag,
	}
	resInstance, err = vultrClient.Instance.Create(context.Background(), instanceOptions)

	if err != nil {
		return nil, fmt.Errorf("error with vultr API: %v", err)
	}
	return resInstance, nil
}

func CreateVKE(vultrClient *govultr.Client, i Infra) (vkeCluster *govultr.Cluster, err error) {
	clusterOptions := &govultr.ClusterReq{
		Label:     i.Kubernetes.Label,
		Version:   i.Kubernetes.Version,
		Region:    i.Kubernetes.Region,
		NodePools: i.Kubernetes.NodePool,
		//, fields{mapping: Mapping{tempField: testSliceCity3}
	}

	vkeCluster, err = vultrClient.Kubernetes.CreateCluster(context.Background(), clusterOptions)
	if err != nil {
		return nil, fmt.Errorf("error with vultr API: %v", err)
	}

	return vkeCluster, nil

}
