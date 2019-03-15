package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	ctx "github.com/hortonworks/cloud-haunter/context"
	"github.com/hortonworks/cloud-haunter/types"
	"gopkg.in/yaml.v2"
)

// IsAnyMatch looks for any of the given tag in types.Tag
func IsAnyMatch(haystack map[string]string, needles ...string) bool {
	for k, v := range haystack {
		haystack[strings.ToLower(k)] = v
	}
	for _, n := range needles {
		if _, ok := haystack[n]; ok {
			return true
		}
	}
	return false
}

// IsAnyStartsWith looks any tag start with given needles
func IsAnyStartsWith(haystack map[string]string, needles ...string) bool {
	for k := range haystack {
		if IsStartsWith(k, needles...) {
			return true
		}
	}
	return false
}

// IsStartsWith looks input start with given needles
func IsStartsWith(hay string, needles ...string) bool {
	for _, n := range needles {
		if strings.Index(hay, n) == 0 {
			return true
		}
	}
	return false
}

// ConvertTimeRFC3339 converts RFC3339 format string to time.Time
func ConvertTimeRFC3339(stringTime string) (time.Time, error) {
	return time.Parse(time.RFC3339, stringTime)
}

// ConvertTimeLayout converts a string in the format of the layout to time.Time
func ConvertTimeLayout(layout, timeString string) (time.Time, error) {
	return time.Parse(layout, timeString)
}

// ConvertTimeUnix parses a unix timestamp (seconds since epoch start) from string to time.Time
func ConvertTimeUnix(unixTimestamp string) time.Time {
	timestamp, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		log.Warnf("[util.ConvertTimeUnix] cannot convert time: %s, err: %s", unixTimestamp, err)
		timestamp = 0
	}
	return time.Unix(timestamp, 0)
}

// ConvertTags converts a map of tags to types.Tag
func ConvertTags(tagMap map[string]*string) types.Tags {
	tags := make(types.Tags, 0)
	for k, v := range tagMap {
		tags[strings.ToLower(k)] = *v
	}
	return tags
}

// LoadFilterConfig loads and unmarshalls filter config YAML
func LoadFilterConfig(location string) (*types.FilterConfig, error) {
	raw, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	config := &types.FilterConfig{}
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}
	log.Debugf("[UTIL] Filter config loaded:\n%s", raw)
	return config, nil
}

// GetCloudAccountNames returns the name of the configured cloud accounts
func GetCloudAccountNames() map[types.CloudType]string {
	var accounts = make(map[types.CloudType]string)
	for cType, initFunc := range ctx.CloudProviders {
		accounts[cType] = initFunc().GetAccountName()
	}
	return accounts
}

// SplitListToMap splits comma separated list to key:true map
func SplitListToMap(list string) (resp map[string]bool) {
	resp = map[string]bool{}
	for _, i := range strings.Split(list, ",") {
		if len(i) != 0 {
			trimmed := strings.Trim(i, " ")
			resp[strings.ToLower(trimmed)] = true
			resp[strings.ToUpper(trimmed)] = true
		}
	}
	return
}

// CovertJsonToString converts a struct to json string
func CovertJsonToString(source interface{}) (*string, error) {
	j, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	return &(&types.S{S: string(j)}).S, nil
}

// GetFilterNames returns the name of the applied filters separated by colon
func GetFilterNames(filters []types.FilterType) string {
	if len(filters) == 0 {
		return "noFilter"
	}
	fNames := make([]string, 0)
	for _, f := range filters {
		fNames = append(fNames, f.String())
	}
	return strings.Join(fNames, ",")
}
