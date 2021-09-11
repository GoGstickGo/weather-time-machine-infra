package methods

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

func DeleteSSH(vultrClient *govultr.Client, sshID string) (err error) {
	err = vultrClient.SSHKey.Delete(context.Background(), sshID)
	if err != nil {
		return fmt.Errorf("error with vultr API(DeleteSSH): %v", err)
	}
	return err
}

func DeletePrivateNetwork(vultrClient *govultr.Client, networkID string) (err error) {
	err = vultrClient.Network.Delete(context.Background(), networkID)
	if err != nil {
		return fmt.Errorf("error with vultr API(DeletePrivateNetwork): %v", err)
	}
	return err
}
