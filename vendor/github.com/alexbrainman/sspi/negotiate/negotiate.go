// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package negotiate provides access to the Microsoft Negotiate SSP Package.
//
package negotiate

import (
	"errors"
	"syscall"
	"time"
	"unsafe"

	"github.com/alexbrainman/sspi"
)

// Tastra: maybe (if possible) move all winapi related out of sspi and into sspi/internal/winapi

// PackageInfo contains Negotiate SSP package description.
var PackageInfo *sspi.PackageInfo

func init() {
	var err error
	PackageInfo, err = sspi.QueryPackageInfo(sspi.NEGOSSP_NAME)
	if err != nil {
		panic("failed to fetch Negotiate package info: " + err.Error())
	}
}

func acquireCredentials(principalName string, creduse uint32, ai *sspi.SEC_WINNT_AUTH_IDENTITY) (*sspi.Credentials, error) {
	c, err := sspi.AcquireCredentials(principalName, sspi.NEGOSSP_NAME, creduse, (*byte)(unsafe.Pointer(ai)))
	if err != nil {
		return nil, err
	}
	return c, nil
}

// AcquireCurrentUserCredentials acquires credentials of currently
// logged on user. These will be used by the client to authenticate
// itself to the server. It will also be used by the server
// to impersonate the user.
func AcquireCurrentUserCredentials() (*sspi.Credentials, error) {
	return acquireCredentials("", sspi.SECPKG_CRED_OUTBOUND, nil)
}

// Tastra: see if I can share this common ntlm and negotiate code

// AcquireUserCredentials acquires credentials of user described by
// domain, username and password. These will be used by the client to
// authenticate itself to the server. It will also be used by the
// server to impersonate the user.
func AcquireUserCredentials(domain, username, password string) (*sspi.Credentials, error) {
	if len(username) == 0 {
		return nil, errors.New("username parameter cannot be empty")
	}
	d, err := syscall.UTF16FromString(domain)
	if err != nil {
		return nil, err
	}
	u, err := syscall.UTF16FromString(username)
	if err != nil {
		return nil, err
	}
	p, err := syscall.UTF16FromString(password)
	if err != nil {
		return nil, err
	}
	ai := sspi.SEC_WINNT_AUTH_IDENTITY{
		User:           &u[0],
		UserLength:     uint32(len(u) - 1), // do not count terminating 0
		Domain:         &d[0],
		DomainLength:   uint32(len(d) - 1), // do not count terminating 0
		Password:       &p[0],
		PasswordLength: uint32(len(p) - 1), // do not count terminating 0
		Flags:          sspi.SEC_WINNT_AUTH_IDENTITY_UNICODE,
	}
	return acquireCredentials("", sspi.SECPKG_CRED_OUTBOUND, &ai)
}

// AcquireServerCredentials acquires server credentials that will
// be used to authenticate clients.
// The principalName parameter is passed to the underlying call to
// the winapi AcquireCredentialsHandle function (and specifies the
// name of the principal whose credentials the underlying handle
// will reference).
// As a special case, using an empty string for the principal name
// will require the credential of the user under whose security context
// the current process is running.
func AcquireServerCredentials(principalName string) (*sspi.Credentials, error) {
	return acquireCredentials(principalName, sspi.SECPKG_CRED_INBOUND, nil)
}

func updateContext(c *sspi.Context, dst, src []byte, targetName *uint16) (authCompleted bool, n int, err error) {
	var inBuf, outBuf [1]sspi.SecBuffer
	inBuf[0].Set(sspi.SECBUFFER_TOKEN, src)
	inBufs := &sspi.SecBufferDesc{
		Version:      sspi.SECBUFFER_VERSION,
		BuffersCount: 1,
		Buffers:      &inBuf[0],
	}
	outBuf[0].Set(sspi.SECBUFFER_TOKEN, dst)
	outBufs := &sspi.SecBufferDesc{
		Version:      sspi.SECBUFFER_VERSION,
		BuffersCount: 1,
		Buffers:      &outBuf[0],
	}
	ret := c.Update(targetName, outBufs, inBufs)
	switch ret {
	case sspi.SEC_E_OK:
		// session established -> return success
		return true, int(outBuf[0].BufferSize), nil
	case sspi.SEC_I_COMPLETE_NEEDED, sspi.SEC_I_COMPLETE_AND_CONTINUE:
		ret = sspi.CompleteAuthToken(c.Handle, outBufs)
		if ret != sspi.SEC_E_OK {
			return false, 0, ret
		}
	case sspi.SEC_I_CONTINUE_NEEDED:
	default:
		return false, 0, ret
	}
	return false, int(outBuf[0].BufferSize), nil
}

