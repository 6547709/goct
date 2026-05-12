package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/usb_device"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// UsbDeviceOps 定义 USB 设备操作。
type UsbDeviceOps interface {
	ListUsbDevices(ctx context.Context, opts ListOpts) ([]UsbDevice, error)
	MountUsbDevice(ctx context.Context, usbID string, vmID string) (TaskRef, error)
	UnmountUsbDevice(ctx context.Context, usbID string) (TaskRef, error)
}

// ---------- ListUsbDevices ----------

func (c *defaultClient) ListUsbDevices(ctx context.Context, opts ListOpts) ([]UsbDevice, error) {
	params := usb_device.NewGetUsbDevicesParams()
	params.SetContext(ctx)

	where := &models.UsbDeviceWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Host = &models.HostWhereInput{Cluster: &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}}
		hasWhere = true
	}

	body := &models.GetUsbDevicesRequestBody{}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.UsbDevice.GetUsbDevices(params)
	if err != nil {
		return nil, fmt.Errorf("list usb devices: %w", err)
	}
	out := make([]UsbDevice, 0, len(resp.Payload))
	for _, d := range resp.Payload {
		out = append(out, toUsbDevice(d))
	}
	return out, nil
}

// ---------- MountUsbDevice ----------

func (c *defaultClient) MountUsbDevice(ctx context.Context, usbID string, vmID string) (TaskRef, error) {
	params := usb_device.NewMountUsbDeviceParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.UsbDeviceMountParams{
		Data: &models.UsbDeviceMountParamsData{
			VMID: pointy.String(vmID),
		},
		Where: &models.UsbDeviceWhereInput{
			ID: pointy.String(usbID),
		},
	})

	resp, err := c.api.UsbDevice.MountUsbDevice(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("mount usb device: %w", err)
	}
	return toWithTaskUsbDeviceTaskRef(resp.Payload), nil
}

// ---------- UnmountUsbDevice ----------

func (c *defaultClient) UnmountUsbDevice(ctx context.Context, usbID string) (TaskRef, error) {
	params := usb_device.NewUnmountUsbDeviceParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.UsbDeviceUnmountParams{
		Where: &models.UsbDeviceWhereInput{
			ID: pointy.String(usbID),
		},
	})

	resp, err := c.api.UsbDevice.UnmountUsbDevice(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("unmount usb device: %w", err)
	}
	return toWithTaskUsbDeviceTaskRef(resp.Payload), nil
}

// toUsbDevice 把 SDK models.UsbDevice 转成内部 UsbDevice 模型。
func toUsbDevice(d *models.UsbDevice) UsbDevice {
	out := UsbDevice{}
	if d.ID != nil {
		out.ID = *d.ID
	}
	if d.Name != nil {
		out.Name = *d.Name
	}
	if d.Description != nil {
		out.Description = *d.Description
	}
	if d.LocalID != nil {
		out.LocalID = *d.LocalID
	}
	if d.Manufacturer != nil {
		out.Manufacturer = *d.Manufacturer
	}
	if d.Status != nil {
		out.Status = string(*d.Status)
	}
	if d.UsbType != nil {
		out.UsbType = *d.UsbType
	}
	if d.Size != nil {
		out.Size = *d.Size
	}
	if d.Binded != nil {
		out.Binded = *d.Binded
	}
	if d.Host != nil {
		if d.Host.ID != nil {
			out.HostID = *d.Host.ID
		}
		if d.Host.Name != nil {
			out.HostName = *d.Host.Name
		}
	}
	if d.VM != nil {
		if d.VM.ID != nil {
			out.VMID = *d.VM.ID
		}
		if d.VM.Name != nil {
			out.VMName = *d.VM.Name
		}
	}
	if d.LocalCreatedAt != nil {
		out.LocalCreatedAt = *d.LocalCreatedAt
	}
	return out
}

// toWithTaskUsbDeviceTaskRef 从 WithTaskUsbDevice 提取 TaskRef。
func toWithTaskUsbDeviceTaskRef(items []*models.WithTaskUsbDevice) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	task := items[0]
	if task.TaskID == nil {
		return TaskRef{}
	}
	ref := TaskRef{ID: *task.TaskID}
	if task.Data != nil && task.Data.ID != nil {
		ref.EntityID = *task.Data.ID
	}
	ref.EntityKind = "usb_device"
	return ref
}
