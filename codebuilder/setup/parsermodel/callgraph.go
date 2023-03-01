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
package parsermodel

//函数定义
type FuncDesc struct {
	File    string //文件路径
	Path    string //文件路径
	Package string //package名
	Name    string //函数名，格式为Package.Func
}

//描述一个函数调用N个函数的一对多关系
type CallerRelation struct {
	Caller  FuncDesc
	Callees []FuncDesc
}

type CallGraph struct {
	RootFunc          *FuncDesc
	CallerRelationMap map[string]*CallerRelation
}

//描述关键函数的一条反向调用关系
type CalledRelation struct {
	Callees []FuncDesc
	CanFix  bool //该调用关系能反向找到gin.Context即可以自动修复
}

type MWTNode struct {
	Key      string
	Value    FuncDesc
	N        int
	Children []*MWTNode
}
