package vps

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/internal/testutil"
)

func TestBigStorageRepository_GetBigStorages(t *testing.T) {
	const apiResponse = `{ "bigStorages": [ { "name": "example-bigstorage", "description": "Big storage description", "diskSize": 2147483648, "offsiteBackups": true, "vpsName": "example-vps", "status": "active", "isLocked": false, "availabilityZone": "ams0", "serial": "e7e12b3c7c6602973ac7" } ] } `
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	all, err := repo.GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(all))

	assert.Equal(t, "example-bigstorage", all[0].Name)
	assert.Equal(t, "Big storage description", all[0].Description)
	assert.EqualValues(t, 2147483648, all[0].DiskSize)
	assert.Equal(t, true, all[0].OffsiteBackups)
	assert.Equal(t, "example-vps", all[0].VpsName)
	assert.EqualValues(t, "active", all[0].Status)
	assert.Equal(t, false, all[0].IsLocked)
	assert.Equal(t, "ams0", all[0].AvailabilityZone)
	assert.Equal(t, "e7e12b3c7c6602973ac7", all[0].Serial)
}

func TestBigStorageRepository_GetSelection(t *testing.T) {
	const apiResponse = `{ "bigStorages": [ { "name": "example-bigstorage", "description": "Big storage description", "diskSize": 2147483648, "offsiteBackups": true, "vpsName": "example-vps", "status": "active", "isLocked": false, "availabilityZone": "ams0", "serial": "e7e12b3c7c6602973ac7"} ] } `
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages?page=1&pageSize=25", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	all, err := repo.GetSelection(1, 25)
	require.NoError(t, err)
	require.Equal(t, 1, len(all))

	assert.Equal(t, "example-bigstorage", all[0].Name)
	assert.Equal(t, "Big storage description", all[0].Description)
	assert.EqualValues(t, 2147483648, all[0].DiskSize)
	assert.Equal(t, true, all[0].OffsiteBackups)
	assert.Equal(t, "example-vps", all[0].VpsName)
	assert.EqualValues(t, "active", all[0].Status)
	assert.Equal(t, false, all[0].IsLocked)
	assert.Equal(t, "ams0", all[0].AvailabilityZone)
	assert.Equal(t, "e7e12b3c7c6602973ac7", all[0].Serial)
}

func TestBigStorageRepository_GetBigStorageByName(t *testing.T) {
	const apiResponse = `{ "bigStorage": { "name": "example-bigstorage", "description": "Big storage description", "diskSize": 2147483648, "offsiteBackups": true, "vpsName": "example-vps", "status": "active", "isLocked": false, "availabilityZone": "ams0", "serial": "e7e12b3c7c6602973ac7" } } `
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	bigstorage, err := repo.GetByName("example-bigstorage")
	require.NoError(t, err)
	assert.Equal(t, "example-bigstorage", bigstorage.Name)
	assert.Equal(t, "Big storage description", bigstorage.Description)
	assert.EqualValues(t, 2147483648, bigstorage.DiskSize)
	assert.Equal(t, true, bigstorage.OffsiteBackups)
	assert.Equal(t, "example-vps", bigstorage.VpsName)
	assert.EqualValues(t, "active", bigstorage.Status)
	assert.Equal(t, false, bigstorage.IsLocked)
	assert.Equal(t, "ams0", bigstorage.AvailabilityZone)
	assert.Equal(t, "e7e12b3c7c6602973ac7", bigstorage.Serial)
}

func TestBigStorageRepository_OrderBigStorage(t *testing.T) {
	const expectedRequest = `{"size":8,"offsiteBackups":true,"availabilityZone":"ams0","vpsName":"example-vps","description":"test-description"}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	order := BigStorageOrder{Size: 8, OffsiteBackups: true, AvailabilityZone: "ams0", VpsName: "example-vps", Description: "test-description"}
	err := repo.Order(order)

	require.NoError(t, err)
}

func TestBigStorageRepository_UpgradeBigStorage(t *testing.T) {
	const expectedRequest = `{"bigStorageName":"example-bigstorage","size":8,"offsiteBackups":true}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	err := repo.Upgrade("example-bigstorage", 8, true)

	require.NoError(t, err)
}

func TestBigStorageRepository_UpdateBigStorage(t *testing.T) {
	const expectedRequest = `{"bigStorage":{"name":"example-bigstorage","description":"Big storage description","diskSize":2147483648,"offsiteBackups":true,"vpsName":"example-vps","status":"active","serial":"e7e12b3c7c6602973ac7","isLocked":false,"availabilityZone":"ams0"}}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	bigStorage := BigStorage{
		Name:             "example-bigstorage",
		Description:      "Big storage description",
		DiskSize:         2147483648,
		OffsiteBackups:   true,
		VpsName:          "example-vps",
		Status:           BigStorageStatusActive,
		IsLocked:         false,
		AvailabilityZone: "ams0",
		Serial:           "e7e12b3c7c6602973ac7",
	}
	err := repo.Update(bigStorage)

	require.NoError(t, err)
}

func TestBigStorageRepository_DetachVpsFromBigStorage(t *testing.T) {
	const expectedRequest = `{"bigStorage":{"name":"example-bigstorage","description":"Big storage description","diskSize":2147483648,"offsiteBackups":true,"vpsName":"example-vps","status":"active","serial":"e7e12b3c7c6602973ac7","isLocked":false,"availabilityZone":"ams0"}}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	bigStorage := BigStorage{
		Name:             "example-bigstorage",
		Description:      "Big storage description",
		DiskSize:         2147483648,
		OffsiteBackups:   true,
		Status:           "active",
		IsLocked:         false,
		AvailabilityZone: "ams0",
		Serial:           "e7e12b3c7c6602973ac7",
	}
	err := repo.AttachToVps("example-vps", bigStorage)
	require.NoError(t, err)
}

