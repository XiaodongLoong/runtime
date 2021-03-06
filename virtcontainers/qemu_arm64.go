// Copyright (c) 2018 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"context"
	"time"

	govmmQemu "github.com/intel/govmm/qemu"
)

type qemuArm64 struct {
	// inherit from qemuArchBase, overwrite methods if needed
	qemuArchBase
}

const defaultQemuPath = "/usr/bin/qemu-system-aarch64"

const defaultQemuMachineType = QemuVirt

const qmpMigrationWaitTimeout = 10 * time.Second

const defaultQemuMachineOptions = "usb=off,accel=kvm,gic-version=host"

var defaultGICVersion = uint32(3)

var qemuPaths = map[string]string{
	QemuVirt: defaultQemuPath,
}

var kernelParams = []Param{
	{"console", "hvc0"},
	{"console", "hvc1"},
	{"iommu.passthrough", "0"},
}

var supportedQemuMachines = []govmmQemu.Machine{
	{
		Type:    QemuVirt,
		Options: defaultQemuMachineOptions,
	},
}

//In qemu, maximum number of vCPUs depends on the GIC version, or on how
//many redistributors we can fit into the memory map.
//related codes are under github.com/qemu/qemu/hw/arm/virt.c(Line 135 and 1306 in stable-2.11)
//for now, qemu only supports v2 and v3, we treat v4 as v3 based on
//backward compatibility.
var gicList = map[uint32]uint32{
	uint32(2): uint32(8),
	uint32(3): uint32(123),
	uint32(4): uint32(123),
}

// MaxQemuVCPUs returns the maximum number of vCPUs supported
func MaxQemuVCPUs() uint32 {
	return gicList[defaultGICVersion]
}

func newQemuArch(config HypervisorConfig) qemuArch {
	machineType := config.HypervisorMachineType
	if machineType == "" {
		machineType = defaultQemuMachineType
	}

	q := &qemuArm64{
		qemuArchBase{
			machineType:           machineType,
			memoryOffset:          config.MemOffset,
			qemuPaths:             qemuPaths,
			supportedQemuMachines: supportedQemuMachines,
			kernelParamsNonDebug:  kernelParamsNonDebug,
			kernelParamsDebug:     kernelParamsDebug,
			kernelParams:          kernelParams,
			disableNvdimm:         config.DisableImageNvdimm,
			dax:                   true,
		},
	}

	q.handleImagePath(config)

	return q
}

func (q *qemuArm64) bridges(number uint32) {
	q.Bridges = genericBridges(number, q.machineType)
}

// appendBridges appends to devices the given bridges
func (q *qemuArm64) appendBridges(devices []govmmQemu.Device) []govmmQemu.Device {
	return genericAppendBridges(devices, q.Bridges, q.machineType)
}

func (q *qemuArm64) appendImage(devices []govmmQemu.Device, path string) ([]govmmQemu.Device, error) {
	if !q.disableNvdimm {
		return q.appendNvdimmImage(devices, path)
	}
	return q.appendBlockImage(devices, path)
}

func (q *qemuArm64) setIgnoreSharedMemoryMigrationCaps(_ context.Context, _ *govmmQemu.QMP) error {
	// x-ignore-shared not support in arm64 for now
	return nil
}
