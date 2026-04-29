package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm"
	vmsnapshot "github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_snapshot"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// SnapshotOps 定义 VM 快照操作。
type SnapshotOps interface {
	ListSnapshots(ctx context.Context, vmID string) ([]Snapshot, error)
	CreateSnapshot(ctx context.Context, vmID, name string) (TaskRef, error)
	RevertSnapshot(ctx context.Context, vmID, snapID string) (TaskRef, error)
	DeleteSnapshot(ctx context.Context, snapID string) (TaskRef, error)
}

// ---------- ListSnapshots ----------

func (c *defaultClient) ListSnapshots(ctx context.Context, vmID string) ([]Snapshot, error) {
	params := vmsnapshot.NewGetVMSnapshotsParams()
	params.SetContext(ctx)
	body := &models.GetVMSnapshotsRequestBody{}
	if vmID != "" {
		body.Where = &models.VMSnapshotWhereInput{
			VM: &models.VMWhereInput{ID: pointy.String(vmID)},
		}
	}
	params.SetRequestBody(body)

	resp, err := c.api.VMSnapshot.GetVMSnapshots(params)
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	out := make([]Snapshot, 0, len(resp.Payload))
	for _, s := range resp.Payload {
		out = append(out, toSnapshot(s))
	}
	return out, nil
}

// ---------- CreateSnapshot ----------

func (c *defaultClient) CreateSnapshot(ctx context.Context, vmID, name string) (TaskRef, error) {
	params := vmsnapshot.NewCreateVMSnapshotParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMSnapshotCreationParams{
		Data: []*models.VMSnapshotCreationParamsData{
			{
				VMID: pointy.String(vmID),
				Name: pointy.String(name),
			},
		},
	})
	resp, err := c.api.VMSnapshot.CreateVMSnapshot(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create snapshot: %w", err)
	}
	return firstSnapshotTaskRef(resp.Payload), nil
}

// ---------- RevertSnapshot ----------

func (c *defaultClient) RevertSnapshot(ctx context.Context, vmID, snapID string) (TaskRef, error) {
	params := vm.NewRollbackVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMRollbackParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMRollbackParamsData{
			SnapshotID: pointy.String(snapID),
		},
	})
	resp, err := c.api.VM.RollbackVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("revert snapshot %s: %w", snapID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- DeleteSnapshot ----------

func (c *defaultClient) DeleteSnapshot(ctx context.Context, snapID string) (TaskRef, error) {
	params := vmsnapshot.NewDeleteVMSnapshotParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMSnapshotDeletionParams{
		Where: &models.VMSnapshotWhereInput{ID: pointy.String(snapID)},
	})
	resp, err := c.api.VMSnapshot.DeleteVMSnapshot(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("delete snapshot %s: %w", snapID, err)
	}
	return firstDeleteSnapshotTaskRef(resp.Payload), nil
}

// ---------- helpers ----------

func toSnapshot(s *models.VMSnapshot) Snapshot {
	out := Snapshot{}
	if s.ID != nil {
		out.ID = *s.ID
	}
	if s.Name != nil {
		out.Name = *s.Name
	}
	if s.Description != nil {
		out.Description = *s.Description
	}
	if s.LocalCreatedAt != nil {
		out.CreatedAt = *s.LocalCreatedAt
	}
	if s.VM != nil && s.VM.ID != nil {
		out.VMID = *s.VM.ID
	}
	return out
}

func firstSnapshotTaskRef(items []*models.WithTaskVMSnapshot) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	ref := TaskRef{EntityKind: "VMSnapshot"}
	it := items[0]
	if it.TaskID != nil {
		ref.ID = *it.TaskID
	}
	if it.Data != nil && it.Data.ID != nil {
		ref.EntityID = *it.Data.ID
	}
	return ref
}

func firstDeleteSnapshotTaskRef(items []*models.WithTaskDeleteVMSnapshot) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	ref := TaskRef{EntityKind: "VMSnapshot"}
	it := items[0]
	if it.TaskID != nil {
		ref.ID = *it.TaskID
	}
	if it.Data != nil && it.Data.ID != nil {
		ref.EntityID = *it.Data.ID
	}
	return ref
}
