package hyperv

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/flynnhandley/psremote/hvremote"
	"github.com/hashicorp/packer/common/powershell"
)

// PS5Driver because I have to
type PS5Driver struct {
	hyperv *hvremote.HyperVCmd
}

// NewPS5Driver because I have to
func NewPS5Driver(userName, password, computerName string) (Driver, error) {
	appliesTo := "Applies to Windows 8.1, Windows PowerShell 4.0, Windows Server 2012 R2 only"

	// Check this is Windows
	if runtime.GOOS != "windows" {
		err := fmt.Errorf("%s", appliesTo)
		return nil, err
	}

	ps4Driver := &PS5Driver{}
	var hv, _ = hvremote.NewHyperVCmd(userName, password, computerName)
	ps4Driver.hyperv = hv

	return ps4Driver, nil
}

// Returns hash using given algorithm
func (d *PS5Driver) InvokeCommand(scriptBlock string, params map[string]string) (string, error) {
	return d.hyperv.InvokeCommand(scriptBlock, params)
}

// Returns hash using given algorithm
func (d *PS5Driver) TestConnectivity() error {
	return d.hyperv.TestConnectivity()
}

// Returns hash using given algorithm
func (d *PS5Driver) AttachBootVHD(vmID, source string) (string, error) {
	return d.hyperv.AttachBootVHD(vmID, source)
}

// Returns hash using given algorithm
func (d *PS5Driver) NewVhd(vmID, vhdName string, diskSize int64) (string, error) {
	return d.hyperv.NewVhd(vmID, vhdName, diskSize)
}

// Sends a file to the remote session
func (d *PS5Driver) PutFile(source, dest string) error {
	return d.hyperv.PutFile(source, dest)
}

// Gets a file from the remote session
func (d *PS5Driver) GetFile(source, dest string) error {
	return d.hyperv.GetFile(source, dest)
}

// Returns hash using given algorithm
func (d *PS5Driver) Download(source, dest, hash, algorithm string) (string, error) {
	return d.hyperv.Download(source, dest, hash, algorithm)
}

func (d *PS5Driver) IsRunning(vmName string) (bool, error) {
	return d.hyperv.IsRunning(vmName)
}

func (d *PS5Driver) IsOff(vmName string) (bool, error) {
	return d.hyperv.IsOff(vmName)
}

func (d *PS5Driver) Uptime(vmName string) (uint64, error) {
	return d.hyperv.Uptime(vmName)
}

// Start starts a VM specified by the name given.
func (d *PS5Driver) Start(vmName string) error {
	return d.hyperv.StartVirtualMachine(vmName)
}

// Stop stops a VM specified by the name given.
func (d *PS5Driver) Stop(vmName string) error {
	return d.hyperv.StopVirtualMachine(vmName)
}

func (d *PS5Driver) Verify() error {

	if err := d.verifyPSVersion(); err != nil {
		return err
	}

	if err := d.verifyPSHypervModule(); err != nil {
		return err
	}

	if err := d.verifyHypervPermissions(); err != nil {
		return err
	}

	return nil
}

// Get mac address for VM.
func (d *PS5Driver) Mac(vmName string) (string, error) {
	res, err := d.hyperv.Mac(vmName)

	if err != nil {
		return res, err
	}

	if res == "" {
		err := fmt.Errorf("%s", "No mac address.")
		return res, err
	}

	return res, err
}

// Get ip address for mac address.
func (d *PS5Driver) IpAddress(mac string) (string, error) {
	res, err := d.hyperv.IpAddress(mac)

	if err != nil {
		return res, err
	}

	if res == "" {
		err := fmt.Errorf("%s", "No ip address.")
		return res, err
	}
	return res, err
}

// Get host name from ip address
func (d *PS5Driver) GetHostName(ip string) (string, error) {
	return powershell.GetHostName(ip)
}

// Finds the IP address of a host adapter connected to switch
func (d *PS5Driver) GetHostAdapterIpAddressForSwitch(switchName string) (string, error) {
	res, err := d.hyperv.GetHostAdapterIpAddressForSwitch(switchName)

	if err != nil {
		return res, err
	}

	if res == "" {
		err := fmt.Errorf("%s", "No ip address.")
		return res, err
	}
	return res, err
}

