/*
Copyright 2022-present The Ztalab Authors.
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

package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ztdbp/ZACA/database/mysql/cfssl-model/model"
)

// GetSelfKeypair
func GetSelfKeypair(db *gorm.DB, argId uint32) (record *model.SelfKeypair, err error) {
	record = &model.SelfKeypair{}
	if err = db.Where("id = ?", argId).First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return record, nil
}
