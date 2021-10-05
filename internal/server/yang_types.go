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

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

var resourceRefPaths = []*gnmi.Path{
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "aggregate", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "aggregate", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-address", Key: map[string]string{
				"address":          "",
				"network-instance": "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-address", Key: map[string]string{
				"address":          "",
				"network-instance": "",
				"tenant":           "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-address", Key: map[string]string{
				"address":          "",
				"network-instance": "",
				"tenant":           "",
			}},
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-address", Key: map[string]string{
				"address":          "",
				"network-instance": "",
				"tenant":           "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
			{Name: "aggregate", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "aggregate", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "rir", Key: map[string]string{
				"name": "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "nddo-ipam"},
			{Name: "ipam"},
			{Name: "rir", Key: map[string]string{
				"name": "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
}
