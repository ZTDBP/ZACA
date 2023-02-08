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

package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ztdbp/ZACA/database/mysql/cfssl-model/model"
)

// GetAllForbid
func GetAllForbid(db *gorm.DB, page, pagesize int, order string) (results []*model.Forbid, totalRows int64, err error) {

	resultOrm := db.Model(&model.Forbid{})
	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return results, totalRows, nil
		}
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetForbid
func GetForbid(db *gorm.DB) (record *model.Forbid, err error) {
	record = &model.Forbid{}
	if err = db.First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return record, err
	}

	return record, nil
}

// AddForbid
func AddForbid(db *gorm.DB, record *model.Forbid) (result *model.Forbid, RowsAffected int64, err error) {
	query := db.Save(record)
	if err = query.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, query.RowsAffected, nil
}
