package devices

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

const (
	HeaderByteSize             = 9
	ContentLengthPadding       = 8
	RedirectSessionLengthBytes = 13
	RedirectionSessionReply    = 4
)

type DeviceConnection struct {
	Conn          *websocket.Conn
	wsmanMessages wsman.Messages
	Device        dto.Device
	Direct        bool
	Mode          string
	Challenge     client.AuthChallenge
}

func (uc *UseCase) Redirect(c context.Context, conn *websocket.Conn, guid, mode string) error {
	// grab device info from db
	device, err := uc.GetByID(c, guid, "")
	if err != nil {
		return err
	}

	key := device.GUID + "-" + mode
	// setup wsman messages with support for talking on 16994 over tcp
	var deviceConnection *DeviceConnection
	if _, ok := uc.redirConnections[key]; ok {
		deviceConnection = uc.redirConnections[key]
	} else {
		wsmanConnection := uc.redirection.SetupWsmanClient(*device, true, true)
		deviceConnection = &DeviceConnection{
			Conn:          conn,
			wsmanMessages: wsmanConnection,
			Device:        *device,
			Direct:        false,
			Mode:          mode,
			Challenge: client.AuthChallenge{
				Username: device.Username,
				Password: device.Password,
			},
		}
		uc.redirConnections[key] = deviceConnection
	}

	err = uc.redirection.RedirectConnect(c, deviceConnection)
	if err != nil {
		return err
	}

	// To Do: scoop the errors out of this for logging
	go uc.listenToDevice(c, deviceConnection)
	go uc.listenToBrowser(c, deviceConnection)

	return nil
}

func (uc *UseCase) listenToDevice(c context.Context, deviceConnection *DeviceConnection) {
	for {
		// setup listener for response from device
		data, err := uc.redirection.RedirectListen(c, deviceConnection) // calls Receive()
		if err != nil {
			break
		}

		if len(data) == 0 {
			continue
		}

		toSend := data
		if !deviceConnection.Direct {
			toSend, deviceConnection.Direct = processDeviceData(toSend, &deviceConnection.Challenge)
		}
		// Write message back to browser
		err = deviceConnection.Conn.WriteMessage(websocket.BinaryMessage, toSend)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				_ = fmt.Errorf("interceptor - listenToDevice - websocket closed unexpectedly (writing to browser): %w", err)

				uc.redirection.RedirectClose(c, deviceConnection)
				delete(uc.redirConnections, deviceConnection.Device.GUID+"-"+deviceConnection.Mode)
			}

			return
		}
	}
}

func (uc *UseCase) listenToBrowser(c context.Context, deviceConnection *DeviceConnection) {
	for {
		_, msg, err := deviceConnection.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				_ = fmt.Errorf("interceptor - listenToBrowser - websocket closed unexpectedly (reading from browser): %w", err)

				uc.redirection.RedirectClose(c, deviceConnection)
				delete(uc.redirConnections, deviceConnection.Device.GUID+"-"+deviceConnection.Mode)
			}

			break
		}

		toSend := msg
		if !deviceConnection.Direct {
			toSend = processBrowserData(msg, &deviceConnection.Challenge)
		}
		// Send the message to the TCP Connection on the device
		err = uc.redirection.RedirectSend(c, deviceConnection, toSend) // calls send
		if err != nil {
			_ = fmt.Errorf("interceptor - listenToBrowser - error sending message to device: %w", err)
		}
	}
}

func processBrowserData(msg []byte, challenge *client.AuthChallenge) []byte {
	switch msg[0] {
	case RedirectionCommandsStartRedirectionSession:
		return msg[0:8]
	case RedirectionCommandsEndRedirectionSession:
		return msg[0:4]
	case RedirectionCommandsAuthenticateSession:
		return handleAuthenticationSession(msg, challenge)
	default:
	}

	return nil
}

func processDeviceData(msg []byte, challenge *client.AuthChallenge) ([]byte, bool) {
	switch msg[0] {
	case RedirectionCommandsStartRedirectionSessionReply:
		return handleStartRedirectionSessionReply(msg), false
	case RedirectionCommandsAuthenticateSessionReply:
		return handleAuthenticateSessionReply(msg, challenge)
	default:
	}

	return nil, false
}

func handleStartRedirectionSessionReply(msg []byte) []byte {
	if len(msg) < RedirectionSessionReply {
		return []byte("")
	}

	if msg[1:2][0] == uint8(0) {
		if len(msg) < RedirectSessionLengthBytes {
			return []byte("")
		}

		oemLen := int(msg[12:13][0])
		if len(msg) < RedirectSessionLengthBytes+oemLen {
			return []byte("")
		}

		r := msg[0 : RedirectSessionLengthBytes+oemLen]

		return r
	}

	return []byte("")
}

func allZero(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}

	return true
}

