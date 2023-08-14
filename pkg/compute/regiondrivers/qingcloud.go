// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package regiondrivers

import (
	"context"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type SQingCloudRegionDriver struct {
	SManagedVirtualizationRegionDriver
}

func init() {
	driver := SQingCloudRegionDriver{}
	models.RegisterRegionDriver(&driver)
}

func (self *SQingCloudRegionDriver) GetProvider() string {
	return api.CLOUD_PROVIDER_QINGCLOUD
}

func (self *SQingCloudRegionDriver) IsAllowSecurityGroupNameRepeat() bool {
	return true
}

func (self *SQingCloudRegionDriver) GenerateSecurityGroupName(name string) string {
	return name
}

func (self *SQingCloudRegionDriver) IsSecurityGroupBelongVpc() bool {
	return true
}

func (self *SQingCloudRegionDriver) ValidateCreateSnapshotData(ctx context.Context, userCred mcclient.TokenCredential, disk *models.SDisk, storage *models.SStorage, input *api.SnapshotCreateInput) error {
	return nil
}
