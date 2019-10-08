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

func CountParts(path string, splitter byte) (count uint8) {
	l := len(path)
	switch l {
	case 0:
	case 1:
		if path[0] != splitter {
			count = 1
		}
	default:
		ended := false
		for n := 0; n < l; n++ {
			matched := path[n] != splitter
			if matched && !ended {
				count++
			}
			ended = matched
		}
	}
	return
}

func NextPart(path string, from int, splitter byte) (part string, next int, ended bool) {
	l := len(path)
	switch l {
	case from:
		next, ended = from, true
	case from + 1:
		if path[from] == splitter {
			next, ended = from, true
		} else {
			part, next, ended = path, l, true
		}
	default:
		findFrom, findPart, findEnding := 0, 1, 2
		status := findFrom
		to := 0
		for next = from; next < l; next++ {
			if status == findFrom {
				if path[next] != splitter{
					from = next
					status = findPart
				}
			}else if status == findPart {
				if path[next] == splitter {
					to = next
					status = findEnding
				}
			}else{
				if path[next] != splitter{
					next = next -1
					break
				}
			}
		}
		if status != 0 && to == 0{
			to = next
		}
		if to - 1 >= from {
			part = path[from:to]
		}
		ended = bool(l == next)
	}
	return
}

func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}
