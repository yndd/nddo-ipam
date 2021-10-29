/*
Copyright 2021 NDDO.

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

package server

var hPathElements = map[string]interface{}{
	"/ipam": map[string]interface{}{
		"tenant": nil,
	},
	"/ipam/tenant": map[string]interface{}{
		"network-instance": nil,
	},
	"/ipam/tenant/network-instance": map[string]interface{}{
		"ip-prefix":  nil,
		"ip-address": nil,
		"ip-raange":  nil,
	},
	"/ipam/instance/network-instance/ip-prefix":  nil,
	"/ipam/instance/network-instance/ip-range":   nil,
	"/ipam/instance/network-instance/ip-address": nil,
}
