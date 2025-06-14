package kubernetes

import (
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/assi010/gotransip/v6/internal/testutil"
	"github.com/assi010/gotransip/v6/rest"
	"github.com/assi010/gotransip/v6/vps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// To be compatible with < Go 1.20
const dateOnlyFormat = "2006-01-02"

func TestRepository_GetClusters(t *testing.T) {
	const apiResponse = `{"clusters":[{"name":"k888k","description":"production cluster","isLocked":true,"isBlocked": false},{"name":"aiceayoo","description":"development cluster","isLocked":false,"isBlocked":true}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetClusters()
	require.NoError(t, err)

	if assert.Equal(t, 2, len(list)) {
		assert.Equal(t, "k888k", list[0].Name)
		assert.Equal(t, "production cluster", list[0].Description)
		assert.True(t, list[0].IsLocked)
		assert.False(t, list[0].IsBlocked)
		assert.Equal(t, "aiceayoo", list[1].Name)
		assert.Equal(t, "development cluster", list[1].Description)
		assert.False(t, list[1].IsLocked)
		assert.True(t, list[1].IsBlocked)
	}
}

func TestRepository_GetClusterByName(t *testing.T) {
	const apiResponse = `{"cluster":{"name":"k888k","description":"production cluster","isLocked":true,"isBlocked": false, "version": "1.24.10", "endpoint": "https://kaas.transip.dev"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	cluster, err := repo.GetClusterByName("k888k")
	require.NoError(t, err)

	assert.Equal(t, "k888k", cluster.Name)
	assert.Equal(t, "production cluster", cluster.Description)
	assert.Equal(t, "https://kaas.transip.dev", cluster.Endpoint)
	assert.Equal(t, "1.24.10", cluster.Version)
	assert.True(t, cluster.IsLocked)
	assert.False(t, cluster.IsBlocked)
}

func TestRepository_CreateCluster(t *testing.T) {
	const expectedRequestBody = `{"description":"production cluster"}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequestBody}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	order := ClusterOrder{
		Description: "production cluster",
	}

	err := repo.CreateCluster(order)
	require.NoError(t, err)
}

func TestRepository_UpdateCluster(t *testing.T) {
	const expectedRequest = `{"cluster":{"name":"k888k","description":"staging cluster","version":"1.24.10","endpoint":"","isLocked":false,"isBlocked":false}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	clusterToUpdate := Cluster{
		Name:        "k888k",
		Description: "staging cluster",
		Version:     "1.24.10",
	}

	err := repo.UpdateCluster(clusterToUpdate)

	require.NoError(t, err)
}

func TestRepository_RemoveCluster(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k", ExpectedMethod: "DELETE", StatusCode: 204}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.RemoveCluster("k888k")
	require.NoError(t, err)
}

func TestRepository_GetKubeConfig(t *testing.T) {
	const apiResponse = `{"kubeConfig": {"encodedYaml": "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6IFtdCmNvbnRleHRzOiBbXQpraW5kOiBDb25maWcKcHJlZmVyZW5jZXM6IHt9CnVzZXJzOiBbXQoK"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/kubeconfig", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	config, err := repo.GetKubeConfig("k888k")
	require.NoError(t, err)

	assert.Contains(t, config, "apiVersion: v1")
}

func TestRepository_GetNodePools(t *testing.T) {
	const apiResponse = `{"nodePools":[{"uuid":"402c2f84-c37d-9388-634d-00002b7c6a82","description":"frontend","desiredNodeCount":3,"nodeSpec":"vps-bladevps-x4","availabilityZone":"ams0","labels":{"foo":"bar"},"taints":[{"key":"foo","value":"bar","effect":"NoSchedule"}],"nodes":[{"uuid":"76743b28-f779-3e68-6aa1-00007fbb911d","nodePoolUuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","status":"active"}]}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetNodePools("k888k")
	require.NoError(t, err)

	if assert.Equal(t, 1, len(list)) {
		assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", list[0].UUID)
		assert.Equal(t, "frontend", list[0].Description)
		assert.Equal(t, 3, list[0].DesiredNodeCount)
		assert.Equal(t, "vps-bladevps-x4", list[0].NodeSpec)
		assert.Equal(t, "ams0", list[0].AvailabilityZone)
		if assert.Equal(t, 1, len(list[0].Nodes)) {
			assert.Equal(t, NodeStatusActive, list[0].Nodes[0].Status)
			assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", list[0].Nodes[0].UUID)
			assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", list[0].Nodes[0].NodePoolUUID)
			assert.Equal(t, "k888k", list[0].Nodes[0].ClusterName)
		}
	}
}

func TestRepository_GetNodePool(t *testing.T) {
	const apiResponse = `{"nodePool":{"uuid":"402c2f84-c37d-9388-634d-00002b7c6a82","description":"frontend","desiredNodeCount":3,"nodeSpec":"vps-bladevps-x4","availabilityZone":"ams0","labels":{"foo":"bar"},"taints":[{"key":"foo","value":"bar","effect":"NoSchedule"}],"nodes":[{"uuid":"76743b28-f779-3e68-6aa1-00007fbb911d","nodePoolUuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","status":"active"}]}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	nodePool, err := repo.GetNodePool("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82")
	require.NoError(t, err)

	assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", nodePool.UUID)
	assert.Equal(t, "frontend", nodePool.Description)
	assert.Equal(t, 3, nodePool.DesiredNodeCount)
	assert.Equal(t, "vps-bladevps-x4", nodePool.NodeSpec)
	assert.Equal(t, "ams0", nodePool.AvailabilityZone)
	if assert.Equal(t, 1, len(nodePool.Nodes)) {
		assert.Equal(t, NodeStatusActive, nodePool.Nodes[0].Status)
		assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", nodePool.Nodes[0].UUID)
		assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", nodePool.Nodes[0].NodePoolUUID)
		assert.Equal(t, "k888k", nodePool.Nodes[0].ClusterName)
	}
}

func TestRepository_AddNodePool(t *testing.T) {
	const expectedRequestBody = `{"description":"frontend","desiredNodeCount":3,"nodeSpec":"vps-bladevps-x4","availabilityZone":"ams0"}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequestBody}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	order := NodePoolOrder{
		ClusterName:      "k888k",
		Description:      "frontend",
		DesiredNodeCount: 3,
		NodeSpec:         "vps-bladevps-x4",
		AvailabilityZone: "ams0",
	}

	err := repo.AddNodePool(order)
	require.NoError(t, err)
}

func TestRepository_UpdateNodePool(t *testing.T) {
	const expectedRequest = `{"nodePool":{"uuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","description":"backend","desiredNodeCount":4,"nodeSpec":"vps-bladevps-x8","availabilityZone":"ams0"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	nodePoolToUpdate := NodePool{
		UUID:             "402c2f84-c37d-9388-634d-00002b7c6a82",
		ClusterName:      "k888k",
		Description:      "backend",
		DesiredNodeCount: 4,
		NodeSpec:         "vps-bladevps-x8",
		AvailabilityZone: "ams0",
	}

	err := repo.UpdateNodePool(nodePoolToUpdate)

	require.NoError(t, err)
}

func TestRepository_RemoveNodePool(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82", ExpectedMethod: "DELETE", StatusCode: 204}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.RemoveNodePool("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82")
	require.NoError(t, err)
}

func TestRepository_GetNodes(t *testing.T) {
	const apiResponse = `{"nodes":[{"uuid":"76743b28-f779-3e68-6aa1-00007fbb911d","nodePoolUuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","status":"active","ipAddresses":[{"address":"37.97.254.6","subnetMask":"255.255.255.0","type":"external"}]}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/nodes", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetNodes("k888k")
	require.NoError(t, err)

	if assert.Equal(t, 1, len(list)) {
		assert.Equal(t, NodeStatusActive, list[0].Status)
		assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", list[0].UUID)
		assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", list[0].NodePoolUUID)
		assert.Equal(t, "k888k", list[0].ClusterName)
		if assert.Equal(t, 1, len(list[0].IPAddresses)) {
			assert.Equal(t, "37.97.254.6", list[0].IPAddresses[0].Address.String())
			assert.Equal(t, "255.255.255.0", list[0].IPAddresses[0].Netmask.String())
			assert.Equal(t, NodeAddressTypeExternal, list[0].IPAddresses[0].Type)
		}
	}
}

func TestRepository_GetNodesByNodePoolUUID(t *testing.T) {
	const apiResponse = `{"nodes":[{"uuid":"76743b28-f779-3e68-6aa1-00007fbb911d","nodePoolUuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","status":"active","ipAddresses":[{"address":"37.97.254.6","subnetMask":"255.255.255.0","type":"external"}]}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/nodes?nodePoolUuid=402c2f84-c37d-9388-634d-00002b7c6a82", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetNodesByNodePoolUUID("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82")
	require.NoError(t, err)

	if assert.Equal(t, 1, len(list)) {
		assert.Equal(t, NodeStatusActive, list[0].Status)
		assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", list[0].UUID)
		assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", list[0].NodePoolUUID)
		assert.Equal(t, "k888k", list[0].ClusterName)
		if assert.Equal(t, 1, len(list[0].IPAddresses)) {
			assert.Equal(t, "37.97.254.6", list[0].IPAddresses[0].Address.String())
			assert.Equal(t, "255.255.255.0", list[0].IPAddresses[0].Netmask.String())
			assert.Equal(t, NodeAddressTypeExternal, list[0].IPAddresses[0].Type)
		}
	}
}

func TestRepository_GetNode(t *testing.T) {
	const apiResponse = `{"node":{"uuid":"76743b28-f779-3e68-6aa1-00007fbb911d","nodePoolUuid":"402c2f84-c37d-9388-634d-00002b7c6a82","clusterName":"k888k","status":"active","ipAddresses":[{"address":"37.97.254.6","subnetMask":"255.255.255.0","type":"external"}]}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/nodes/76743b28-f779-3e68-6aa1-00007fbb911d", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	node, err := repo.GetNode("k888k", "76743b28-f779-3e68-6aa1-00007fbb911d")
	require.NoError(t, err)

	assert.Equal(t, NodeStatusActive, node.Status)
	assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", node.UUID)
	assert.Equal(t, "402c2f84-c37d-9388-634d-00002b7c6a82", node.NodePoolUUID)
	assert.Equal(t, "k888k", node.ClusterName)
	if assert.Equal(t, 1, len(node.IPAddresses)) {
		assert.Equal(t, "37.97.254.6", node.IPAddresses[0].Address.String())
		assert.Equal(t, "255.255.255.0", node.IPAddresses[0].Netmask.String())
		assert.Equal(t, NodeAddressTypeExternal, node.IPAddresses[0].Type)
	}
}

func TestRepository_RebootNode(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/nodes/76743b28-f779-3e68-6aa1-00007fbb911d", ExpectedMethod: "PATCH", ExpectedRequest: "{\"action\":\"reboot\"}"}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	t.Run("test valid reboot", func(t *testing.T) {
		server.StatusCode = 204
		err := repo.RebootNode("k888k", "76743b28-f779-3e68-6aa1-00007fbb911d")
		require.NoError(t, err)
	})

	t.Run("test nonexistent node", func(t *testing.T) {
		server.StatusCode = 404
		server.Response = `{"error": "Node with uuid '76743b28-f779-3e68-6aa1-00007fbb911d' not found"}`
		err := repo.RebootNode("k888k", "76743b28-f779-3e68-6aa1-00007fbb911d")

		if assert.Error(t, err) {
			assert.Equal(t, &rest.Error{
				Message:    "Node with uuid '76743b28-f779-3e68-6aa1-00007fbb911d' not found",
				StatusCode: 404,
			}, err)
		}
	})

	t.Run("test locked node", func(t *testing.T) {
		server.StatusCode = 409
		server.Response = `{"error": "Actions on Node '76743b28-f779-3e68-6aa1-00007fbb911d' are temporary disabled"}`
		err := repo.RebootNode("k888k", "76743b28-f779-3e68-6aa1-00007fbb911d")

		if assert.Error(t, err) {
			assert.Equal(t, &rest.Error{
				Message:    "Actions on Node '76743b28-f779-3e68-6aa1-00007fbb911d' are temporary disabled",
				StatusCode: 409,
			}, err)
		}
	})
}

func TestRepository_GetNodeStatistics(t *testing.T) {
	const apiResponse = `
	{
		"usage": {
		  "cpu": [
				{
				"percentage": 3.11,
				"date": 1500538995
				}
			],
		  "disk": [
				{
				"iopsRead": 0.27,
				"iopsWrite": 0.13,
				"date": 1500538995
				}
			],
		  "network": [
				{
				"mbitOut": 100.2,
				"mbitIn": 249.93,
				"date": 1500538995
				}
			]
		}
	}`

	values := url.Values{
		"dateTimeStart": []string{"1500538995"},
		"dateTimeEnd":   []string{"1500542619"},
	}

	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/nodes/76743b28-f779-3e68-6aa1-00007fbb911d/stats?" + values.Encode(),
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	statistics, err := repo.GetNodeStatistics(
		"k888k",
		"76743b28-f779-3e68-6aa1-00007fbb911d",
		[]vps.UsageType{},
		vps.UsagePeriod{
			TimeStart: 1500538995,
			TimeEnd:   1500542619,
		})
	require.NoError(t, err)

	cpuStatistics := statistics.CPU
	if assert.NotEmpty(t, cpuStatistics) {
		assert.Equal(t, vps.UsageDataCPU{Date: 1500538995, Percentage: 3.11}, cpuStatistics[0])
	}

	diskStatistics := statistics.Disk
	if assert.NotEmpty(t, diskStatistics) {
		assert.Equal(t, vps.UsageDataDisk{Date: 1500538995, IopsRead: 0.27, IopsWrite: 0.13}, diskStatistics[0])
	}

	networkStatistics := statistics.Network
	if assert.NotEmpty(t, networkStatistics) {
		assert.Equal(t, vps.UsageDataNetwork{Date: 1500538995, MbitOut: 100.2, MbitIn: 249.93}, networkStatistics[0])
	}
}

func TestRepository_GetBlockStorageVolumes(t *testing.T) {
	const apiResponse = `{"volumes":[{"uuid":"220887f0-db1a-76a9-2332-00004f589b19","name":"custom-2c3501ab-5a45-34e9-c289-00002b084a0c","sizeInGib":20,"type":"hdd","availabilityZone":"ams0","status":"available","nodeUuid":"76743b28-f779-3e68-6aa1-00007fbb911d","serial":"a4d857d3fe5e814f34bb"}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/block-storages", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetBlockStorageVolumes("k888k")
	require.NoError(t, err)

	if assert.Equal(t, 1, len(list)) {
		assert.Equal(t, "220887f0-db1a-76a9-2332-00004f589b19", list[0].UUID)
		assert.Equal(t, "custom-2c3501ab-5a45-34e9-c289-00002b084a0c", list[0].Name)
		assert.Equal(t, 20, list[0].SizeInGiB)
		assert.Equal(t, BlockStorageTypeHDD, list[0].Type)
		assert.Equal(t, "ams0", list[0].AvailabilityZone)
		assert.Equal(t, BlockStorageStatusAvailable, list[0].Status)
		assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", list[0].NodeUUID)
		assert.Equal(t, "a4d857d3fe5e814f34bb", list[0].Serial)
	}
}

func TestRepository_GetBlockStorageVolume(t *testing.T) {
	const apiResponse = `{"volume":{"uuid":"220887f0-db1a-76a9-2332-00004f589b19","name":"custom-2c3501ab-5a45-34e9-c289-00002b084a0c","sizeInGib":20,"type":"hdd","availabilityZone":"ams0","status":"available","nodeUuid":"76743b28-f779-3e68-6aa1-00007fbb911d","serial":"a4d857d3fe5e814f34bb"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/block-storages/custom-2c3501ab-5a45-34e9-c289-00002b084a0c", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	volume, err := repo.GetBlockStorageVolume("k888k", "custom-2c3501ab-5a45-34e9-c289-00002b084a0c")
	require.NoError(t, err)

	assert.Equal(t, "220887f0-db1a-76a9-2332-00004f589b19", volume.UUID)
	assert.Equal(t, "custom-2c3501ab-5a45-34e9-c289-00002b084a0c", volume.Name)
	assert.Equal(t, 20, volume.SizeInGiB)
	assert.Equal(t, BlockStorageTypeHDD, volume.Type)
	assert.Equal(t, "ams0", volume.AvailabilityZone)
	assert.Equal(t, BlockStorageStatusAvailable, volume.Status)
	assert.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", volume.NodeUUID)
	assert.Equal(t, "a4d857d3fe5e814f34bb", volume.Serial)
}

func TestRepository_AddBlockStorageVolume(t *testing.T) {
	const expectedRequestBody = `{"name":"custom-2c3501ab-5a45-34e9-c289-00002b084a0c","sizeInGib":200,"type":"hdd","availabilityZone":"ams0"}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/block-storages", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequestBody}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	order := BlockStorageOrder{
		ClusterName:      "k888k",
		Name:             "custom-2c3501ab-5a45-34e9-c289-00002b084a0c",
		SizeInGiB:        200,
		Type:             BlockStorageTypeHDD,
		AvailabilityZone: "ams0",
	}

	err := repo.AddBlockStorageVolume(order)
	require.NoError(t, err)
}

func TestRepository_UpdateBlockStorageVolume(t *testing.T) {
	const expectedRequest = `{"volume":{"uuid":"220887f0-db1a-76a9-2332-00004f589b19","name":"custom-2c3501ab-5a45-34e9-c289-00002b084a0c","clusterName":"k888k","sizeInGib":20,"type":"hdd","availabilityZone":"ams0","status":"available","nodeUuid":"76743b28-f779-3e68-6aa1-00007fbb911d","serial":"a4d857d3fe5e814f34bb"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/block-storages/custom-2c3501ab-5a45-34e9-c289-00002b084a0c", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.UpdateBlockStorageVolume(BlockStorage{
		ClusterName:      "k888k",
		UUID:             "220887f0-db1a-76a9-2332-00004f589b19",
		Name:             "custom-2c3501ab-5a45-34e9-c289-00002b084a0c",
		SizeInGiB:        20,
		Type:             BlockStorageTypeHDD,
		AvailabilityZone: "ams0",
		Status:           BlockStorageStatusAvailable,
		NodeUUID:         "76743b28-f779-3e68-6aa1-00007fbb911d",
		Serial:           "a4d857d3fe5e814f34bb",
	})

	require.NoError(t, err)
}

func TestRepository_RemoveBlockStorageVolume(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/block-storages/custom-2c3501ab-5a45-34e9-c289-00002b084a0c", ExpectedMethod: "DELETE", StatusCode: 204}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.RemoveBlockStorageVolume("k888k", "custom-2c3501ab-5a45-34e9-c289-00002b084a0c")
	require.NoError(t, err)
}

func TestRepository_GetBlockStorageStatistics(t *testing.T) {
	const apiResponse = `
	{
		"usage": [
			{
				"iopsRead": 0.27,
				"iopsWrite": 0.13,
				"date": 1500538995
			}
		]
	}`

	values := url.Values{
		"dateTimeStart": []string{"1500538995"},
		"dateTimeEnd":   []string{"1500542619"},
	}

	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/block-storages/custom-2c3501ab-5a45-34e9-c289-00002b084a0c/stats?" + values.Encode(),
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	statistics, err := repo.GetBlockStorageStatistics(
		"k888k",
		"custom-2c3501ab-5a45-34e9-c289-00002b084a0c",
		vps.UsagePeriod{
			TimeStart: 1500538995,
			TimeEnd:   1500542619,
		})
	require.NoError(t, err)

	if assert.NotEmpty(t, statistics) {
		assert.Equal(t, vps.UsageDataDisk{Date: 1500538995, IopsRead: 0.27, IopsWrite: 0.13}, statistics[0])
	}
}

func TestRepository_GetLoadBalancers(t *testing.T) {
	const apiResponse = `{"loadBalancers":[{"uuid":"220887f0-db1a-76a9-2332-00004f589b19","name":"lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80","status":"active","ipv4Address":"37.97.254.7","ipv6Address":"2a01:7c8:3:1337::1"}]}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	list, err := repo.GetLoadBalancers("k888k")
	require.NoError(t, err)

	if assert.Equal(t, 1, len(list)) {
		assert.Equal(t, "220887f0-db1a-76a9-2332-00004f589b19", list[0].UUID)
		assert.Equal(t, "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", list[0].Name)
		assert.Equal(t, LoadBalancerStatusActive, list[0].Status)
		assert.Equal(t, "37.97.254.7", list[0].IPv4Address.String())
		assert.Equal(t, "2a01:7c8:3:1337::1", list[0].IPv6Address.String())
	}
}

func TestRepository_GetLoadBalancer(t *testing.T) {
	const apiResponse = `{"loadBalancer":{"uuid":"220887f0-db1a-76a9-2332-00004f589b19","name":"lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80","status":"active","ipv4Address":"37.97.254.7","ipv6Address":"2a01:7c8:3:1337::1"}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers/lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	lb, err := repo.GetLoadBalancer("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80")
	require.NoError(t, err)

	assert.Equal(t, "220887f0-db1a-76a9-2332-00004f589b19", lb.UUID)
	assert.Equal(t, "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", lb.Name)
	assert.Equal(t, LoadBalancerStatusActive, lb.Status)
	assert.Equal(t, "37.97.254.7", lb.IPv4Address.String())
	assert.Equal(t, "2a01:7c8:3:1337::1", lb.IPv6Address.String())
}

func TestRepository_CreateLoadBalancer(t *testing.T) {
	const expectedRequestBody = `{"name":"lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80"}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers", ExpectedMethod: "POST", StatusCode: 201, ExpectedRequest: expectedRequestBody}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.CreateLoadBalancer("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80")
	require.NoError(t, err)
}

func TestRepository_UpdateLoadBalanmcer(t *testing.T) {
	const expectedRequest = `{"loadBalancerConfig":{"loadBalancingMode":"cookie","stickyCookieName":"PHPSESSID","healthCheckInterval":3000,"httpHealthCheckPath":"/status.php","httpHealthCheckPort":443,"httpHealthCheckSsl":true,"ipSetup":"ipv6to4","ptrRecord":"frontend.example.com","tlsMode":"tls12","ipAddresses":["10.3.37.1","10.3.38.1"],"portConfiguration":[{"name":"Website Traffic","sourcePort":80,"targetPort":8080,"mode":"http","endpointSslMode":"off"}]}}`

	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers/lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: expectedRequest}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.UpdateLoadBalancer("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", LoadBalancerConfig{
		LoadBalancingMode:   LoadBalancingModeCookie,
		StickyCookieName:    "PHPSESSID",
		HealthCheckInterval: 3000,
		HTTPHealthCheckPath: "/status.php",
		HTTPHealthCheckPort: 443,
		HTTPHealthCheckSSL:  true,
		IPSetup:             IPSetupIPv6to4,
		PTRRecord:           "frontend.example.com",
		TLSMode:             TLSModeMinTLS12,
		IPAddresses:         []net.IP{net.ParseIP("10.3.37.1"), net.ParseIP("10.3.38.1")},
		PortConfigurations: []PortConfiguration{
			{
				Name:            "Website Traffic",
				SourcePort:      80,
				TargetPort:      8080,
				Mode:            PortConfigurationModeHTTP,
				EndpointSSLMode: PortConfigurationEndpointSSLModeOff,
			},
		},
	})

	require.NoError(t, err)
}

func TestRepository_RemoveLoadBalancer(t *testing.T) {
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers/lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", ExpectedMethod: "DELETE", StatusCode: 204}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	err := repo.RemoveLoadBalancer("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80")
	require.NoError(t, err)
}

func TestRepository_GetLoadBalancerStatusReports(t *testing.T) {
	const apiResponse = `
	{
		"statusReports": [
			{
				"nodeUuid": "76743b28-f779-3e68-6aa1-00007fbb911d",
				"nodeIpAddress": "136.10.14.1",
				"port": 80,
				"ipVersion": 4,
				"loadBalancerName": "lb0",
				"loadBalancerIp": "136.144.151.255",
				"state": "up",
				"lastChange": "2019-09-29 16:51:18"
			}
		]
	}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers/lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80/status-reports", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	reports, err := repo.GetLoadBalancerStatusReports("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80")
	require.NoError(t, err)

	if assert.NotEmpty(t, reports) {
		require.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", reports[0].NodeUUID)
		require.Equal(t, net.ParseIP("136.10.14.1"), reports[0].NodeIPAddress)
		require.Equal(t, 80, reports[0].Port)
		require.Equal(t, 4, reports[0].IPVersion)
		require.Equal(t, "lb0", reports[0].LoadBalancerName)
		require.Equal(t, net.ParseIP("136.144.151.255"), reports[0].LoadBalancerIP)
		require.Equal(t, LoadBalancerStateUp, reports[0].State)
		require.Equal(t, "2019-09-29 16:51:18", reports[0].LastChange.Format(time.DateTime))
	}
}

func TestRepository_GetLoadBalancerStatusReportsForNode(t *testing.T) {
	const apiResponse = `
	{
		"statusReports": [
			{
				"nodeUuid": "76743b28-f779-3e68-6aa1-00007fbb911d",
				"nodeIpAddress": "136.10.14.1",
				"port": 80,
				"ipVersion": 4,
				"loadBalancerName": "lb0",
				"loadBalancerIp": "136.144.151.255",
				"state": "up",
				"lastChange": "2019-09-29 16:51:18"
			}
		]
	}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/load-balancers/lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80/status-reports/76743b28-f779-3e68-6aa1-00007fbb911d", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()
	repo := Repository{Client: *client}

	reports, err := repo.GetLoadBalancerStatusReportsForNode("k888k", "lb-bbb0ddf8-8aeb-4f35-85ff-4e198a0faf80", "76743b28-f779-3e68-6aa1-00007fbb911d")
	require.NoError(t, err)

	if assert.NotEmpty(t, reports) {
		require.Equal(t, "76743b28-f779-3e68-6aa1-00007fbb911d", reports[0].NodeUUID)
		require.Equal(t, net.ParseIP("136.10.14.1"), reports[0].NodeIPAddress)
		require.Equal(t, 80, reports[0].Port)
		require.Equal(t, 4, reports[0].IPVersion)
		require.Equal(t, "lb0", reports[0].LoadBalancerName)
		require.Equal(t, net.ParseIP("136.144.151.255"), reports[0].LoadBalancerIP)
		require.Equal(t, LoadBalancerStateUp, reports[0].State)
		require.Equal(t, "2019-09-29 16:51:18", reports[0].LastChange.Format(time.DateTime))
	}
}

func TestRepository_GetNodePoolTaints(t *testing.T) {
	const apiResponse = `{"taints":[{"key": "test-key", "value":"test-value", "effect":"NoSchedule", "modifiable": false}]}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82/taints", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()

	repo := Repository{Client: *client}
	taints, err := repo.GetTaints("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82")
	require.NoError(t, err)
	if assert.Equal(t, 1, len(taints)) {
		assert.Equal(t, "test-key", taints[0].Key)
		assert.Equal(t, "test-value", taints[0].Value)
		assert.Equal(t, "NoSchedule", taints[0].Effect)
		assert.Equal(t, false, taints[0].Modifiable)
	}
}

func TestRepository_GetNodePoolLabels(t *testing.T) {
	const apiResponse = `{"labels":[{"key": "test-key", "value":"test-value", "modifiable": false}]}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82/labels", ExpectedMethod: "GET", StatusCode: 200, Response: apiResponse}
	client, tearDown := server.GetClient()
	defer tearDown()

	repo := Repository{Client: *client}
	labels, err := repo.GetLabels("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82")
	require.NoError(t, err)
	if assert.Equal(t, 1, len(labels)) {
		assert.Equal(t, "test-key", labels[0].Key)
		assert.Equal(t, "test-value", labels[0].Value)
		assert.Equal(t, false, labels[0].Modifiable)
	}
}

func TestRepository_SetNodePoolLabels(t *testing.T) {
	const apiRequest = `{"labels":[{"key":"test-key","value":"test-value","modifiable":false}]}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82/labels", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: apiRequest}
	client, tearDown := server.GetClient()
	defer tearDown()

	repo := Repository{Client: *client}
	err := repo.SetLabels("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82", []Label{{Key: "test-key", Value: "test-value"}})
	require.NoError(t, err)
}

func TestRepository_SetNodePoolTaints(t *testing.T) {
	const apiRequest = `{"taints":[{"key":"test-key","value":"test-value","effect":"NoSchedule","modifiable":false}]}`
	server := testutil.MockServer{T: t, ExpectedURL: "/kubernetes/clusters/k888k/node-pools/402c2f84-c37d-9388-634d-00002b7c6a82/taints", ExpectedMethod: "PUT", StatusCode: 204, ExpectedRequest: apiRequest}
	client, tearDown := server.GetClient()
	defer tearDown()

	repo := Repository{Client: *client}
	err := repo.SetTaints("k888k", "402c2f84-c37d-9388-634d-00002b7c6a82", []Taint{{Key: "test-key", Value: "test-value", Effect: "NoSchedule"}})
	require.NoError(t, err)
}

func TestRepository_UpgradeCluster(t *testing.T) {
	const apiRequest = `{"action":"upgrade","version":"1.27.0"}`
	server := testutil.MockServer{
		T:               t,
		ExpectedURL:     "/kubernetes/clusters/k888k",
		ExpectedMethod:  "PATCH",
		StatusCode:      204,
		ExpectedRequest: apiRequest,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	err := repo.UpgradeCluster("k888k", "1.27.0")
	require.NoError(t, err)
}

func TestRepository_ResetCluster(t *testing.T) {
	const apiRequest = `{"action":"reset","confirmation":"k888k"}`
	server := testutil.MockServer{
		T:               t,
		ExpectedURL:     "/kubernetes/clusters/k888k",
		ExpectedMethod:  "PATCH",
		StatusCode:      204,
		ExpectedRequest: apiRequest,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	err := repo.ResetCluster("k888k", "k888k")
	require.NoError(t, err)
}

func TestRepository_TestGetReleases(t *testing.T) {
	const apiResponse = `{"releases":[{"version": "1.23.5","releaseDate": "2022-03-11","maintenanceModeDate": "2022-12-28","endOfLifeDate": "2023-02-28"}]}`
	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/releases",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	releases, err := repo.GetReleases()
	if assert.Equal(t, 1, len(releases)) {
		assert.Equal(t, "1.23.5", releases[0].Version)
		assert.Equal(t, "2022-03-11", releases[0].ReleaseDate.Format(dateOnlyFormat))
		assert.Equal(t, "2022-12-28", releases[0].MaintenanceModeDate.Format(dateOnlyFormat))
		assert.Equal(t, "2023-02-28", releases[0].EndOfLifeDate.Format(dateOnlyFormat))
	}
	require.NoError(t, err)
}

func TestRepository_TestGetRelease(t *testing.T) {
	const apiResponse = `{"release":{"version": "1.23.5","releaseDate": "2022-03-11","maintenanceModeDate": "2022-12-28","endOfLifeDate": "2023-02-28"}}`
	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/releases/1.23.5",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	release, err := repo.GetRelease("1.23.5")
	assert.Equal(t, "1.23.5", release.Version)
	assert.Equal(t, "2022-03-11", release.ReleaseDate.Format(dateOnlyFormat))
	assert.Equal(t, "2022-12-28", release.MaintenanceModeDate.Format(dateOnlyFormat))
	assert.Equal(t, "2023-02-28", release.EndOfLifeDate.Format(dateOnlyFormat))
	require.NoError(t, err)
}

func TestRepository_TestGetCompatibleReleases(t *testing.T) {
	const apiResponse = `{"releases":[{"isCompatibleUpgrade":true,"version": "1.23.5","releaseDate": "2022-03-11","maintenanceModeDate": "2022-12-28","endOfLifeDate": "2023-02-28"}]}`
	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/releases",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	releases, err := repo.GetCompatibleReleases("k888k")
	if assert.Equal(t, 1, len(releases)) {
		assert.Equal(t, "1.23.5", releases[0].Version)
		assert.Equal(t, "2022-03-11", releases[0].ReleaseDate.Format(dateOnlyFormat))
		assert.Equal(t, "2022-12-28", releases[0].MaintenanceModeDate.Format(dateOnlyFormat))
		assert.Equal(t, "2023-02-28", releases[0].EndOfLifeDate.Format(dateOnlyFormat))
	}
	require.NoError(t, err)
}

func TestRepository_TestGetCompatibleRelease(t *testing.T) {
	const apiResponse = `{"release":{"isCompatibleUpgrade":true,"version": "1.23.5","releaseDate": "2022-03-11","maintenanceModeDate": "2022-12-28","endOfLifeDate": "2023-02-28"}}`
	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/releases/1.23.5",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	release, err := repo.GetCompatibleRelease("k888k", "1.23.5")
	assert.Equal(t, "1.23.5", release.Version)
	assert.Equal(t, "2022-03-11", release.ReleaseDate.Format(dateOnlyFormat))
	assert.Equal(t, "2022-12-28", release.MaintenanceModeDate.Format(dateOnlyFormat))
	assert.Equal(t, "2023-02-28", release.EndOfLifeDate.Format(dateOnlyFormat))
	require.NoError(t, err)
}

func TestRepository_TestGetEvents(t *testing.T) {
	const apiResponse = `
	{
		"events": [
			{
				"name": "kube-proxy-g9ldg.175d7f60d241f2c8",
				"namespace": "default",
				"type": "Warning",
				"message": "Node is not ready",
				"reason": "NodeNotReady",
				"count": 6,
				"creationTimestamp": 1683641890,
				"firstTimestamp": 1683641890,
				"lastTimestamp": 1683641890,
				"involvedObjectKind": "Pod",
				"involvedObjectName": "kube-proxy-g9ldg",
				"sourceComponent": "kubelet"
			}
		]
	}`

	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/events",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	events, err := repo.GetEvents("k888k")

	require.NoError(t, err)
	if assert.NotEmpty(t, events) {
		require.Equal(t, "kube-proxy-g9ldg.175d7f60d241f2c8", events[0].Name)
		require.Equal(t, "default", events[0].Namespace)
		require.Equal(t, "Warning", events[0].Type)
		require.Equal(t, "Node is not ready", events[0].Message)
		require.Equal(t, "NodeNotReady", events[0].Reason)
		require.Equal(t, 6, events[0].Count)
		require.Equal(t, 1683641890, events[0].CreationTimestamp)
		require.Equal(t, 1683641890, events[0].FirstTimestamp)
		require.Equal(t, 1683641890, events[0].LastTimestamp)
		require.Equal(t, "Pod", events[0].InvolvedObjectKind)
		require.Equal(t, "kube-proxy-g9ldg", events[0].InvolvedObjectName)
		require.Equal(t, "kubelet", events[0].SourceComponent)
	}
}

func TestRepository_TestGetEventsByNamespace(t *testing.T) {
	const apiResponse = `
	{
		"events": [
			{
				"name": "kube-proxy-g9ldg.175d7f60d241f2c8",
				"namespace": "default",
				"type": "Warning",
				"message": "Node is not ready",
				"reason": "NodeNotReady",
				"count": 6,
				"creationTimestamp": 1683641890,
				"firstTimestamp": 1683641890,
				"lastTimestamp": 1683641890,
				"involvedObjectKind": "Pod",
				"involvedObjectName": "kube-proxy-g9ldg",
				"sourceComponent": "kubelet"
			}
		]
	}
	`

	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/events?namespace=default",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	events, err := repo.GetEventsByNamespace("k888k", "default")

	require.NoError(t, err)
	if assert.NotEmpty(t, events) {
		require.Equal(t, "kube-proxy-g9ldg.175d7f60d241f2c8", events[0].Name)
		require.Equal(t, "default", events[0].Namespace)
		require.Equal(t, "Warning", events[0].Type)
		require.Equal(t, "Node is not ready", events[0].Message)
		require.Equal(t, "NodeNotReady", events[0].Reason)
		require.Equal(t, 6, events[0].Count)
		require.Equal(t, 1683641890, events[0].CreationTimestamp)
		require.Equal(t, 1683641890, events[0].FirstTimestamp)
		require.Equal(t, 1683641890, events[0].LastTimestamp)
		require.Equal(t, "Pod", events[0].InvolvedObjectKind)
		require.Equal(t, "kube-proxy-g9ldg", events[0].InvolvedObjectName)
		require.Equal(t, "kubelet", events[0].SourceComponent)
	}
}

func TestRepository_TestGetEventByName(t *testing.T) {
	const apiResponse = `
	{
		"event": {
			"name": "kube-proxy-g9ldg.175d7f60d241f2c8",
			"namespace": "default",
			"type": "Warning",
			"message": "Node is not ready",
			"reason": "NodeNotReady",
			"count": 6,
			"creationTimestamp": 1683641890,
			"firstTimestamp": 1683641890,
			"lastTimestamp": 1683641890,
			"involvedObjectKind": "Pod",
			"involvedObjectName": "kube-proxy-g9ldg",
			"sourceComponent": "kubelet"
		}
	}
	`
	server := testutil.MockServer{
		T:              t,
		ExpectedURL:    "/kubernetes/clusters/k888k/events/kube-proxy-g9ldg.175d7f60d241f2c8",
		ExpectedMethod: "GET",
		StatusCode:     200,
		Response:       apiResponse,
	}

	client, teardown := server.GetClient()
	defer teardown()

	repo := Repository{Client: *client}
	event, err := repo.GetEventByName("k888k", "kube-proxy-g9ldg.175d7f60d241f2c8")

	require.NoError(t, err)
	require.Equal(t, "kube-proxy-g9ldg.175d7f60d241f2c8", event.Name)
	require.Equal(t, "default", event.Namespace)
	require.Equal(t, "Warning", event.Type)
	require.Equal(t, "Node is not ready", event.Message)
	require.Equal(t, "NodeNotReady", event.Reason)
	require.Equal(t, 6, event.Count)
	require.Equal(t, 1683641890, event.CreationTimestamp)
	require.Equal(t, 1683641890, event.FirstTimestamp)
	require.Equal(t, 1683641890, event.LastTimestamp)
	require.Equal(t, "Pod", event.InvolvedObjectKind)
	require.Equal(t, "kube-proxy-g9ldg", event.InvolvedObjectName)
	require.Equal(t, "kubelet", event.SourceComponent)
}
