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

package model

import (
	"time"
)

// OcspResponses struct is a row record of the ocsp_responses table in the cap database
type OcspResponses struct {
	SerialNumber           string    `gorm:"primary_key;column:serial_number;type:varchar;size:128;" json:"serial_number" db:"serial_number"`
	AuthorityKeyIdentifier string    `gorm:"primary_key;column:authority_key_identifier;type:varchar;size:128;" json:"authority_key_identifier" db:"authority_key_identifier"`
	Body                   string    `gorm:"column:body;type:text;size:65535;" json:"body" db:"body"`
	Expiry                 time.Time `gorm:"column:expiry;type:timestamp;" json:"expiry" db:"expiry"`
}

// TableName sets the insert table name for this struct type
func (o *OcspResponses) TableName() string {
	return "ocsp_responses"
}
