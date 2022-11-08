/*
Copyright 2022-present The ZTDBP Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keymanager

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ztdbp/ZACA/pkg/logger"
	cfssl_client "github.com/ztdbp/cfssl/api/client"
	"github.com/ztdbp/cfssl/cli/genkey"
	"github.com/ztdbp/cfssl/csr"
	"github.com/ztdbp/cfssl/signer"

	"github.com/ztdbp/ZACA/core"
)

// RemoteSigner ...
type RemoteSigner struct {
	logger *logger.Logger
}

// NewRemoteSigner ...
func NewRemoteSigner() *RemoteSigner {
	return &RemoteSigner{
		logger: logger.Named("remote-signer"),
	}
}

//Run calls the remote CA to sign the certificate and persist it
func (ss *RemoteSigner) Run() error {
	if core.Is.Config.Keymanager.SelfSign {
		return nil
	}
	key, cert, _ := GetKeeper().GetCachedSelfKeyPairPEM()
	if key != nil && cert != nil {
		ss.logger.Info("The certificate already exists. Skip the remote signing process")
		return nil
	}
	ss.logger.Warn("There is no certificate. You will sign the certificate remotely")
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(getIntermediateCSRTemplate())
	if err != nil {
		ss.logger.Errorf("key, csr Production error: %v", err)
		return err
	}

	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Profile: "intermediate",
	}

	signReqBytes, _ := jsoniter.Marshal(&signReq)
	err = GetKeeper().RootClient.DoWithRetry(func(remote *cfssl_client.AuthRemote) error {
		certResp, err := remote.Sign(signReqBytes)
		if err != nil {
			return err
		}
		cert = certResp
		return nil
	})
	if err != nil {
		ss.logger.Errorf("initca Create error: %v", err)
		return err
	}
	ss.logger.With("key", string(key), "cert", string(cert)).Debugf("Self signed certificate completed")
	if err = GetKeeper().SetKeyPairPEM(key, cert); err != nil {
		ss.logger.Errorf("Error saving certificate: %v", err)
		return err
	}

	return nil
}
