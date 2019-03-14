package action

import (
	"fmt"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	ctx "github.com/hortonworks/cloud-haunter/context"
	"github.com/hortonworks/cloud-haunter/types"
)

func init() {
	ctx.Actions[types.StopAction] = new(stopAction)
}

type stopAction struct {
}

func (s stopAction) Execute(op types.OpType, filters []types.FilterType, items []types.CloudItem) {
	instancesPerCloud := map[types.CloudType][]*types.Instance{}
	for _, item := range items {
		switch t := item.GetItem().(type) {
		case types.Instance:
			instancesPerCloud[item.GetCloudType()] = append(instancesPerCloud[item.GetCloudType()], item.(*types.Instance))
		default:
			log.Debugf("[STOP] Ignoring cloud item: %s, because it's not an instance, but a %s", t, item.GetType())
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(instancesPerCloud))
	for cloud, instances := range instancesPerCloud {
		go func(cloud types.CloudType, instances []*types.Instance) {
			defer wg.Done()

			log.Infof("[STOP] Stop %d instances on %s: %s", len(instances), cloud, strings.Join(getInstanceNames(instances), ","))
			if errors := ctx.CloudProviders[cloud]().StopInstances(instances); len(errors) != 0 {
				for _, err := range errors {
					log.Errorf("[STOP] Failed to stop instances on cloud: %s, err: %s", cloud, err.Error())
				}
			}
		}(cloud, instances)
	}

	wg.Wait()
}

func getInstanceNames(instances []*types.Instance) []string {
	result := make([]string, len(instances))
	for _, inst := range instances {
		result = append(result, fmt.Sprintf("%s:%s", inst.ID, inst.Name))
	}
	return result
}
