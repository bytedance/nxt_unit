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
package contexthelper

import (
	"context"
)

type builderVectorKey struct {
}

var BuilderVectorKey = builderVectorKey{}

func SetBuilderVector(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, BuilderVectorKey, id)
}

func GetBuilderVector(ctx context.Context) (string, bool) {
	value := ctx.Value(BuilderVectorKey)
	id, ok := value.(string)
	if !ok {
		return "", false
	}
	return id, true
}
