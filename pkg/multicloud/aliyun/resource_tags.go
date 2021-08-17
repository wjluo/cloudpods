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

package aliyun

import (
	"fmt"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/cloudprovider"
)

func (self *SRegion) tagRequest(serviceType, action string, params map[string]string) (jsonutils.JSONObject, error) {
	switch serviceType {
	case ALIYUN_SERVICE_ECS:
		return self.ecsRequest(action, params)
	case ALIYUN_SERVICE_VPC:
		return self.vpcRequest(action, params)
	case ALIYUN_SERVICE_RDS:
		return self.rdsRequest(action, params)
	case ALIYUN_SERVICE_SLB:
		return self.lbRequest(action, params)
	case ALIYUN_SERVICE_KVS:
		return self.kvsRequest(action, params)
	case ALIYUN_SERVICE_NAS:
		return self.nasRequest(action, params)
	default:
		return nil, fmt.Errorf("invalid service type")
	}
}

// 资源类型。取值范围：
// disk, instance, image, securitygroup, snapshot
func (self *SRegion) ListTags(serviceType string, resourceType string, resourceId string) ([]SAliyunTag, error) {
	tags := []SAliyunTag{}
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["ResourceType"] = resourceType
	params["ResourceId.1"] = resourceId
	nextToken := ""
	for {
		if len(nextToken) > 0 {
			params["NextToken"] = nextToken
		}
		resp, err := self.tagRequest(serviceType, "ListTagResources", params)
		if err != nil {
			return nil, errors.Wrapf(err, "%s ListTagResources %s", serviceType, params)
		}
		part := []SAliyunTag{}
		err = resp.Unmarshal(&part, "TagResources", "TagResource")
		if err != nil {
			return nil, errors.Wrapf(err, "resp.Unmarshal")
		}
		tags = append(tags, part...)
		nextToken, _ = resp.GetString("NextToken")
		if len(nextToken) == 0 {
			break
		}
	}
	return tags, nil
}

func (self *SRegion) UntagResource(serviceType string, resourceType string, resId string, keys []string) error {
	if len(resId) == 0 || len(keys) == 0 {
		return nil
	}

	params := map[string]string{
		"RegionId":     self.RegionId,
		"ResourceId.1": resId,
		"ResourceType": resourceType,
	}
	for i, key := range keys {
		params[fmt.Sprintf("TagKey.%d", i+1)] = key
	}

	_, err := self.tagRequest(serviceType, "UntagResources", params)
	return errors.Wrapf(err, "UntagResources %s", params)
}

func (self *SRegion) SetResourceTags(serviceType string, resourceType string, resId string, tags map[string]string, replace bool) error {
	err := self.TagResource(serviceType, resourceType, resId, tags)
	if err != nil {
		return errors.Wrapf(err, "TagResource")
	}
	if !replace || len(tags) == 0 {
		return nil
	}
	_, _tags, err := self.ListSysAndUserTags(serviceType, resourceType, resId)
	if err != nil {
		return errors.Wrapf(err, "ListTags")
	}
	tagMaps := map[string]string{}
	for k, v := range tags {
		tagMaps[strings.ToLower(k)] = v
	}
	keys := []string{}
	for k := range _tags {
		if _, ok := tagMaps[strings.ToLower(k)]; !ok {
			keys = append(keys, k)
		}
	}
	return self.UntagResource(serviceType, resourceType, resId, keys)
}

func (self *SRegion) TagResource(serviceType string, resourceType string, resourceId string, tags map[string]string) error {
	if len(tags) > 20 {
		return errors.Wrap(cloudprovider.ErrNotSupported, "tags count exceed 20 for one request")
	}
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["ResourceType"] = resourceType
	params["ResourceId.1"] = resourceId
	i := 0
	for k, v := range tags {
		if strings.HasPrefix(k, "aliyun") ||
			strings.HasPrefix(k, "acs:") ||
			strings.HasPrefix(k, "http://") ||
			strings.HasPrefix(k, "https://") ||
			strings.HasPrefix(v, "http://") ||
			strings.HasPrefix(v, "https://") ||
			strings.HasPrefix(v, "acs:") {
			continue
		}
		params[fmt.Sprintf("Tag.%d.Key", i+1)] = k
		params[fmt.Sprintf("Tag.%d.Value", i+1)] = v
		i++
	}
	action := "TagResources"
	if len(tags) == 0 {
		action = "UntagResources"
		params["All"] = "true"
	}
	_, err := self.tagRequest(serviceType, action, params)
	if err != nil {
		return errors.Wrapf(err, "%s %s %s", action, resourceId, params)
	}
	return nil
}

func (self *SRegion) ListSysAndUserTags(serviceType string, resourceType string, resourceId string) (map[string]string, map[string]string, error) {
	tags, err := self.ListTags(serviceType, resourceType, resourceId)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "ListTags(%s, %s)", resourceType, resourceId)
	}
	sys, user := map[string]string{}, map[string]string{}
	for _, tag := range tags {
		if strings.HasPrefix(tag.TagKey, "aliyun") || strings.HasPrefix(tag.TagKey, "acs:") {
			sys[tag.TagKey] = tag.TagValue
			continue
		}
		user[tag.TagKey] = tag.TagValue
	}
	return sys, user, nil
}
