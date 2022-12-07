package dan

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	current "github.com/containernetworking/cni/pkg/types/100"

	"github.com/containernetworking/cni/pkg/types"
)

type DirectAttachableNetworkType string

const DirectAttachableNetworkTypeTap DirectAttachableNetworkType = "tap"
const DirectAttachableNetworkTypePassthrough DirectAttachableNetworkType = "passthrough"
const DirectAttachableNetworkTypeDPDK DirectAttachableNetworkType = "dpdk"

type DirectAttachableNetwork struct {
	// tell runtime about the network hardware parts
	NetworkType            DirectAttachableNetworkType `json:"networkType"`
	DeviceName             string                      `json:"deviceName"`
	ContainerInterfaceName string                      `json:"containerInterfaceName"`
	DPDKSocketPath         string                      `json:"dpdkSocketPath"`
	HWAddr                 string                      `json:"hwaddr"`
	KernelPath             string                      `json:"kernelPath"`
	PCIAddr                string                      `json:"pciAddr"`

	// in Direct Attachable Network, this maybe set by CNI plugins
	Interfaces  []*current.Interface `json:"interfaces,omitempty"`
	IPs         []*current.IPConfig  `json:"ips,omitempty"`
	Routes      []*types.Route       `json:"routes,omitempty"`
	DNS         types.DNS            `json:"dns,omitempty"`
	Annotations map[string]string    `json:"annotations,omitempty"`
}

func FromResult(networkType DirectAttachableNetworkType, device, containerInfName string, r *current.Result) *DirectAttachableNetwork {
	dan := &DirectAttachableNetwork{
		NetworkType:            networkType,
		DeviceName:             device,
		ContainerInterfaceName: containerInfName,
		Interfaces:             r.Interfaces,
		IPs:                    r.IPs,
		Routes:                 r.Routes,
		DNS:                    r.DNS,
		Annotations:            r.Annotations,
	}

	return dan
}

func Log(format string, args ...interface{}) {
	f, err := os.OpenFile("/tmp/cni.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	fmt.Fprintf(f, format+"\n", args...)
}

func (dan *DirectAttachableNetwork) Save(metaFile string) error {
	body, _ := json.MarshalIndent(dan, "", " ")
	path := filepath.Dir(metaFile)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return os.WriteFile(metaFile, body, 0644)
}
