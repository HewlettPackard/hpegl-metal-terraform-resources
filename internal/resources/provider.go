// (C) Copyright 2020-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	pollInterval = 3 * time.Second
)

var (
	resourceDefaultTimeouts *schema.ResourceTimeout
)

func init() {
	d := time.Minute * 60
	resourceDefaultTimeouts = &schema.ResourceTimeout{
		Create:  schema.DefaultTimeout(d),
		Update:  schema.DefaultTimeout(d),
		Delete:  schema.DefaultTimeout(d),
		Default: schema.DefaultTimeout(d),
	}
}
