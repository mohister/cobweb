/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package router

type Node struct {
	Pattern string
	Nodes []*Node
}

//func (n *Node)AddNode(pattern string, handler http.Handler)  {
//	target := n
//	var nn *node
//	// todo: 增加层级，更好的命中handler
//	path = internal.Trim(path, '/')
//	var err error
//	for next, part := 0, ""; err == nil; {
//		next, part, err = internal.Traverse(path, next, '/')
//		if err != nil && err != io.EOF {
//			panic(err)
//		}
//		switch part[0] {
//		case ':':
//			if len(part) == 1 || !checkPart(part[1:]) {
//				panic("Wide pattern must be param")
//			}
//			if err == io.EOF && target.paramHandle {
//				panic("Ending point method mush be only")
//			}
//			for i := 0; i < len(target.params); i++ {
//				if target.params[i].pattern == part[1:] {
//					nn = target.params[i]
//					break
//				}
//			}
//			if nn == nil {
//				nn = newNode()
//				target.params = append(target.params, nn)
//			}
//			nn.pattern = part[1:]
//			if err == io.EOF {
//				if nn.method != nil {
//					panic("Ending point method mush be only")
//				}
//				nn.method = handler
//				target.paramHandle = true
//			}
//		case '?':
//			if err != io.EOF || len(part) < 3 || part[1] != ':' &&
//				!checkPart(part[2:]) || target.wide != nil {
//				panic("Wide pattern must be only in ending point")
//			}
//			nn = newNode()
//			nn.pattern = part[2:]
//			nn.method = handler
//			target.wide = nn
//		default:
//			if !checkPart(part) {
//				panic("Wide pattern must be static")
//			}
//			for i := 0; i < len(target.statics); i++ {
//				if target.statics[i].pattern == part {
//					nn = target.statics[i]
//					break
//				}
//			}
//			if nn == nil {
//				nn = newNode()
//				target.statics = append(target.statics, nn)
//			}
//			nn.pattern = part
//			if err == io.EOF {
//				if nn.method != nil {
//					panic("Ending point method mush be only")
//				}
//				nn.method = handler
//			}
//		}
//		target = nn
//		nn = nil
//	}
//}
