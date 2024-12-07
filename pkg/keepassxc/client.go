package keepassxc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/scalarmult"

	"keepassxc-http-tools-go/pkg/utils"
)

var (
	ErrUnspecifiedSocketPath = errors.New("unspecified socket path")
	ErrInvalidPeerKey        = errors.New("invalid peer key")
	ErrNotImplemented        = errors.New("not implemented yet")
)

type Client struct {
	Id              string
	SocketPath      string
	ApplicationName string
	AssocProfile    KeepassxcClientProfile

	socket     net.Conn
	privateKey nacl.Key
	publicKey  nacl.Key
	peerKey    nacl.Key
}

/*
	Connection implementation
*/

// ClientOption type represents an option function for NewClient.
type ClientOption func(*Client) error

// OptApplicationName is an option to NewClient.
// It can be used to override the default application name.
// See constant utils.ApplicationName.
func OptApplicationName(name string) ClientOption {
	return func(client *Client) error {
		client.ApplicationName = name
		return nil
	}
}

// OptSocketPath is an option to NewClient.
// It can be used to override the auto detection of the socket file.
func OptSocketPath(path string) ClientOption {
	return func(client *Client) error {
		client.SocketPath = path
		if _, err := os.Stat(client.SocketPath); err != nil {
			return errors.Join(err, utils.ErrKeepassxcSocketNotFound)
		}
		return nil
	}
}

// NewClient creates a new keepassxc http api client and connect to its socket.
func NewClient(assocProfile KeepassxcClientProfile, options ...ClientOption) (*Client, error) {
	var err error
	client := &Client{
		AssocProfile: assocProfile,
		privateKey:   nacl.NewKey(),
	}
	client.publicKey = scalarmult.Base(client.privateKey)

	for _, option := range options {
		err = option(client)
		if err != nil {
			return nil, err
		}
	}

	if client.ApplicationName == "" {
		client.ApplicationName = utils.ApplicationName
	}

	if client.SocketPath == "" {
		client.SocketPath, err = SocketPath()
		if err != nil {
			return nil, err
		}
	}

	client.Id = client.ApplicationName + utils.NaclNonceToB64(nacl.NewNonce())
	client.socket, err = connect(client.SocketPath)
	if err != nil {
		return nil, err
	}
	if err = client.exchangePublicKeys(); err != nil {
		return nil, err
	}
	if client.AssocProfile.GetAssocKey() == nil {
		if err = client.associate(); err != nil {
			return nil, err
		}
	} else {
		if err = client.testAssociate(); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// exchangePublicKeys is a helper function for NewClient.
// It exchanges encryption keys with the server.
// TODO refactor
func (c *Client) exchangePublicKeys() error {
	resp, err := c.sendMessage(Message{
		"action":    "change-public-keys",
		"publicKey": utils.NaclKeyToB64(c.publicKey),
	}, false)
	if err != nil {
		return err
	}
	if peerKey, ok := resp["publicKey"]; ok {
		c.peerKey = utils.B64ToNaclKey(peerKey.(string))
		return nil
	}
	return utils.ErrKeepassxcKeyExchangeFailed
}

// associate is a helper function for NewClient.
// It tells the server to associate the key with the given profile.
// TODO refactor
func (c *Client) associate() error {
	assocKey := nacl.NewKey()
	resp, err := c.sendMessage(Message{
		"action": "associate",
		"key":    utils.NaclKeyToB64(c.publicKey),
		"idKey":  utils.NaclKeyToB64(assocKey),
	}, true)
	if err != nil {
		return err
	}
	if v, ok := resp["message"]; ok {
		if msg, ok := v.(map[string]interface{}); ok {
			if id, ok := msg["id"]; ok {
				return c.AssocProfile.SetAssoc(id.(string), assocKey)
			}
		}
	}
	return errors.New("associate failed")
}

// testAssociate is a helper function for NewClient.
// It tests the association if an assocKey is already present in the profile.
func (c *Client) testAssociate() error {
	_, err := c.sendMessage(Message{
		"action": "test-associate",
		"key":    utils.NaclKeyToB64(c.AssocProfile.GetAssocKey()),
		"id":     c.AssocProfile.GetAssocName(),
	}, true)
	return err
}

// Disconnect from the keepassxc http api socket.
func (c *Client) Disconnect() error {
	if c.socket != nil {
		return c.socket.Close()
	}
	return nil
}

/*
	Messaging implementation
*/

func (c *Client) encryptMessage(msg Message) ([]byte, error) {
	if len(c.peerKey) == 0 {
		return nil, ErrInvalidPeerKey
	}
	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return box.EasySeal(msgData, c.peerKey, c.privateKey), nil
}

func (c *Client) decryptResponse(encryptedMsg []byte) ([]byte, error) {
	if len(c.peerKey) == 0 {
		return nil, ErrInvalidPeerKey
	}
	return box.EasyOpen(encryptedMsg, c.peerKey, c.privateKey)
}

func (c *Client) sendMessage(msg Message, encrypted bool) (Response, error) {
	if encrypted {
		encryptedMsg, err := c.encryptMessage(msg)
		if err != nil {
			return nil, err
		}
		action := msg["action"]
		msg = Message{
			"action":  action,
			"message": base64.StdEncoding.EncodeToString(encryptedMsg[nacl.NonceSize:]),
			"nonce":   base64.StdEncoding.EncodeToString(encryptedMsg[:nacl.NonceSize]),
		}
	} else {
		msg["nonce"] = utils.NaclNonceToB64(nacl.NewNonce())
	}
	msg["clientID"] = c.Id

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_, err = c.socket.Write(data)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	count, err := c.socket.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[0:count]

	var resp Response
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return nil, err
	}

	if err, ok := resp["error"]; ok {
		return nil, fmt.Errorf("%v %s", resp["errorCode"], err.(string))
	}

	if encrypted {
		decoded, err := base64.StdEncoding.DecodeString(resp["nonce"].(string) + resp["message"].(string))
		if err != nil {
			return nil, err
		}
		decryptedMsg, err := c.decryptResponse(decoded)
		if err != nil {
			return nil, err
		}
		var msg map[string]interface{}
		err = json.Unmarshal(decryptedMsg, &msg)
		if err != nil {
			return nil, err
		}
		resp["message"] = msg
	}

	return resp, err
}

func (c *Client) GetLogins(url string) (Entries, error) {
	msg := Message{
		"action": "get-logins",
		"url":    url,
		"keys": []map[string]string{
			{
				"id":  c.AssocProfile.GetAssocName(),
				"key": utils.NaclKeyToB64(c.AssocProfile.GetAssocKey()),
			},
		},
	}
	resp, err := c.sendMessage(msg, true)
	if err != nil {
		return nil, err
	}

	return resp.entries()
}