func handleAuthenticateSessionReply(msg []byte, challenge *client.AuthChallenge) ([]byte, bool) {
	if len(msg) < HeaderByteSize {
		return []byte(""), false
	}

	buf := bytes.NewReader(msg[1:HeaderByteSize])

	var authStatus, authType uint8

	var unknown uint16

	var num uint32

	_ = binary.Read(buf, binary.LittleEndian, &authStatus)
	_ = binary.Read(buf, binary.LittleEndian, &unknown)
	_ = binary.Read(buf, binary.LittleEndian, &authType)
	_ = binary.Read(buf, binary.LittleEndian, &num)

	if len(msg) < HeaderByteSize+int(num) {
		return []byte(""), false
	}

	if authType == AuthenticationTypeDigest && authStatus == AuthenticationStatusFail {
		var realmLength, nonceLength, qopLength uint8

		buf2 := bytes.NewReader(msg[9:])

		_ = binary.Read(buf2, binary.LittleEndian, &realmLength)

		realm := make([]byte, realmLength)
		_ = binary.Read(buf2, binary.LittleEndian, &realm)
		_ = binary.Read(buf2, binary.LittleEndian, &nonceLength)

		nonce := make([]byte, nonceLength)
		_ = binary.Read(buf2, binary.LittleEndian, &nonce)

		_ = binary.Read(buf2, binary.LittleEndian, &qopLength)

		qop := make([]byte, qopLength)
		_ = binary.Read(buf2, binary.LittleEndian, &qop)

		challenge.Realm = string(realm)
		challenge.Nonce = string(nonce)
		challenge.Qop = string(qop)
	} else if authType != AuthenticationTypeQuery && authStatus == AuthenticationStatusSuccess {
		// Intel AMT relayed that authentication was successful, go to direct relay mode in both directions.
		return msg, true
	}

	return msg, false
}

// RandomValueHex generates a random hexadecimal string of the specified length.
func RandomValueHex(length int) (string, error) {
	divideByHalf := 2
	n := (length + 1) / divideByHalf // Calculate the number of bytes needed

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err // Return the error if random byte generation fails
	}

	hexStr := hex.EncodeToString(b) // Convert bytes to a hexadecimal string

	return hexStr[:length], nil // Slice the string to the desired length and return it
}

// Helper function to write length and bytes.
func writeField(buf *bytes.Buffer, field string) {
	if err := binary.Write(buf, binary.BigEndian, uint8(len(field))); err != nil {
		log.Fatal(err)
	}

	if err := binary.Write(buf, binary.BigEndian, []byte(field)); err != nil {
		log.Fatal(err)
	}
}

func handleAuthenticationSession(msg []byte, challenge *client.AuthChallenge) []byte {
	if len(msg) < HeaderByteSize {
		return []byte("")
	}

	if len(msg) == 9 && allZero(msg[1:]) {
		return msg
	}

	buf := bytes.NewReader(msg[1:9])
	// Variable to hold the decoded value
	var status uint8

	var unknown uint16

	var authType uint8

	// Read the binary data into the variable
	_ = binary.Read(buf, binary.BigEndian, &status)
	_ = binary.Read(buf, binary.BigEndian, &unknown)
	_ = binary.Read(buf, binary.BigEndian, &authType)
	// generate auth challenge

	authURL := "/RedirectionService"

	if authType == AuthenticationTypeDigest {
		if challenge.Realm != "" {
			nc := challenge.NonceCount
			randomByteCount := 10
			challenge.CNonce, _ = RandomValueHex(randomByteCount)
			nonceCount := fmt.Sprintf("%08x", nc)
			nonceData := challenge.GetFormattedNonceData(challenge.Nonce)
			response := challenge.ComputeDigestHash("POST", authURL, nonceData)
			challenge.NonceCount++

			var replyBuf bytes.Buffer

			_ = binary.Write(&replyBuf, binary.BigEndian, [5]byte{0x13, 0x00, 0x00, 0x00, 0x04})                                                                                                                                                 //            [5]byte
			_ = binary.Write(&replyBuf, binary.LittleEndian, uint32(len(challenge.Username)+len(challenge.Realm)+len(challenge.Nonce)+len(authURL)+len(challenge.CNonce)+len(nonceCount)+len(response)+len(challenge.Qop)+ContentLengthPadding)) //     uint32

			// Write fields
			writeField(&replyBuf, challenge.Username)
			writeField(&replyBuf, challenge.Realm)
			writeField(&replyBuf, challenge.Nonce)
			writeField(&replyBuf, authURL)
			writeField(&replyBuf, challenge.CNonce)
			writeField(&replyBuf, nonceCount)
			writeField(&replyBuf, response)
			writeField(&replyBuf, challenge.Qop)

			return replyBuf.Bytes()
		}

		return generateEmptyAuth(challenge, authURL)
	}

	return []byte("")
}

func generateEmptyAuth(challenge *client.AuthChallenge, authURL string) []byte {
	var buf bytes.Buffer

	emptyAuth := emptyAuth{
		usernameLength: uint8(len(challenge.Username)),
		authURLPadding: [2]byte{0x00, 0x00},
		authURLLength:  uint8(len(authURL)),
		endPadding:     [4]byte{0x00, 0x00, 0x00, 0x00},
	}

	copy(emptyAuth.username[:], challenge.Username)
	copy(emptyAuth.authURL[:], authURL)

	_ = binary.Write(&buf, binary.BigEndian, [5]byte{0x13, 0x00, 0x00, 0x00, 0x04})                                // header
	_ = binary.Write(&buf, binary.LittleEndian, uint32(len(challenge.Username)+len(authURL)+ContentLengthPadding)) // flip flop endian for content length
	_ = binary.Write(&buf, binary.BigEndian, emptyAuth)

	return buf.Bytes()
}

type emptyAuth struct {
	usernameLength uint8
	username       [5]byte
	authURLPadding [2]byte
	authURLLength  uint8
	authURL        [19]byte
	endPadding     [4]byte
}
