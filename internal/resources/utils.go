// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rest "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
)

const (
	// dsFilter it's the name of the filter key
	dsFilter = "filter"
	// maxETagRetries sets the number of retries for a databse operation when an ETag mismatch error occurs during an update.
	maxETagRetries = 1000
	// minBackoffTime is the minimum time before a retry should be attempted
	minBackoffTime = 5 * time.Millisecond
	// backoffJitterTime is a maximum additional time that backoff can wait for
	backoffJitterTime = 95
)

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

func safeMapStrInt32(s interface{}) map[string]int32 {
	m, _ := s.(map[string]interface{})

	r := make(map[string]int32, len(m))

	for k, v := range m {
		i, _ := v.(int)
		r[k] = int32(i)
	}

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

// wrapResourceError ensures that any non-nil error is wrapped.
func wrapResourceError(err *error, msg string) {
	if err == nil || *err == nil {
		return
	}

	nErr := rest.GenericOpenAPIError{}

	if errors.As(*err, &nErr) {
		*err = fmt.Errorf("%s %s: %w", msg, strings.Trim(nErr.Message(), "\n "), *err)

		return
	}

	*err = fmt.Errorf("%s %w", msg, *err)
}

// convertMap returns map of string key to string value.
func convertMap(in map[string]interface{}) map[string]string {
	ret := make(map[string]string, len(in))
	if len(in) == 0 {
		return ret
	}

	for k, v := range in {
		if s, ok := v.(string); ok {
			ret[k] = s
		}
	}

	return ret
}

// expandStringList takes []interfaces and returns []strings.
func expandStringList(list []interface{}) []string {
	vs := make([]string, 0, len(list))

	for _, v := range list {
		val, ok := v.(string)
		if ok {
			vs = append(vs, val)
		}
	}

	return vs
}

// flattenStringList takes []strings and returns []interfaces.
func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}

	return vs
}

// retryBackoff holds a thread off for a random amount of time.
func retryBackoff() time.Duration {
	wait := minBackoffTime + time.Duration(rand.Intn(backoffJitterTime))*time.Millisecond
	time.Sleep(wait)

	return wait
}
