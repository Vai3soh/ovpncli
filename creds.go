/*
* ovpncli -- Library for wrapping openvpn3 (https://github.com/OpenVPN/openvpn3) functionality in go way.
* Copyright (C) 2022 Vai3soh
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License Version 3
* as published by the Free Software Foundation.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.

* You should have received a copy of the GNU Affero General Public License
* along with this program in the COPYING file.
* If not, see <http://www.gnu.org/licenses/>.
 */
package ovpncli

type clientCreds struct {
	ClientAPI_ProvideCreds
}

type OptionCred func(*clientCreds)

func NewClientCreds(opts ...OptionCred) *clientCreds {
	op := &clientCreds{
		ClientAPI_ProvideCreds: NewClientAPI_ProvideCreds(),
	}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func WithPassword(password string) OptionCred {
	return func(op *clientCreds) {
		op.SetPassword(password)
	}
}

func WithUsername(username string) OptionCred {
	return func(op *clientCreds) {
		op.SetUsername(username)
	}
}

func WithHttpProxyPass(password string) OptionCred {
	return func(op *clientCreds) {
		op.SetHttp_proxy_pass(password)
	}
}

func WithHttpProxyUser(user string) OptionCred {
	return func(op *clientCreds) {
		op.SetHttp_proxy_user(user)
	}
}

func WithResponse(arg string) OptionCred {
	return func(op *clientCreds) {
		op.SetResponse(arg)
	}
}

func WithCachePassword(enable bool) OptionCred {
	return func(op *clientCreds) {
		op.SetCachePassword(enable)
	}
}

func WithDynamicChallengeCookie(arg string) OptionCred {
	return func(op *clientCreds) {
		op.SetDynamicChallengeCookie(arg)
	}
}

func WithReplacePasswordWithSessionID(enable bool) OptionCred {
	return func(op *clientCreds) {
		op.SetReplacePasswordWithSessionID(enable)
	}
}
