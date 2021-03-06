// Copyright (c) 2017-2018 Intel Corporation
// Copyright (c) 2018 Huawei Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package drivers

import (
	"encoding/hex"

	"github.com/kata-containers/runtime/virtcontainers/device/api"
	"github.com/kata-containers/runtime/virtcontainers/device/config"
	"github.com/kata-containers/runtime/virtcontainers/utils"
)

// VhostUserNetDevice is a network vhost-user based device
type VhostUserNetDevice struct {
	config.VhostUserDeviceAttrs
	DeviceInfo *config.DeviceInfo
}

//
// VhostUserNetDevice's implementation of the device interface:
//

// Attach is standard interface of api.Device, it's used to add device to some
// DeviceReceiver
func (device *VhostUserNetDevice) Attach(devReceiver api.DeviceReceiver) (err error) {
	if device.DeviceInfo.Hotplugged {
		return nil
	}

	// generate a unique ID to be used for hypervisor commandline fields
	randBytes, err := utils.GenerateRandomBytes(8)
	if err != nil {
		return err
	}
	id := hex.EncodeToString(randBytes)

	device.ID = id
	device.Type = device.DeviceType()

	defer func() {
		if err == nil {
			device.DeviceInfo.Hotplugged = true
		}
	}()
	return devReceiver.AppendDevice(device)
}

// Detach is standard interface of api.Device, it's used to remove device from some
// DeviceReceiver
func (device *VhostUserNetDevice) Detach(devReceiver api.DeviceReceiver) error {
	if !device.DeviceInfo.Hotplugged {
		return nil
	}

	device.DeviceInfo.Hotplugged = false
	return nil
}

// IsAttached checks if the device is attached
func (device *VhostUserNetDevice) IsAttached() bool {
	return device.DeviceInfo.Hotplugged
}

// DeviceID returns device ID
func (device *VhostUserNetDevice) DeviceID() string {
	return device.ID
}

// DeviceType is standard interface of api.Device, it returns device type
func (device *VhostUserNetDevice) DeviceType() config.DeviceType {
	return config.VhostUserNet
}

// GetDeviceInfo returns device information used for creating
func (device *VhostUserNetDevice) GetDeviceInfo() interface{} {
	device.Type = device.DeviceType()
	return &device.VhostUserDeviceAttrs
}
