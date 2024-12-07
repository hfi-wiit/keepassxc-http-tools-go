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
)