func TestBigStorageRepository_AttachVpsToBigStorage(t *testing.T) {
	const expectedRequest = `{"bigStorage":{"name":"example-bigstorage","description":"Big storage description","diskSize":2147483648,"offsiteBackups":true,"vpsName":"","status":"active","serial":"e7e12b3c7c6602973ac7","isLocked":false,"availabilityZone":"ams0"}}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	bigStorage := BigStorage{
		Name:             "example-bigstorage",
		Description:      "Big storage description",
		DiskSize:         2147483648,
		OffsiteBackups:   true,
		VpsName:          "example-vps",
		Status:           "active",
		IsLocked:         false,
		AvailabilityZone: "ams0",
		Serial:           "e7e12b3c7c6602973ac7",
	}
	err := repo.DetachFromVps(bigStorage)
	require.NoError(t, err)
}

func TestBigStorageRepository_CancelBigStorage(t *testing.T) {
	const expectedRequest = `{"endTime":"end"}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage", ExpectedMethod: "DELETE", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	err := repo.Cancel("example-bigstorage", gotransip.CancellationTimeEnd)
	require.NoError(t, err)
}

func TestBigStorageRepository_GetBigStorageBackups(t *testing.T) {
	const apiResponse = `{ "backups": [ { "id": 1583, "status": "active", "diskSize": 4294967296, "dateTimeCreate": "2019-12-31 09:13:55", "availabilityZone": "ams0" } ] }`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage/backups", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	all, err := repo.GetBackups("example-bigstorage")
	require.NoError(t, err)
	require.Equal(t, 1, len(all))

	assert.EqualValues(t, 1583, all[0].ID)
	assert.EqualValues(t, "active", all[0].Status)
	assert.EqualValues(t, 4294967296, all[0].DiskSize)
	assert.Equal(t, "2019-12-31 09:13:55", all[0].DateTimeCreate.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "ams0", all[0].AvailabilityZone)
}

func TestBigStorageRepository_RevertBigStorageBackup(t *testing.T) {
	const expectedRequest = `{"action":"revert"}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage/backups/123", ExpectedMethod: "PATCH", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	err := repo.RevertBackup("example-bigstorage", 123)
	require.NoError(t, err)
}

func TestBigStorageRepository_RevertBackupToOtherBigStorage(t *testing.T) {
	const expectedRequest = `{"action":"revert","destinationBigStorageName":"example-bigStorage2"}`
	server := testutil.MockServer{T: t, ExpectedURL: "/big-storages/example-bigstorage/backups/123", ExpectedMethod: "PATCH", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	err := repo.RevertBackupToOtherBigStorage("example-bigstorage", 123, "example-bigStorage2")
	require.NoError(t, err)
}

func TestBigStorageRepository_GetBigStorageUsage(t *testing.T) {
	const apiResponse = `{ "usage": [ { "iopsRead": 0.27, "iopsWrite": 0.13, "date": 1574783109 } ] }`

	parameters := url.Values{
		"dateTimeStart": []string{"1500538995"},
		"dateTimeEnd":   []string{"1500542619"},
	}

	expectedURL := "/big-storages/example-bigstorage/usage?" + parameters.Encode()
	server := testutil.MockServer{T: t, ExpectedURL: expectedURL, ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	usageData, err := repo.GetUsage("example-bigstorage", UsagePeriod{TimeStart: 1500538995, TimeEnd: 1500542619})
	require.NoError(t, err)

	require.Equal(t, 1, len(usageData))
	assert.EqualValues(t, 0.27, usageData[0].IopsRead)
	assert.EqualValues(t, 0.13, usageData[0].IopsWrite)
	assert.EqualValues(t, 1574783109, usageData[0].Date)
}

func TestBigStorageRepository_GetBigStorageUsageLast24Hours(t *testing.T) {
	const apiResponse = `{ "usage": [ { "iopsRead": 0.27, "iopsWrite": 0.13, "date": 1574783109 } ] }`

	parameters := url.Values{
		"dateTimeStart": []string{fmt.Sprintf("%d", time.Now().Add(-24*time.Hour).Unix())},
		"dateTimeEnd":   []string{fmt.Sprintf("%d", time.Now().Unix())},
	}

	expectedURL := "/big-storages/example-bigstorage/usage?" + parameters.Encode()
	server := testutil.MockServer{T: t, ExpectedURL: expectedURL, ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := BigStorageRepository{Client: *client}

	usageData, err := repo.GetUsageLast24Hours("example-bigstorage")
	require.NoError(t, err)

	require.Equal(t, 1, len(usageData))
	assert.EqualValues(t, 0.27, usageData[0].IopsRead)
	assert.EqualValues(t, 0.13, usageData[0].IopsWrite)
	assert.EqualValues(t, 1574783109, usageData[0].Date)
}
