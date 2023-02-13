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
	"database/sql"
	"time"
)

// Certificates struct is a row record of the certificates table in the cap database
type Certificates struct {
	SerialNumber           string         `gorm:"primary_key;column:serial_number;type:varchar;size:128;" json:"serial_number" db:"serial_number"`
	AuthorityKeyIdentifier string         `gorm:"primary_key;column:authority_key_identifier;type:varchar;size:128;" json:"authority_key_identifier" db:"authority_key_identifier"`
	CaLabel                sql.NullString `gorm:"column:ca_label;type:varchar;size:128;" json:"ca_label" db:"ca_label"`
	Status                 string         `gorm:"column:status;type:varchar;size:128;" json:"status" db:"status"`
	Reason                 sql.NullInt64  `gorm:"column:reason;type:int;" json:"reason" db:"reason"`
	Expiry                 time.Time      `gorm:"column:expiry;type:timestamp;" json:"expiry" db:"expiry"`
	RevokedAt              time.Time      `gorm:"column:revoked_at;type:timestamp;" json:"revoked_at" db:"revoked_at"`
	Pem                    string         `gorm:"column:pem;type:text;size:65535;" json:"pem" db:"pem"`
	IssuedAt               time.Time      `gorm:"column:issued_at;type:timestamp;" json:"issued_at" db:"issued_at"`
	NotBefore              time.Time      `gorm:"column:not_before;type:timestamp;" json:"not_before" db:"not_before"`
	Metadata               sql.NullString `gorm:"column:metadata;type:json;" json:"metadata" db:"metadata"`
	Sans                   sql.NullString `gorm:"column:sans;type:json;" json:"sans" db:"sans"`
	CommonName             sql.NullString `gorm:"column:common_name;type:text;size:65535;" json:"common_name" db:"common_name"`
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "certificates"
}
