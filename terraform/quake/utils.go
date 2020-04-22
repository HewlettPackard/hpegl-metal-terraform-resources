package quake

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
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
	if fSet, ok := d.GetOk(dsFilter); ok {
		if flts, ok := fSet.(*schema.Set); ok {
			var values []*regexp.Regexp
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
		}
	}
	return
}

func convertStringArr(a []interface{}) (ret []string) {
	for _, v := range a {
		if v == nil {
			continue
		}
		ret = append(ret, v.(string))
	}
	return ret
}
