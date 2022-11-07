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

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Client interface {
	ClientAPI_OpenVPNClient
	deleteClient()
	StartConnection(ctx context.Context)
	StopConnection()
	Reconnection(time int)
	ResumeConnect()
	PauseConnect(reason string)
	CallbackError() error
}

type client struct {
	ClientAPI_OpenVPNClient
	statusErr chan error
	g         *errgroup.Group
}

type clientConfig struct {
	ClientAPI_Config
}

func (c *client) deleteClient() {
	DeleteDirectorClientAPI_OpenVPNClient(c.ClientAPI_OpenVPNClient)
}

func (c *client) controlCancelContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.StopConnection()
			return
		default:
			continue
		}
	}
}

func (c *client) ResumeConnect() {
	c.Resume()
}

// Pause the client -- useful to avoid continuous reconnection attempts
// when network is down.  May be called from a different thread
// when connect() is running.

func (c *client) PauseConnect(reason string) {
	c.Pause(reason)
}

func (c *client) Reconnection(time int) {
	c.Reconnect(time)
}

func (c *client) StartConnection(ctx context.Context) {

	c.g.Go(func() error {
		status := c.Connect()
		if status.GetError() {
			return fmt.Errorf("error with status: [%s]", status.GetMessage())
		}
		return nil
	})

	go c.controlCancelContext(ctx)
	go c.getErrorFromSession()
}

func (c *client) getErrorFromSession() {
	c.statusErr <- c.g.Wait()
}

func (c *client) CallbackError() error {
	for {
		select {
		case err := <-c.statusErr:
			return err
		default:
			continue
		}
	}
}

func (c *client) StopConnection() {
	c.Stop()
}

func DeleteClient(c Client) {
	c.deleteClient()
}

type OverwriteClient interface{}

func NewClient(ocl OverwriteClient) Client {

	cl := NewDirectorClientAPI_OpenVPNClient(ocl)
	cli := &client{
		ClientAPI_OpenVPNClient: cl,
		statusErr:               make(chan error),
		g:                       &errgroup.Group{},
	}
	return cli
}

type Option func(*clientConfig)

func NewClientConfig(opts ...Option) *clientConfig {
	op := &clientConfig{
		ClientAPI_Config: NewClientAPI_Config(),
	}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func WithConfig(config string) Option {
	return func(op *clientConfig) {
		op.SetContent(config)
	}
}

// 1,2
func WithSslDebugLevel(level int) Option {
	return func(op *clientConfig) {
		op.SetSslDebugLevel(level)
	}
}

// yes|no|asym
func WithCompressionMode(mode string) Option {
	return func(op *clientConfig) {
		op.SetCompressionMode(mode)
	}
}

func WithConnTimeout(time int) Option {
	return func(op *clientConfig) {
		op.SetConnTimeout(time)
	}
}

func WithLegacyAlgorithms(enable bool) Option {
	return func(op *clientConfig) {
		op.SetEnableLegacyAlgorithms(enable)
	}
}

func WithNonPreferredDCAlgorithms(enable bool) Option {
	return func(op *clientConfig) {
		op.SetEnableNonPreferredDCAlgorithms(enable)
	}
}

func WithDisableClientCert(enable bool) Option {
	return func(op *clientConfig) {
		op.SetDisableClientCert(enable)
	}
}

func WithClockTickMS(timeMS uint) Option {
	return func(op *clientConfig) {
		op.SetClockTickMS(timeMS)
	}
}

func WithRetryOnAuthFailed(enable bool) Option {
	return func(op *clientConfig) {
		op.SetRetryOnAuthFailed(enable)
	}
}

func WithAllowLocalDnsResolvers(enable bool) Option {
	return func(op *clientConfig) {
		op.SetAllowLocalDnsResolvers(enable)
	}
}

func WithAllowLocalLanAccess(enable bool) Option {
	return func(op *clientConfig) {
		op.SetAllowLocalLanAccess(enable)
	}
}

func WithAllowUnusedAddrFamilies(arg string) Option {
	return func(op *clientConfig) {
		op.SetAllowUnusedAddrFamilies(arg)
	}
}

func WithAltProxy(enable bool) Option {
	return func(op *clientConfig) {
		op.SetAltProxy(enable)
	}
}

func WithAutologinSessions(enable bool) Option {
	return func(op *clientConfig) {
		op.SetAutologinSessions(enable)
	}
}

func WithDco(enable bool) Option {
	return func(op *clientConfig) {
		op.SetDco(enable)
	}
}

func WithEcho(enable bool) Option {
	return func(op *clientConfig) {
		op.SetEcho(enable)
	}
}

func WithExternalPkiAlias(arg string) Option {
	return func(op *clientConfig) {
		op.SetExternalPkiAlias(arg)
	}
}

func WithGenerateTunBuilderCaptureEvent(arg bool) Option {
	return func(op *clientConfig) {
		op.SetGenerate_tun_builder_capture_event(arg)
	}
}

func WithGoogleDnsFallback(arg bool) Option {
	return func(op *clientConfig) {
		op.SetGoogleDnsFallback(arg)
	}
}

func WithGremlinConfig(arg string) Option {
	return func(op *clientConfig) {
		op.SetGremlinConfig(arg)
	}
}

func WithGuiVersion(arg string) Option {
	return func(op *clientConfig) {
		op.SetGuiVersion(arg)
	}
}

func WithHwAddrOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetHwAddrOverride(arg)
	}
}

