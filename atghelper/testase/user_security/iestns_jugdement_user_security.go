/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package user_security

import (
	"context"
)

type Hello struct {
	Url string
}

type AnchorData struct {
	Id int64 `json:"id"`
	// We use JSON 'omitempty' tag for State to reduce memory usage in Redis when a new item anchor relation is created..
	// This is because the value for State when setting or appending a new item anchor relation is by default 0.
	State int16 `json:"state,omitempty"`
}

type Args struct {
	Ctx *struct {
		context.Context
		InterfaceTag bool
	}
	AnchorList []AnchorData
	ItemId     int64
}

type MGetUserSecurityForSearchReq struct {
	BizType  string           `thrift:"bizType,1,required" json:"bizType"`
	Scene    int              `thrift:"scene,2,required" json:"scene"`
	RiskInfo map[int64]*Hello `thrift:"riskInfo,3,required" json:"riskInfo"`
}

func MGetUserSecurityForSearch(ctx context.Context, req *MGetUserSecurityForSearchReq) (*MGetUserSecurityForSearchReq, error) {
	return nil, nil
}
