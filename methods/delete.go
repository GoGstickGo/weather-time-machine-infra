package methods

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

func DeleteSSH(vultrClient *govultr.Client, sshID string) (err error) {
	err = vultrClient.SSHKey.Delete(context.Background(), sshID)
	if err != nil {
		return fmt.Errorf("error with vultr API(DS): %v", err)
	}
	return err
}

func DeletePrivateNetwork(vultrClient *govultr.Client, networkID string) (err error) {
	err = vultrClient.Network.Delete(context.Background(), networkID)
	if err != nil {
		return fmt.Errorf("error with vultr API(DPNetwork): %v", err)
	}
	return err
}

func DeleteInstance(vultrClient *govultr.Client, instanceID string) (err error) {
	err = vultrClient.Instance.Delete(context.Background(), instanceID)
	if err != nil {
		return fmt.Errorf("error with vultr API(DI): %v", err)
	}
	return err
}
