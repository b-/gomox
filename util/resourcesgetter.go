package util

import (
	"context"
	"fmt"

	"github.com/luthermonson/go-proxmox"
)

const ( // resource getter filters
	VmRequestFilter      = "vm"
	StorageRequestFilter = "storage"
	NodeRequestFilter    = "node"
	SdnRequestFilter     = "sdn"
)

//goland:noinspection GoDeprecation
const ( // returned resource types to further filter
	NodeResource    = "node"
	StorageResource = "storage"
	PoolResource    = "pool"
	QemuResource    = "qemu"
	LxcResource     = "lxc"
	OpenVzResource  = "openvz" // deprecated
	SdnResource     = "sdn"
)

func GetVirtualMachineByVMID(ctx context.Context, vmid uint64, client proxmox.Client) (
	vm *proxmox.VirtualMachine,
	err error,
) {
	var node *proxmox.Node

	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, VmRequestFilter)
	if err != nil {
		return nil, err
	}

	for _, rs := range resources {
		if rs.VMID == vmid {
			node, err = client.Node(ctx, rs.Node)
			if err != nil {
				return nil, err
			}
			vm, err = node.VirtualMachine(ctx, int(rs.VMID))
			if err != nil {
				return nil, err
			}
		}
	}

	if vm == nil {
		err = fmt.Errorf("no vm with id found: %d", vmid)
	}

	return vm, err
}

// getResourceListConfig defines the options for
type getResourceListConfig struct {
	filter        string
	furtherFilter []string
}

func (c getResourceListConfig) furtherFilterFunction() []proxmox.ClusterResource {

}

// GetResourceListOption specifies the type of Resources for GetResource to get.
type GetResourceListOption func(c *getResourceListConfig)

type ResourceFilterFunction func() []proxmox.ClusterResource

func FilterRsByType(t string, resources []proxmox.ClusterResource) ResourceFilterFunction {
	return func() []proxmox.ClusterResource {
		var rsList []proxmox.ClusterResource
		var outList []proxmox.ClusterResource
		for _, rs := range resources {
			rsList = append(rsList, rs)
			if rs.Type == t {
				outList = append(outList, rs)
			}
		}
		return outList
	}
}

// WithServerFilter makes GetResourceList apply the given filter constant to the query.
// WithServerFilter(VmServerFilter) or WithServerFilter(StorageServerFilter)
func WithServerFilter(filter string) GetResourceListOption {
	return func(c *getResourceListConfig) { c.filter = filter }
}

// WithFurtherFilter makes GetResourceList further filter the returned Resources
// in order to (for example) return only Resources of the QemuResource type.
// WithFurtherFilter(QemuResource)
func WithFurtherFilter(furtherFilter string) GetResourceListOption {
	return func(c *getResourceListConfig) { c.furtherFilter = append(c.furtherFilter, furtherFilter) }
}
func WithFurtherFilterFunc(function ResourceFilterFunction) GetResourceListOption {
	return func(c *getResourceListConfig) { c.furtherFilterFunction = function }
}

func GetResourceList(
	ctx context.Context,
	client proxmox.Client,
	opts ...GetResourceListOption,
) (
	outList []*proxmox.ClusterResource,
	err error,
) {
	var rsList []*proxmox.ClusterResource

	c := &getResourceListConfig{
		filter:        "",
		furtherFilter: nil,
	}
	for _, opt := range opts {
		opt(c)
	}
	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx, c.filter)
	if err != nil {
		return nil, err
	}
	for _, rs := range resources {
		rsList = append(rsList, rs)
		if c.furtherFilter != nil {
			for _, f := range c.furtherFilter {
				if rs.Type == f {
					outList = append(outList, rs)
				}
			}
		} else {
			outList = append(outList, rs)
		}
	}

	return outList, nil
}

func GetVirtualMachineList(ctx context.Context, client proxmox.Client) (vmList []*proxmox.VirtualMachine, err error) {
	var node *proxmox.Node
	var vm *proxmox.VirtualMachine

	resources, err := GetResourceList(ctx, client, WithServerFilter(VmRequestFilter))
	var rsList []*proxmox.ClusterResource
	for _, rs := range resources {
		node, err = client.Node(ctx, rs.Node)
		if err != nil {
			return nil, err
		}
		vm, err = node.VirtualMachine(ctx, int(rs.VMID))
		rsList = append(rsList, rs)
		if rs.Type == "qemu" {
			vmList = append(vmList, vm)
		}
	}

	return vmList, nil
}