func WithInfo(arg bool) Option {
	return func(op *clientConfig) {
		op.SetInfo(arg)
	}
}

func WithPeerInfo(arg Std_vector_Sl_openvpn_ClientAPI_KeyValue_Sg_) Option {
	return func(op *clientConfig) {
		op.SetPeerInfo(arg)
	}
}

func WithPlatformVersion(arg string) Option {
	return func(op *clientConfig) {
		op.SetPlatformVersion(arg)
	}
}

func WithPortOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetPortOverride(arg)
	}
}

func WithPrivateKeyPassword(arg string) Option {
	return func(op *clientConfig) {
		op.SetPrivateKeyPassword(arg)
	}
}

func WithProtoOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetProtoOverride(arg)
	}
}

func WithProtoVersionOverride(arg int) Option {
	return func(op *clientConfig) {
		op.SetProtoVersionOverride(arg)
	}
}

func WithProxyAllowCleartextAuth(arg bool) Option {
	return func(op *clientConfig) {
		op.SetProxyAllowCleartextAuth(arg)
	}
}

func WithProxyHost(arg string) Option {
	return func(op *clientConfig) {
		op.SetProxyHost(arg)
	}
}

func WithProxyPassword(arg string) Option {
	return func(op *clientConfig) {
		op.SetProxyPassword(arg)
	}
}

func WithProxyPort(arg string) Option {
	return func(op *clientConfig) {
		op.SetProxyPort(arg)
	}
}

func WithProxyUsername(arg string) Option {
	return func(op *clientConfig) {
		op.SetProxyUsername(arg)
	}
}

func WithServerOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetServerOverride(arg)
	}
}

func WithSsoMethods(arg string) Option {
	return func(op *clientConfig) {
		op.SetSsoMethods(arg)
	}
}

func WithSynchronousDnsLookup(arg bool) Option {
	return func(op *clientConfig) {
		op.SetSynchronousDnsLookup(arg)
	}
}

func WithTlsCertProfileOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetTlsCertProfileOverride(arg)
	}
}

func WithTlsCipherList(arg string) Option {
	return func(op *clientConfig) {
		op.SetTlsCipherList(arg)
	}
}

func WithTlsCiphersuitesList(arg string) Option {
	return func(op *clientConfig) {
		op.SetTlsCiphersuitesList(arg)
	}
}

func WithTlsVersionMinOverride(arg string) Option {
	return func(op *clientConfig) {
		op.SetTlsVersionMinOverride(arg)
	}
}

func WithTunPersist(arg bool) Option {
	return func(op *clientConfig) {
		op.SetTunPersist(arg)
	}
}

func WithWinTun(arg bool) Option {
	return func(op *clientConfig) {
		op.SetWintun(arg)
	}
}

func WithDefaultKeyDirection(arg int) Option {
	return func(op *clientConfig) {
		op.SetDefaultKeyDirection(arg)
	}
}
