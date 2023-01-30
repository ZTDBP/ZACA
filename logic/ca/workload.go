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

package ca

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ztdbp/ZACA/database/mysql/cfssl-model/dao"
	"github.com/ztdbp/ZACA/database/mysql/cfssl-model/model"
	"github.com/ztdbp/ZACA/util"
)

const AllCertsCacheKey = "all_certs_cache"

// WorkloadUnit UniqueID Divided workload unit
type WorkloadUnit struct {
	Role          string    `json:"role"`
	ValidNum      int       `json:"valid_num"`       // Number of valid certificates
	FirstIssuedAt time.Time `json:"first_issued_at"` // Date of first issuance of certificate
	UniqueId      string    `json:"unique_id"`
	Forbidden     bool      `json:"forbidden"` // Is it prohibited
}

type WorkloadUnitsParams struct {
	Page, PageSize int
	UniqueId       string
}

func getCerts(db *gorm.DB) ([]*model.Certificates, error) {
	var certs []*model.Certificates
	var err error
	allCerts, ok := util.MapCache.Get(AllCertsCacheKey)
	if !ok {
		certs, _, err = dao.GetAllCertificates(db, 1, 10000, "issued_at desc")
		if err != nil {
			return nil, errors.Wrap(err, "Database query error")
		}
		util.MapCache.SetDefault(AllCertsCacheKey, certs)
	}
	if allCerts != nil {
		certs = allCerts.([]*model.Certificates)
	}
	return certs, nil
}
