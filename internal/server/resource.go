package server

import (
	"sort"
	"strings"
)

type resource struct {
	key  map[string]string
	data interface{}
}

func newResource(key map[string]string) *resource {
	return &resource{
		key:  key,
		data: nil,
	}
}

func (r *resource) GetKey() map[string]string {
	return r.key
}

func (r *resource) GetKeyString() string {
	sb := strings.Builder{}
	i := 0
	type kv struct {
		Key   string
		Value string
	}
	var ss []kv
	for k, v := range r.GetKey() {
		ss = append(ss, kv{k, v})
	}
	// sort the slice keys so we have a determinsitic result
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key > ss[j].Key
	})
	for _, kv := range ss {
		sb.WriteString(kv.Key)
		sb.WriteString("=")
		sb.WriteString(kv.Value)
		if i != len(ss)-1 {
			sb.WriteString(",")
		}
		i++
	}
	return sb.String()
}

func (r *resource) GetData() interface{} {
	return r.data
}