func makeSignature(c *sspi.Context, msg []byte, qop, seqno uint32) ([]byte, error) {
	_, maxSignature, _, _, err := c.Sizes()
	if err != nil {
		return nil, err
	}

	if maxSignature == 0 {
		return nil, errors.New("integrity services are not requested or unavailable")
	}

	var b [2]sspi.SecBuffer
	b[0].Set(sspi.SECBUFFER_DATA, msg)
	b[1].Set(sspi.SECBUFFER_TOKEN, make([]byte, maxSignature))

	ret := sspi.MakeSignature(c.Handle, qop, sspi.NewSecBufferDesc(b[:]), seqno)
	if ret != sspi.SEC_E_OK {
		return nil, ret
	}

	return b[1].Bytes(), nil
}

func encryptMessage(c *sspi.Context, msg []byte, qop, seqno uint32) ([]byte, error) {
	_ /*maxToken*/, maxSignature, cBlockSize, cSecurityTrailer, err := c.Sizes()
	if err != nil {
		return nil, err
	}

	if maxSignature == 0 {
		return nil, errors.New("integrity services are not requested or unavailable")
	}

	var b [3]sspi.SecBuffer
	b[0].Set(sspi.SECBUFFER_TOKEN, make([]byte, cSecurityTrailer))
	b[1].Set(sspi.SECBUFFER_DATA, msg)
	b[2].Set(sspi.SECBUFFER_PADDING, make([]byte, cBlockSize))

	ret := sspi.EncryptMessage(c.Handle, qop, sspi.NewSecBufferDesc(b[:]), seqno)
	if ret != sspi.SEC_E_OK {
		return nil, ret
	}

	r0, r1, r2 := b[0].Bytes(), b[1].Bytes(), b[2].Bytes()
	res := make([]byte, 0, len(r0)+len(r1)+len(r2))
	res = append(res, r0...)
	res = append(res, r1...)
	res = append(res, r2...)

	return res, nil
}

func decryptMessage(c *sspi.Context, msg []byte, seqno uint32) (uint32, []byte, error) {
	var b [2]sspi.SecBuffer
	b[0].Set(sspi.SECBUFFER_STREAM, msg)
	b[1].Set(sspi.SECBUFFER_DATA, []byte{})

	var qop uint32
	ret := sspi.DecryptMessage(c.Handle, sspi.NewSecBufferDesc(b[:]), seqno, &qop)
	if ret != sspi.SEC_E_OK {
		return qop, nil, ret
	}

	return qop, b[1].Bytes(), nil
}

func verifySignature(c *sspi.Context, msg, token []byte, seqno uint32) (uint32, error) {
	var b [2]sspi.SecBuffer
	b[0].Set(sspi.SECBUFFER_DATA, msg)
	b[1].Set(sspi.SECBUFFER_TOKEN, token)

	var qop uint32

	ret := sspi.VerifySignature(c.Handle, sspi.NewSecBufferDesc(b[:]), seqno, &qop)
	if ret != sspi.SEC_E_OK {
		return 0, ret
	}

	return qop, nil
}

// ClientContext is used by the client to manage all steps of Negotiate negotiation.
type ClientContext struct {
	sctxt      *sspi.Context
	targetName *uint16
}

// NewClientContext creates a new client context. It uses client
// credentials cred generated by AcquireCurrentUserCredentials or
// AcquireUserCredentials and SPN to start a client Negotiate
// negotiation sequence. targetName is the service principal name
// (SPN) or the security context of the destination server.
// NewClientContext returns a new token to be sent to the server.
func NewClientContext(cred *sspi.Credentials, targetName string) (cc *ClientContext, outputToken []byte, err error) {
	return NewClientContextWithFlags(cred, targetName, sspi.ISC_REQ_CONNECTION)
}

