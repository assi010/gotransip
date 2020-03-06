package vps

import (
	"github.com/transip/gotransip/v6/ipaddress"
	"github.com/transip/gotransip/v6/product"
)

// vpsWrapper struct contains a Vps in it,
// this is solely used for unmarshalling/marshalling
type vpsWrapper struct {
	Vps Vps `json:"vps"`
}

// vpssWrapper struct contains a list of Vpses in it,
// this is solely used for unmarshalling/marshalling
type vpssWrapper struct {
	Vpss []Vps `json:"vpss"`
}

// vpssOrderWrapper struct contains a list of VpsOrders in it,
// this is solely used for marshalling
type vpssOrderWrapper struct {
	Orders []VpsOrder `json:"vpss"`
}

// cloneRequest is solely used for marshalling a vpsName and an availabilityZone
type cloneRequest struct {
	VpsName          string `json:"vpsName"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// actionWrapper struct contains an action in it,
// this is solely used for marshalling
type actionWrapper struct {
	Action string `json:"action"`
}

// handoverRequest is used to request a handover, this is solely used for marshalling
type handoverRequest struct {
	Action             string `json:"action"`
	TargetCustomerName string `json:"targetCustomerName"`
}

// usageWrapper struct contains Usage in it,
// this is solely used for unmarshalling
type usageWrapper struct {
	Usage Usage `json:"usage"`
}

// vncDataWrapper struct contains VncData in it,
// this is solely used for unmarshalling
type vncDataWrapper struct {
	VncData VncData `json:"vncData"`
}

// addonsWrapper struct contains a list with Addons in it,
// this is solely used for unmarshalling
type addonsWrapper struct {
	Addons Addons `json:"addons"`
}

// addonOrderRequest struct contains a list with Addons in it,
// this is solely used for marshalling
type addonOrderRequest struct {
	Addons []string `json:"addons"`
}

// upgradeRequest struct contains a Product Name in it,
// this is solely used for marshalling
type upgradeRequest struct {
	ProductName string `json:"productName"`
}

// upgradesWrapper struct contains a list with Products in it,
// this is solely used for marshalling
type upgradesWrapper struct {
	Upgrades []product.Product `json:"upgrades"`
}

// operatingSystemsWrapper struct contains a list with OperatingSystems in it,
// this is solely used for marshalling
type operatingSystemsWrapper struct {
	OperatingSystems []OperatingSystem `json:"operatingSystems"`
}

// ipAddressWrapper struct contains an IPAddress in it,
// this is solely used for unmarshalling
type ipAddressWrapper struct {
	IPAddress ipaddress.IPAddress `json:"ipAddress"`
}

// snapshotWrapper struct contains a Snapshot in it,
// this is solely used for unmarshalling
type snapshotWrapper struct {
	Snapshot Snapshot `json:"snapshot"`
}

// snapshotWrapper struct contains a list of Snapshots in it,
// this is solely used for unmarshalling
type snapshotsWrapper struct {
	Snapshots []Snapshot `json:"snapshot"`
}

// addIpRequest struct contains an IPAddress in it,
// this is solely used for marshalling
type addIpRequest struct {
	IPAddress string `json:"ipAddress"`
}

// createSnapshotRequest is used to marshal a request for creating a snapshot on a vps
// this is solely used for marshalling
type createSnapshotRequest struct {
	Description    string `json:"description"`
	ShouldStartVps bool   `json:"shouldStartVps"`
}

// revertSnapshotRequest is used to marshal a request for reverting a snapshot to a vps
// this is solely used for marshalling
type revertSnapshotRequest struct {
	DestinationVpsName string `json:"destinationVpsName"`
}

// installRequest struct contains a list with OperatingSystems in it,
// this is solely used for marshalling
type installRequest struct {
	OperatingSystemName string `json:"operatingSystemName"`
	Hostname            string `json:"hostname,omitempty"`
	Base64InstallText   string `json:"base64InstallText,omitempty"`
}

// Vps struct for Vps
type Vps struct {
	// The custom tags added to this VPS
	Tags []string `json:"tags,omitempty"`
	// The name of the availability zone the VPS is in
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// The VPS cpu count
	Cpus float32 `json:"cpus,omitempty"`
	// The amount of snapshots that is used on this VPS
	CurrentSnapshots float32 `json:"currentSnapshots,omitempty"`
	// The name that can be set by customer
	Description string `json:"description,omitempty"`
	// The VPS disk size in kB
	DiskSize float32 `json:"diskSize,omitempty"`
	// The VPS main ipAddress
	IpAddress string `json:"ipAddress,omitempty"`
	// If the VPS is administratively blocked
	IsBlocked bool `json:"isBlocked,omitempty"`
	// If this VPS is locked by the customer
	IsCustomerLocked bool `json:"isCustomerLocked,omitempty"`
	// Whether or not another process is already doing stuff with this VPS
	IsLocked bool `json:"isLocked,omitempty"`
	// The VPS macaddress
	MacAddress string `json:"macAddress,omitempty"`
	// The maximum amount of snapshots for this VPS
	MaxSnapshots float32 `json:"maxSnapshots,omitempty"`
	// The VPS memory size in kB
	MemorySize float32 `json:"memorySize,omitempty"`
	// The unique VPS name
	Name string `json:"name"`
	// The VPS OperatingSystem
	OperatingSystem string `json:"operatingSystem,omitempty"`
	// The product name
	ProductName string `json:"productName,omitempty"`
	// The VPS status, either 'created', 'installing', 'running', 'stopped' or 'paused'
	Status string `json:"status,omitempty"`
}

// VncData struct for VpsVncData
type VncData struct {
	// Location of the VNC Proxy
	Host string `json:"host,omitempty"`
	// Password to setup up the VNC connection (changes dynamically)
	Password string `json:"password,omitempty"`
	// Websocket path including the token
	Path string `json:"path,omitempty"`
	// token to identify the VPS to connect to (changes dynamically)
	Token string `json:"token,omitempty"`
	// Complete websocket URL
	Url string `json:"url,omitempty"`
}

// VpsUsageDataNetwork struct for VpsUsageDataNetwork
type VpsUsageDataNetwork struct {
	// Date of the entry, by default in UNIX timestamp format
	Date float32 `json:"date"`
	// The amount of inbound traffic in Mbps for this usage entry
	MbitIn float32 `json:"mbitIn"`
	// The amount of outbound traffic in Mbps for this usage entry
	MbitOut float32 `json:"mbitOut"`
}

// UsageDataDisk struct for UsageDataDisk
type UsageDataDisk struct {
	// Date of the entry, by default in UNIX timestamp format
	Date float32 `json:"date"`
	// The read IOPS for this entry
	IopsRead float32 `json:"iopsRead"`
	// The write IOPS for this entry
	IopsWrite float32 `json:"iopsWrite"`
}

// VpsUsageDataCpu struct for VpsUsageDataCpu
type VpsUsageDataCpu struct {
	// Date of the entry, by default in UNIX timestamp format
	Date float32 `json:"date"`
	// The percentage of CPU usage for this entry
	Percentage float32 `json:"percentage"`
}

// VpsOrder struct for VpsOrder
type VpsOrder struct {
	// Name of the product
	ProductName string `json:"productName"`
	// The name of the operating system to install
	OperatingSystem string `json:"operatingSystem"`
	// The name of the availability zone where the vps should be created
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// Array with additional addons
	Addons []string `json:"addons,omitempty"`
	// The name for the host, only needed for installing a preinstallable control panel image
	Hostname string `json:"hostname,omitempty"`
	// The description of the VPS
	Description string `json:"description,omitempty"`
	// Base64 encoded preseed / kickstart instructions, when installing unattended
	Base64InstallText string `json:"base64InstallText,omitempty"`
}

// Addons struct for Addons
type Addons struct {
	// A list of non cancellable active addons
	Active []product.Product `json:"active,omitempty"`
	// A list of available addons that you can order
	Available []product.Product `json:"available,omitempty"`
	// A list of addons that you can cancel
	Cancellable []product.Product `json:"cancellable,omitempty"`
}

// Backup struct for Backup
type Backup struct {
	// The name of the availability zone the backup is in
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// The backup creation date
	DateTimeCreate string `json:"dateTimeCreate"`
	// The backup disk size in kB
	DiskSize float32 `json:"diskSize"`
	// The backup id
	Id float32 `json:"id"`
	// The backup operatingSystem
	OperatingSystem string `json:"operatingSystem"`
	// Status of the backup ('active', 'creating', 'reverting', 'deleting', 'pendingDeletion', 'syncing', 'moving')
	Status string `json:"status,omitempty"`
}

// Snapshots struct for Snapshots
type Snapshots struct {
	// The snapshot creation date
	Snapshots []Snapshot `json:"snapshots,omitempty"`
}

// Snapshot struct for Snapshot
type Snapshot struct {
	// The snapshot creation date
	DateTimeCreate string `json:"dateTimeCreate,omitempty"`
	// The snapshot description
	Description string `json:"description,omitempty"`
	// The size of the snapshot in kB
	DiskSize float32 `json:"diskSize,omitempty"`
	// The snapshot name
	Name string `json:"name,omitempty"`
	// The snapshot OperatingSystem
	OperatingSystem string `json:"operatingSystem,omitempty"`
	// The snapshot status ('active', 'creating', 'reverting', 'deleting', 'pendingDeletion', 'syncing', 'moving')
	Status string `json:"status,omitempty"`
}

// OperatingSystems struct for OperatingSystems
type OperatingSystems struct {
	// OperatingSystems
	OperatingSystems string `json:"operatingSystems,omitempty"`
}

// OperatingSystem struct for OperatingSystem
type OperatingSystem struct {
	// Description
	Description string `json:"description,omitempty"`
	// Is a preinstallable image
	IsPreinstallableImage bool `json:"isPreinstallableImage,omitempty"`
	// The operating system name
	Name string `json:"name"`
	// The monthly price of the operating system in cents
	Price float32 `json:"price,omitempty"`
	// The version of the operating system
	Version string `json:"version,omitempty"`
}

// Usage struct for Usage
type Usage struct {
	Cpu     []VpsUsageDataCpu     `json:"cpu"`
	Disk    []UsageDataDisk       `json:"disk"`
	Network []VpsUsageDataNetwork `json:"network"`
}

// Upgrades struct for Upgrades
type Upgrades struct {
	Upgrades []product.Product `json:"upgrades"`
}
