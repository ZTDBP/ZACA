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

package revoke

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ztalab/zta-tools/pkg/logger"
	"github.com/ztdbp/cfssl/api"
	"github.com/ztdbp/cfssl/certdb"
	cf_err "github.com/ztdbp/cfssl/errors"
	"github.com/ztdbp/cfssl/helpers"
	"github.com/ztdbp/cfssl/hook"
	"github.com/ztdbp/cfssl/ocsp"
	"gorm.io/gorm"

	"github.com/ztdbp/ZACA/core"
	"github.com/ztdbp/ZACA/database/mysql/cfssl-model/model"
	"github.com/ztdbp/ZACA/logic/events"
	"github.com/ztdbp/ZACA/pkg/signature"
	"github.com/ztdbp/ZACA/util"
)

type Handler struct {
	dbAccessor certdb.Accessor
	logger     *logger.Logger
}

func NewHandler(dbAccessor certdb.Accessor) http.Handler {
	return &api.HTTPHandler{
		Handler: &Handler{
			dbAccessor: dbAccessor,
			logger:     logger.Named("revoke"),
		},
		Methods: []string{"POST"},
	}
}

type JsonRevokeRequest struct {
	Serial  string `json:"serial"`
	AKI     string `json:"authority_key_id"`
	Reason  string `json:"reason"`
	Nonce   string `json:"nonce"`
	Sign    string `json:"sign"`
	AuthKey string `json:"auth_key"`
	Profile string `json:"profile"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()

	var req JsonRevokeRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return cf_err.NewBadRequestString("Unable to parse revocation request")
	}

	if len(req.Serial) == 0 {
		return cf_err.NewBadRequestString("serial number is required but not provided")
	}

	certRecord := &model.Certificates{}
	if err := core.Is.Db.Where("serial_number = ? AND authority_key_identifier = ?", req.Serial, req.AKI).First(certRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.logger.With("sn", req.Serial, "aki", req.AKI).Warn("Certificate does not exist")
		} else {
			h.logger.With("sn", req.Serial, "aki", req.AKI).Errorf("Certificate acquisition error: %v", err)
		}
		return cf_err.NewBadRequest(err)
	}

	if hook.EnableVaultStorage {
		pem, err := core.Is.VaultSecret.GetCertPEM(req.Serial)
		if err != nil {
			h.logger.With("sn", req.Serial, "aki", req.AKI).Warnf("Vault Get error: %v", err)
		} else {
			certRecord.Pem = *pem
		}
	}

	cert, err := helpers.ParseCertificatePEM([]byte(certRecord.Pem))
	if err != nil {
		h.logger.With("sn", req.Serial, "aki", req.AKI).Errorf("Certificate PEM parsing error: %v", err)
		return cf_err.NewBadRequest(err)
	}

	var valid bool
	if req.AuthKey == "" {
		v := signature.NewVerifier(cert.PublicKey)
		valid, err = v.Verify([]byte(req.Nonce), req.Sign)
		if err != nil {
			h.logger.With("sn", req.Serial, "aki", req.AKI).Warnf("Validation error: %v", err)
			return cf_err.NewBadRequest(err)
		}
	} else {
		if req.Profile == "" {
			return cf_err.NewBadRequest(errors.New("profile Unspecified"))
		}
		if authKey, ok := core.Is.Config.Singleca.CfsslConfig.AuthKeys[req.Profile]; ok {
			if authKey.Key == req.AuthKey {
				valid = true
			}
		}
	}

	if !valid {
		h.logger.With("sn", req.Serial, "aki", req.AKI).Warnf("Certificate cannot correspond: %v", err)
		return cf_err.NewBadRequest(err)
	}

	var reasonCode int
	reasonCode, err = ocsp.ReasonStringToCode("keycompromise")
	if err != nil {
		return cf_err.NewBadRequestString("Invalid reason code")
	}

	// Delete the certificate corresponding to vault
	if hook.EnableVaultStorage {
		if err := core.Is.VaultSecret.DeleteCertPEM(req.Serial); err != nil {
			h.logger.With("sn", req.Serial, "aki", req.AKI).Warnf("Vault Delete error: %v", err)
		}
	}

	err = h.dbAccessor.RevokeCertificate(req.Serial, req.AKI, reasonCode)
	if err != nil {
		h.logger.With("sn", req.Serial, "aki", req.AKI).Warnf("Database operation error: %v", err)
		return err
	}

	AddMetricsPoint(cert)

	events.NewWorkloadLifeCycle("self-revoke", events.OperatorSDK, events.CertOp{
		UniqueId: cert.Subject.CommonName,
		SN:       req.Serial,
		AKI:      req.AKI,
	}).Log()

	h.logger.With("sn", req.Serial, "aki", req.AKI, "uri", util.GetSanURI(cert)).Info("Workload Active revocation of certificate")

	result := map[string]string{}
	return api.SendResponse(w, result)
}
