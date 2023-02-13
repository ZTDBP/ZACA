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
	"github.com/guregu/null"
	"time"
)

// Forbid struct is a row record of the forbid table in the cap database
type Forbid struct {
	ID        uint32    `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id" db:"id"`
	UniqueID  string    `gorm:"column:unique_id;type:varchar;size:40;" json:"unique_id" db:"unique_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;" json:"updated_at" db:"updated_at"`
	DeletedAt null.Time `gorm:"column:deleted_at;type:timestamp;" json:"deleted_at" db:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (f *Forbid) TableName() string {
	return "forbid"
}
