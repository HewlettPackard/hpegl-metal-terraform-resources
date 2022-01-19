// (C) Copyright 2016-2022 Hewlett Packard Enterprise Development LP.

package quake

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/quake-client/pkg/terraform/configuration"
)

const dsFilter = "filter"

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

type filter struct {
	name   string
	values []*regexp.Regexp
}

func (f filter) match(name string, value interface{}) bool {
	if f.name != name {
		return false
	}
	if s, ok := value.(string); ok {
		for _, v := range f.values {
			if v.MatchString(s) {
				return true
			}
		}
	}
	return false
}

func getFilters(d *schema.ResourceData) (filters []filter, err error) {
	fSet, ok := d.GetOk(dsFilter)
	if !ok {
		return
	}
	flts, ok := fSet.(*schema.Set)
	if !ok {
		return
	}
	values := []*regexp.Regexp{}
	for _, f := range flts.List() {
		m := f.(map[string]interface{})
		if name, ok := m["name"].(string); ok {
			for _, v := range m["values"].([]interface{}) {
				if value, ok := v.(string); ok {
					r, err := regexp.Compile(value)
					if err != nil {
						return nil, fmt.Errorf("%s for %q", err.Error(), name)
					}
					values = append(values, r)
				}
			}
			filters = append(filters, filter{
				name:   name,
				values: values,
			})
		}
	}
	return
}

func convertStringArr(a []interface{}) []string {
	ret := make([]string, len(a))

	for idx, v := range a {
		if v == nil {
			continue
		}

		ret[idx] = v.(string)
	}

	return ret
}

// Function to fish Config out of meta.  meta can be of two forms:
// *configuration.Config - if so return it
// map[string]interface{} - if so check if *configuration.Config is
//    present at key configuration.KeyForGLClientMap(), if not
//    return error
// If neither of the above, return an error
func getConfigFromMeta(meta interface{}) (*configuration.Config, error) {
	switch v := meta.(type) {

	case *configuration.Config:
		return v, nil

	case map[string]interface{}:
		// Check if Config can be found where expected in the map, if not return an error
		switch ct := v[configuration.KeyForGLClientMap()].(type) {
		case *configuration.Config:
			return ct, nil
		}
	}
	// Don't know what form meta is in, return an error

	errStr := "cannot find client, make sure that a .gltform file or a hpegl stanza metal block is present"

	return nil, fmt.Errorf(errStr)
}

func safeString(s interface{}) string {
	r, _ := s.(string)
	return r
}

func safeInt(s interface{}) int {
	r, _ := s.(int)
	return r
}

func safeFloat(s interface{}) float64 {
	r, _ := s.(float64)
	return r
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}

	diff := make([]string, 0)

	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}
