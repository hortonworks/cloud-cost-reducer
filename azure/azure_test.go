package azure

import (
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest"
	ctx "github.com/hortonworks/cloud-haunter/context"
	"github.com/hortonworks/cloud-haunter/types"

	"github.com/stretchr/testify/assert"
)

type callInfo struct {
	invocations []interface{}
}

func TestProviderInit(t *testing.T) {
	provider := azureProvider{}

	authorizer := func() (autorest.Authorizer, error) {
		return autorest.NullAuthorizer{}, nil
	}

	provider.init("AZURE_SUBSCRIPTION_ID", authorizer)

	assert.Equal(t, "AZURE_SUBSCRIPTION_ID", provider.subscriptionID)
	assert.NotNil(t, provider.vmClient)
	assert.NotNil(t, provider.vmClient.Authorizer)
}

func Test_givenTimestampIsInTags_whenGetCreationTimeFromTags_thenReturnsConvertedTimestamp(t *testing.T) {
	testValues := struct {
		timeAsUnixTimeStamp string
		timeAsTime          time.Time
	}{
		timeAsUnixTimeStamp: "1527240203",
		timeAsTime:          time.Date(2018, 5, 25, 11, 23, 23, 0, time.Local),
	}
	tags := types.Tags{ctx.AzureCreationTimeLabel: testValues.timeAsUnixTimeStamp}
	callInfo, stubConverterFunc := getStubConvertTimeUnixByTime(testValues.timeAsTime)

	getCreationTimeFromTags(tags, stubConverterFunc)

	assert.Equal(t, len(callInfo.invocations), 1)
	assert.Equal(t, callInfo.invocations[0].(string), testValues.timeAsUnixTimeStamp)
}

func Test_givenTimestampNotInTags_whenGetCreationTimeFromTags_thenReturnsEpochZeroTime(t *testing.T) {
	callInfo, stubConverterFunc := getStubConvertTimeUnixEpochZero()

	getCreationTimeFromTags(types.Tags{}, stubConverterFunc)

	assert.Equal(t, len(callInfo.invocations), 1)
	assert.Equal(t, callInfo.invocations[0].(string), "0")
}

func TestGetResourceGroupName(t *testing.T) {
	resourceGroupName := getResourceGroupName("/subscriptions/<sub_id>/resourceGroups/<rg_name>/providers/Microsoft.Compute/virtualMachines/<inst_name>")

	assert.Equal(t, "<rg_name>", resourceGroupName)
}

func TestGetResourceGroupNameNotFound(t *testing.T) {
	resourceGroupName := getResourceGroupName("")

	assert.Equal(t, "", resourceGroupName)
}

func getStubConvertTimeUnixByTime(timeAsTime time.Time) (*callInfo, func(string) time.Time) {
	cInfo := callInfo{invocations: make([]interface{}, 0, 3)}
	return &cInfo, func(unixTimestamp string) time.Time {
		cInfo.invocations = append(cInfo.invocations, unixTimestamp)
		return timeAsTime
	}
}

func getStubConvertTimeUnixEpochZero() (*callInfo, func(string) time.Time) {
	cInfo := callInfo{invocations: make([]interface{}, 0, 3)}
	return &cInfo, func(unixTimestamp string) time.Time {
		cInfo.invocations = append(cInfo.invocations, unixTimestamp)
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	}
}