// NewClientContextWithFlags creates a new client context. It uses client
// credentials cred generated by AcquireCurrentUserCredentials or
// AcquireUserCredentials and SPN to start a client Negotiate
// negotiation sequence. targetName is the service principal name
// (SPN) or the security context of the destination server.
// The flags parameter is used to indicate requests for the context
// (for example sspi.ISC_REQ_CONFIDENTIALITY|sspi.ISC_REQ_REPLAY_DETECT)
// NewClientContextWithFlags returns a new token to be sent to the server.
func NewClientContextWithFlags(cred *sspi.Credentials, targetName string, flags uint32) (cc *ClientContext, outputToken []byte, err error) {
	var tname *uint16
	if len(targetName) > 0 {
		p, err2 := syscall.UTF16FromString(targetName)
		if err2 != nil {
			return nil, nil, err2
		}
		if len(p) > 0 {
			tname = &p[0]
		}
	}
	otoken := make([]byte, PackageInfo.MaxToken)
	c := sspi.NewClientContext(cred, flags)

	authCompleted, n, err2 := updateContext(c, otoken, nil, tname)
	if err2 != nil {
		return nil, nil, err2
	}
	if authCompleted {
		c.Release()
		return nil, nil, errors.New("negotiate authentication should not be completed yet")
	}
	if n == 0 {
		c.Release()
		return nil, nil, errors.New("negotiate token should not be empty")
	}
	otoken = otoken[:n]
	return &ClientContext{sctxt: c, targetName: tname}, otoken, nil
}

// Release free up resources associated with client context c.
func (c *ClientContext) Release() error {
	if c == nil {
		return nil
	}
	return c.sctxt.Release()
}

// Expiry returns c expiry time.
func (c *ClientContext) Expiry() time.Time {
	return c.sctxt.Expiry()
}

// Update advances client part of Negotiate negotiation c. It uses
// token received from the server and returns true if client part
// of authentication is complete. It also returns new token to be
// sent to the server.
func (c *ClientContext) Update(token []byte) (authCompleted bool, outputToken []byte, err error) {
	otoken := make([]byte, PackageInfo.MaxToken)
	authDone, n, err2 := updateContext(c.sctxt, otoken, token, c.targetName)
	if err2 != nil {
		return false, nil, err2
	}
	if n == 0 && !authDone {
		return false, nil, errors.New("negotiate token should not be empty")
	}
	otoken = otoken[:n]
	return authDone, otoken, nil
}

// Sizes queries the client context for the sizes used in per-message
// functions. It returns the maximum token size used in authentication
// exchanges, the maximum signature size, the preferred integral size of
// messages, the size of any security trailer, and any error.
func (c *ClientContext) Sizes() (uint32, uint32, uint32, uint32, error) {
	return c.sctxt.Sizes()
}

// MakeSignature uses the established client context to create a signature
// for the given message using the provided quality of protection flags and
// sequence number. It returns the signature token in addition to any error.
func (c *ClientContext) MakeSignature(msg []byte, qop, seqno uint32) ([]byte, error) {
	return makeSignature(c.sctxt, msg, qop, seqno)
}

// VerifySignature uses the established client context and signature token
// to check that the provided message hasn't been tampered or received out
// of sequence. It returns any quality of protection flags and any error
// that occurred.
func (c *ClientContext) VerifySignature(msg, token []byte, seqno uint32) (uint32, error) {
	return verifySignature(c.sctxt, msg, token, seqno)
}

// EncryptMessage uses the established client context to encrypt a message
// using the provided quality of protection flags and sequence number.
// It returns the signature token in addition to any error.
// IMPORTANT: the input msg parameter is updated in place by the low-level windows api
// so must be copied if the initial content should not be modified.
func (c *ClientContext) EncryptMessage(msg []byte, qop, seqno uint32) ([]byte, error) {
	return encryptMessage(c.sctxt, msg, qop, seqno)
}

// DecryptMessage uses the established client context to decrypt a message
// using the provided sequence number.
// It returns the quality of protection flag and the decrypted message in addition to any error.
func (c *ClientContext) DecryptMessage(msg []byte, seqno uint32) (uint32, []byte, error) {
	return decryptMessage(c.sctxt, msg, seqno)
}

// VerifyFlags determines if all flags used to construct the client context
// were honored (see NewClientContextWithFlags).  It should be called after c.Update.
func (c *ClientContext) VerifyFlags() error {
	return c.sctxt.VerifyFlags()
}

// VerifySelectiveFlags determines if the given flags were honored (see NewClientContextWithFlags).
// It should be called after c.Update.
func (c *ClientContext) VerifySelectiveFlags(flags uint32) error {
	return c.sctxt.VerifySelectiveFlags(flags)
}

// ServerContext is used by the server to manage all steps of Negotiate
// negotiation. Once authentication is completed the context can be
// used to impersonate client.
type ServerContext struct {
	sctxt *sspi.Context
}

