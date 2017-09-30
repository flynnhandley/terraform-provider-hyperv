package hyperv

// A driver is able to talk to HyperV and perform certain
// operations with it. Some of the operations on here may seem overly
// specific, but they were built specifically in mind to handle features
// of the HyperV builder for Packer, and to abstract differences in
// versions out of the builder steps, so sometimes the methods are
// extremely specific.
type Driver interface {
	// Downloads a file to remote host
	Download(string, string, string, string) (string, error)

	// Executes arbitrary commands on the remote host
	InvokeCommand(string, map[string]string) (string, error)

	// Copies a file from the local host to the remote host using bits transfer
	PutFile(string, string) error

	// Copies a file from the remote host to the local host using bits transfer
	GetFile(string, string) error

	// Returns error if a WinRM connection cannot be established to the remote host
	TestConnectivity() error

	// Returns hash using given algorithm
	Hash(string, string) (string, error)

	// Checks if the VM named is running.
	IsRunning(string) (bool, error)

	// Checks if the VM named is off.
	IsOff(string) (bool, error)

	//How long has VM been on
	Uptime(vmName string) (uint64, error)

	// Start starts a VM specified by the name given.
	Start(string) error

	// Stop stops a VM specified by the name given.
	Stop(string) error

	// Creates a VHDX in the VM's configuration directoy and attaches to the VM
	NewVhd(string, string, int64) (string, error)

	NewDiskFromImagePath(string, string, string) (string, error)

	NewDifferencingDisk(string, string, string) (string, error)
	NewDiskFromImageURL(string, string, string) (string, error)

	// Verify checks to make sure that this driver should function
	// properly. If there is any indication the driver can't function,
	// this will return an error.
	Verify() error

	// Finds the MAC address of the NIC nic0
	Mac(string) (string, error)

	SetVirtualMachineRemoveNetworkBoot(string) error

	// Finds the IP address of a VM connected that uses DHCP by its MAC address
	IpAddress(string) (string, error)

	// Finds the hostname for the ip address
	GetHostName(string) (string, error)

	// Finds the IP address of a host adapter connected to switch
	GetHostAdapterIpAddressForSwitch(string) (string, error)

	// Type scan codes to virtual keyboard of vm
	TypeScanCodes(string, string) error

	//Get the ip address for network adaptor
	GetVirtualMachineNetworkAdapterAddress(string) (string, error)

	//Set the vlan to use for switch
	SetNetworkAdapterVlanId(string, string) error

	//Set the vlan to use for machine
	SetVirtualMachineVlanId(string, string) error

	UntagVirtualMachineNetworkAdapterVlan(string, string) error

	CreateExternalVirtualSwitch(string, string) error

	GetVirtualSwitchId(map[string]string) (string, error)

	GetVirtualMachineSwitchName(string) (string, error)

	AddVMNetworkAdapter(string, string, string, string) error

	ConnectVirtualMachineNetworkAdapterToSwitch(string, string) error

	CreateVirtualSwitch(string, string) (string, error)

	DeleteVirtualSwitch(string) error

	CreateVirtualMachine(string, string, int64, string, int) (string, error)

	DeleteVirtualMachine(string) error

	GetVirtualMachineId(map[string]string) (string, error)

	SetVirtualMachineCpuCount(string, int) error

	SetVirtualMachineMacSpoofing(string, bool) error

	SetVirtualMachineDynamicMemory(string, bool) error

	SetVirtualMachineSecureBoot(string, bool) error

	SetVirtualMachineVirtualizationExtensions(string, bool) error

	EnableVirtualMachineIntegrationService(string, string) error

	ExportVirtualMachine(string, string) error

	CompactDisks(string, string) error

	CopyExportedVirtualMachine(string, string, string, string) error

	RestartVirtualMachine(string) error

	CreateDvdDrive(string, string, uint) (uint, uint, error)

	MountDvdDrive(string, string, uint, uint) error

	SetBootDvdDrive(string, uint, uint, uint) error

	UnmountDvdDrive(string, uint, uint) error

	DeleteDvdDrive(string, uint, uint) error

	MountFloppyDrive(string, string) error

	UnmountFloppyDrive(string) error
}
