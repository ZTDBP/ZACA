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

package workload

import (
	"github.com/ztdbp/ZACA/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ztdbp/ZACA/core"
)

type Logic struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewLogic() *Logic {
	return &Logic{
		db:     core.Is.Db,
		logger: logger.Named("logic").SugaredLogger,
	}
}
