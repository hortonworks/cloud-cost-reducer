package operation

import (
	"errors"
	"testing"
	"time"

	ctx "github.com/hortonworks/cloud-haunter/context"
	"github.com/hortonworks/cloud-haunter/types"
	"github.com/stretchr/testify/assert"
)

func init() {
	ctx.CloudProviders[types.DUMMY] = func() types.CloudProvider {
		return dummyProvider{}
	}
}
func TestCollect(t *testing.T) {
	clouds := []types.CloudType{types.DUMMY}
	getter := func(p types.CloudProvider) (i []types.CloudItem, e error) {
		return []types.CloudItem{types.Instance{Name: "instance"}}, nil
	}

	itemsChan, _ := collect(clouds, getter)

	items := [][]types.CloudItem{
		<-itemsChan,
		<-itemsChan,
	}

	assert.Equal(t, "instance", items[0][0].GetName())
}

func TestWait(t *testing.T) {
	itemsChan := make(chan []types.CloudItem, 10)
	errChan := make(chan error, 5)

	go func() {
		items := wait(itemsChan, errChan, "error message")

		assert.Equal(t, "instance", items[0].GetName())
	}()

	itemsChan <- []types.CloudItem{types.Instance{Name: "instance"}}
	errChan <- errors.New("some error")

	time.Sleep(100 * time.Millisecond)

	close(itemsChan)
	close(errChan)
}

func TestConvertInstancesToCloudItems(t *testing.T) {
	instances := []*types.Instance{{}}

	items := convertInstancesToCloudItems(instances)

	assert.Equal(t, len(instances), len(items))
}

type dummyProvider struct {
}

func (p dummyProvider) GetAccountName() string {
	return "dummy"
}

func (p dummyProvider) GetInstances() (i []*types.Instance, e error) {
	return
}

func (p dummyProvider) TerminateInstances(i []*types.Instance) (e []error) {
	return
}

func (p dummyProvider) StopInstances(i []*types.Instance) (e []error) {
	return
}

func (p dummyProvider) GetAccesses() (a []*types.Access, e error) {
	return
}

func (p dummyProvider) GetDatabases() ([]*types.Database, error) {
	return nil, nil
}

func (p dummyProvider) GetDisks() ([]*types.Disk, error) {
	return nil, nil
}

func (p dummyProvider) DeleteDisks([]*types.Disk) []error {
	return nil
}

func (p dummyProvider) GetImages() ([]*types.Image, error) {
	return nil, nil
}
