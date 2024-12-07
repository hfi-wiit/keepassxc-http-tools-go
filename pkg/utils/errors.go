package utils

import "errors"

var (
	// keepassxc lib generic base error
	ErrKeepassxc = errors.New("keepassxc error")
	// keepassxc lib socket detection error
	ErrKeepassxcSocketNotFound = errors.Join(errors.New("keepassxc socket not found"), ErrKeepassxc)
	// keepassxc lib invalid api response error
	ErrKeepassxcInvalidResponse = errors.Join(errors.New("keepassxc invalid api response"), ErrKeepassxc)
	// keepassxc lib key exchange failed error
	ErrKeepassxcKeyExchangeFailed = errors.Join(errors.New("keepassxc key exchange failed"), ErrKeepassxc)
	// keepassxc lib profile association failed error
	ErrKeepassxcAssocFailed = errors.Join(errors.New("keepassxc profile association failed"), ErrKeepassxc)
	// keepassxc lib profile association test failed error
	ErrKeepassxcTestAssocFailed = errors.Join(errors.New("keepassxc profile association test failed"), ErrKeepassxc)
	// keepassxc lib message encryption error
	ErrKeepassxcEncryptionFailed = errors.Join(errors.New("keepassxc failed to encrypt message"), ErrKeepassxc)
	// keepassxc lib message decryption error
	ErrKeepassxcDecryptionFailed = errors.Join(errors.New("keepassxc failed to decrypt message"), ErrKeepassxc)
	// keepassxc lib send message error
	ErrKeepassxcSendMessageFailed = errors.Join(errors.New("keepassxc failed send the message"), ErrKeepassxc)
)