// Type scan codes to virtual keyboard of vm
func (d *PS5Driver) TypeScanCodes(vmName string, scanCodes string) error {
	return d.hyperv.TypeScanCodes(vmName, scanCodes)
}

// Returns hash using given algorithm
func (d *PS5Driver) Hash(path, algorithm string) (string, error) {
	return d.hyperv.Hash(path, algorithm)
}

// Get network adapter address
func (d *PS5Driver) GetVirtualMachineNetworkAdapterAddress(vmName string) (string, error) {
	return d.hyperv.GetVirtualMachineNetworkAdapterAddress(vmName)
}

//Set the vlan to use for switch
func (d *PS5Driver) SetNetworkAdapterVlanId(switchName string, vlanId string) error {
	return d.hyperv.SetNetworkAdapterVlanId(switchName, vlanId)
}

//Set the vlan to use for machine
func (d *PS5Driver) SetVirtualMachineVlanId(vmName string, vlanId string) error {
	return d.hyperv.SetVirtualMachineVlanId(vmName, vlanId)
}

func (d *PS5Driver) UntagVirtualMachineNetworkAdapterVlan(vmName string, switchName string) error {
	return d.hyperv.UntagVirtualMachineNetworkAdapterVlan(vmName, switchName)
}

func (d *PS5Driver) CreateExternalVirtualSwitch(vmName string, switchName string) error {
	return d.hyperv.CreateExternalVirtualSwitch(vmName, switchName)
}

func (d *PS5Driver) GetVirtualMachineSwitchName(vmName string) (string, error) {
	return d.hyperv.GetVirtualMachineSwitchName(vmName)
}

func (d *PS5Driver) GetVirtualMachineId(params map[string]string) (string, error) {
	return d.hyperv.GetVirtualMachineId(params)
}

func (d *PS5Driver) GetVirtualSwitchId(params map[string]string) (string, error) {
	return d.hyperv.GetVirtualSwitchId(params)
}

func (d *PS5Driver) AddVMNetworkAdapter(vmId, name, switchName, vlanId string) error {
	return d.hyperv.AddVMNetworkAdapter(vmId, name, switchName, vlanId)
}

func (d *PS5Driver) ConnectVirtualMachineNetworkAdapterToSwitch(vmName string, switchName string) error {
	return d.hyperv.ConnectVirtualMachineNetworkAdapterToSwitch(vmName, switchName)
}

func (d *PS5Driver) DeleteVirtualSwitch(switchId string) error {
	return d.hyperv.DeleteVirtualSwitch(switchId)
}

func (d *PS5Driver) CreateVirtualSwitch(switchName string, switchType string) (string, error) {
	return d.hyperv.CreateVirtualSwitch(switchName, switchType)
}

func (d *PS5Driver) CreateVirtualMachine(vmName string, path string, ram int64, switchName string, generation int) (string, error) {
	return d.hyperv.CreateVirtualMachine(vmName, path, ram, switchName, generation)
}

func (d *PS5Driver) DeleteVirtualMachine(vmId string) error {
	return d.hyperv.DeleteVirtualMachine(vmId)
}

func (d *PS5Driver) SetVirtualMachineCpuCount(vmId string, cpu int) error {
	return d.hyperv.SetVirtualMachineCpuCount(vmId, cpu)
}

func (d *PS5Driver) SetVirtualMachineMacSpoofing(vmName string, enable bool) error {
	return d.hyperv.SetVirtualMachineMacSpoofing(vmName, enable)
}

func (d *PS5Driver) SetVirtualMachineDynamicMemory(vmName string, enable bool) error {
	return d.hyperv.SetVirtualMachineDynamicMemory(vmName, enable)
}

func (d *PS5Driver) SetVirtualMachineSecureBoot(vmName string, enable bool) error {
	return d.hyperv.SetVirtualMachineSecureBoot(vmName, enable)
}

func (d *PS5Driver) SetVirtualMachineVirtualizationExtensions(vmName string, enable bool) error {
	return d.hyperv.SetVirtualMachineVirtualizationExtensions(vmName, enable)
}

func (d *PS5Driver) EnableVirtualMachineIntegrationService(vmName string, integrationServiceName string) error {
	return d.hyperv.EnableVirtualMachineIntegrationService(vmName, integrationServiceName)
}

func (d *PS5Driver) ExportVirtualMachine(vmName string, path string) error {
	return d.hyperv.ExportVirtualMachine(vmName, path)
}

