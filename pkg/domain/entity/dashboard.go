package entity

import (
	"net/url"
	"sort"
	"strings"
)

type DashboardOpts struct {
	DashboardID string
	Variables   map[string][]string
}

// EncodeVariables encodes the values into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (o DashboardOpts) EncodeVariables() string {
	if o.Variables == nil {
		return ""
	}

	var buf strings.Builder
	keys := make([]string, 0, len(o.Variables))
	for k := range o.Variables {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		vs := o.Variables[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}

			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}

	return buf.String()
}
