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

package ipamlogic

import "github.com/openconfig/gnmi/proto/gnmi"

var refPaths2 = []*gnmi.Path{
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "rir", Key: map[string]string{
				"name": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "rir", Key: map[string]string{
				"name": "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
			{Name: "network-instance", Key: map[string]string{
				"name": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
			{Name: "network-instance", Key: map[string]string{
				"name": "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
			{Name: "network-instance", Key: map[string]string{
				"name": "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"prefix": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ipam"},
			{Name: "tenant", Key: map[string]string{
				"name": "",
			}},
			{Name: "network-instance", Key: map[string]string{
				"name": "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"prefix": "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
}