// NewServerContext creates new server context. It uses server
// credentials created by AcquireServerCredentials and token from
// the client to start server Negotiate negotiation sequence.
// It also returns new token to be sent to the client.
func NewServerContext(cred *sspi.Credentials, token []byte) (sc *ServerContext, authDone bool, outputToken []byte, err error) {
	otoken := make([]byte, PackageInfo.MaxToken)
	c := sspi.NewServerContext(cred, sspi.ASC_REQ_CONNECTION)
	authDone, n, err2 := updateContext(c, otoken, token, nil)
	if err2 != nil {
		return nil, false, nil, err2
	}
	otoken = otoken[:n]
	return &ServerContext{sctxt: c}, authDone, otoken, nil
}

// Release free up resources associated with server context c.
func (c *ServerContext) Release() error {
	if c == nil {
		return nil
	}
	return c.sctxt.Release()
}

// Expiry returns c expiry time.
func (c *ServerContext) Expiry() time.Time {
	return c.sctxt.Expiry()
}

// Update advances server part of Negotiate negotiation c. It uses
// token received from the client and returns true if server part
// of authentication is complete. It also returns new token to be
// sent to the client.
func (c *ServerContext) Update(token []byte) (authCompleted bool, outputToken []byte, err error) {
	otoken := make([]byte, PackageInfo.MaxToken)
	authDone, n, err2 := updateContext(c.sctxt, otoken, token, nil)
	if err2 != nil {
		return false, nil, err2
	}
	if n == 0 && !authDone {
		return false, nil, errors.New("negotiate token should not be empty")
	}
	otoken = otoken[:n]
	return authDone, otoken, nil
}

const _SECPKG_ATTR_NATIVE_NAMES = 13

type _SecPkgContext_NativeNames struct {
	ClientName *uint16
	ServerName *uint16
}

// GetUsername returns the username corresponding to the authenticated client
func (c *ServerContext) GetUsername() (string, error) {
	var ns _SecPkgContext_NativeNames
	ret := sspi.QueryContextAttributes(c.sctxt.Handle, _SECPKG_ATTR_NATIVE_NAMES, (*byte)(unsafe.Pointer(&ns)))
	if ret != sspi.SEC_E_OK {
		return "", ret
	}
	sspi.FreeContextBuffer((*byte)(unsafe.Pointer(ns.ServerName)))
	defer sspi.FreeContextBuffer((*byte)(unsafe.Pointer(ns.ClientName)))
	return syscall.UTF16ToString((*[2 << 20]uint16)(unsafe.Pointer(ns.ClientName))[:]), nil
}

// ImpersonateUser changes current OS thread user. New user is
// the user as specified by client credentials.
func (c *ServerContext) ImpersonateUser() error {
	return c.sctxt.ImpersonateUser()
}

// RevertToSelf stops impersonation. It changes current OS thread
// user to what it was before ImpersonateUser was executed.
func (c *ServerContext) RevertToSelf() error {
	return c.sctxt.RevertToSelf()
}

// Sizes queries the server context for the sizes used in per-message
// functions. It returns the maximum token size used in authentication
// exchanges, the maximum signature size, the preferred integral size of
// messages, the size of any security trailer, and any error.
func (c *ServerContext) Sizes() (uint32, uint32, uint32, uint32, error) {
	return c.sctxt.Sizes()
}

// MakeSignature uses the established server context to create a signature
// for the given message using the provided quality of protection flags and
// sequence number. It returns the signature token in addition to any error.
func (c *ServerContext) MakeSignature(msg []byte, qop, seqno uint32) ([]byte, error) {
	return makeSignature(c.sctxt, msg, qop, seqno)
}

// VerifySignature uses the established server context and signature token
// to check that the provided message hasn't been tampered or received out
// of sequence. It returns any quality of protection flags and any error
// that occurred.
func (c *ServerContext) VerifySignature(msg, token []byte, seqno uint32) (uint32, error) {
	return verifySignature(c.sctxt, msg, token, seqno)
}

// EncryptMessage uses the established server context to encrypt a message
// using the provided quality of protection flags and sequence number.
// It returns the signature token in addition to any error.
// IMPORTANT: the input msg parameter is updated in place by the low-level windows api
// so must be copied if the initial content should not be modified.
func (c *ServerContext) EncryptMessage(msg []byte, qop, seqno uint32) ([]byte, error) {
	return encryptMessage(c.sctxt, msg, qop, seqno)
}

// DecryptMessage uses the established server context to decrypt a message
// using the provided sequence number.
// It returns the quality of protection flag and the decrypted message in addition to any error.
func (c *ServerContext) DecryptMessage(msg []byte, seqno uint32) (uint32, []byte, error) {
	return decryptMessage(c.sctxt, msg, seqno)
}