func (d *PS5Driver) CompactDisks(expPath string, vhdDir string) error {
	return d.hyperv.CompactDisks(expPath, vhdDir)
}

func (d *PS5Driver) CopyExportedVirtualMachine(expPath string, outputPath string, vhdDir string, vmDir string) error {
	return d.hyperv.CopyExportedVirtualMachine(expPath, outputPath, vhdDir, vmDir)
}

func (d *PS5Driver) RestartVirtualMachine(vmName string) error {
	return d.hyperv.RestartVirtualMachine(vmName)
}

func (d *PS5Driver) CreateDvdDrive(vmName string, isoPath string, generation uint) (uint, uint, error) {
	return d.hyperv.CreateDvdDrive(vmName, isoPath, generation)
}

func (d *PS5Driver) MountDvdDrive(vmName string, path string, controllerNumber uint, controllerLocation uint) error {
	return d.hyperv.MountDvdDrive(vmName, path, controllerNumber, controllerLocation)
}

func (d *PS5Driver) SetBootDvdDrive(vmName string, controllerNumber uint, controllerLocation uint, generation uint) error {
	return d.hyperv.SetBootDvdDrive(vmName, controllerNumber, controllerLocation, generation)
}

func (d *PS5Driver) UnmountDvdDrive(vmName string, controllerNumber uint, controllerLocation uint) error {
	return d.hyperv.UnmountDvdDrive(vmName, controllerNumber, controllerLocation)
}

func (d *PS5Driver) DeleteDvdDrive(vmName string, controllerNumber uint, controllerLocation uint) error {
	return d.hyperv.DeleteDvdDrive(vmName, controllerNumber, controllerLocation)
}

func (d *PS5Driver) MountFloppyDrive(vmName string, path string) error {
	return d.hyperv.MountFloppyDrive(vmName, path)
}

func (d *PS5Driver) UnmountFloppyDrive(vmName string) error {
	return d.hyperv.UnmountFloppyDrive(vmName)
}

func (d *PS5Driver) verifyPSVersion() error {

	log.Printf("Enter method: %s", "verifyPSVersion")
	// check PS is available and is of proper version
	versionCmd := "$host.version.Major"

	var ps powershell.PowerShellCmd
	cmdOut, err := ps.Output(versionCmd)
	if err != nil {
		return err
	}

	versionOutput := strings.TrimSpace(cmdOut)
	log.Printf("%s output: %s", versionCmd, versionOutput)

	ver, err := strconv.ParseInt(versionOutput, 10, 32)

	if err != nil {
		return err
	}

	if ver < 4 {
		err := fmt.Errorf("%s", "Windows PowerShell version 4.0 or higher is expected")
		return err
	}

	return nil
}

func (d *PS5Driver) verifyPSHypervModule() error {

	log.Printf("Enter method: %s", "verifyPSHypervModule")

	versionCmd := "function foo(){try{ $commands = Get-Command -Module Hyper-V;if($commands.Length -eq 0){return $false} }catch{return $false}; return $true} foo"

	var ps powershell.PowerShellCmd
	cmdOut, err := ps.Output(versionCmd)
	if err != nil {
		return err
	}

	res := strings.TrimSpace(cmdOut)

	if res == "False" {
		err := fmt.Errorf("%s", "PS Hyper-V module is not loaded. Make sure Hyper-V feature is on.")
		return err
	}

	return nil
}

func (d *PS5Driver) verifyHypervPermissions() error {

	log.Printf("Enter method: %s", "verifyHypervPermissions")

	//SID:S-1-5-32-578 = 'BUILTIN\Hyper-V Administrators'
	//https://support.microsoft.com/en-us/help/243330/well-known-security-identifiers-in-windows-operating-systems
	hypervAdminCmd := "([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole('S-1-5-32-578')"

	var ps powershell.PowerShellCmd
	cmdOut, err := ps.Output(hypervAdminCmd)
	if err != nil {
		return err
	}

	res := strings.TrimSpace(cmdOut)

	if res == "False" {
		isAdmin, _ := powershell.IsCurrentUserAnAdministrator()

		if !isAdmin {
			err := fmt.Errorf("%s", "Current user is not a member of 'Hyper-V Administrators' or 'Administrators' group")
			return err
		}
	}

	return nil
}
