package tasks

import (
	"net"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/baremetal/types"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type IBaremetal interface {
	GetId() string
	GetZoneId() string
	GetTaskQueue() *TaskQueue
	GetSSHConfig() (*types.SSHConfig, error)
	TestSSHConfig() bool
	GetAdminNic() *types.Nic
	GetName() string
	GetClientSession() *mcclient.ClientSession
	SaveDesc(desc jsonutils.JSONObject) error
	GetNicByMac(net.HardwareAddr) *types.Nic
	GetRawIPMIConfig() *types.IPMIInfo
	GetIPMINic(mac net.HardwareAddr) *types.Nic
	SetExistingIPMIIPAddr(ipAddr string)

	SyncStatus(status, reason string)
	AutoSyncStatus()
	SyncAllStatus(status string)
	AutoSyncAllStatus()

	GetPowerStatus() (string, error)
	DoPowerShutdown(soft bool) error
	DoPXEBoot() error
	DoDiskBoot() error
}
