package redirection

import (
	"context"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

func (uc *UseCase) SetupWsmanClient(device dto.Device, isRedirection, logAMTMessages bool) wsman.Messages {
	clientParams := client.Parameters{
		Target:            device.Hostname,
		Username:          device.Username,
		Password:          device.Password,
		UseDigest:         true,
		UseTLS:            device.UseTLS,
		SelfSignedAllowed: device.AllowSelfSigned,
		LogAMTMessages:    logAMTMessages,
		IsRedirection:     isRedirection,
	}

	return wsman.NewMessages(clientParams)
}

func (uc *UseCase) RedirectConnect(_ context.Context, deviceConnection *DeviceConnection) error {
	err := deviceConnection.wsmanMessages.Client.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RedirectSend(_ context.Context, deviceConnection *DeviceConnection, data []byte) error {
	err := deviceConnection.wsmanMessages.Client.Send(data)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RedirectListen(_ context.Context, deviceConnection *DeviceConnection) ([]byte, error) {
	data, err := deviceConnection.wsmanMessages.Client.Receive()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (uc *UseCase) RedirectClose(_ context.Context, deviceConnection *DeviceConnection) error {
	err := deviceConnection.wsmanMessages.Client.CloseConnection()
	if err != nil {
		return err
	}

	return nil
}
