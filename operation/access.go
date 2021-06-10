package operation

import (
	ctx "github.com/blentz/cloud-haunter/context"
	"github.com/blentz/cloud-haunter/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	ctx.Operations[types.CloudAccess] = access{}
}

type access struct {
}

func (o access) Execute(clouds []types.CloudType) []types.CloudItem {
	log.Debugf("Collecting old accesses on: [%s]", clouds)
	accessChan, errChan := o.collect(clouds)
	return wait(accessChan, errChan, "[ACCESS] Failed to collect accesses")
}

func (o access) collect(clouds []types.CloudType) (chan []types.CloudItem, chan error) {
	return collect(clouds, func(provider types.CloudProvider) ([]types.CloudItem, error) {
		accesses, err := provider.GetAccesses()
		if err != nil {
			return nil, err
		}
		return o.convertToCloudItems(accesses), nil
	})
}

func (o access) convertToCloudItems(accesses []*types.Access) []types.CloudItem {
	var items []types.CloudItem
	for _, access := range accesses {
		items = append(items, access)
	}
	return items
}
