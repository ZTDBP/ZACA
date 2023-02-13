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

// SelfKeypair struct is a row record of the self_keypair table in the cap database
type SelfKeypair struct {
	ID          uint32         `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id" db:"id"`
	Name        string         `gorm:"column:name;type:varchar;size:40;" json:"name" db:"name"`
	PrivateKey  sql.NullString `gorm:"column:private_key;type:text;size:65535;" json:"private_key" db:"private_key"`
	Certificate sql.NullString `gorm:"column:certificate;type:text;size:65535;" json:"certificate" db:"certificate"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;" json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;" json:"updated_at" db:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (s *SelfKeypair) TableName() string {
	return "self_keypair"
}
