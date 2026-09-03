//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

// TestE2EAutomationDashboards exercises every dashboard method against a live
// Directus instance: single and bulk create, single and bulk read, single and
// keyed and batch patch, and single and bulk delete.
func TestE2EAutomationDashboards(t *testing.T) {
	c := itestClient(t)

	// CreateDashboard
	created, err := c.CreateDashboard(&Dashboard{Name: uniqueName("e2e_dash"), Icon: "space_dashboard"}, nil)
	if err != nil {
		t.Fatalf("CreateDashboard: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateDashboard: expected non-empty id")
	}
	t.Cleanup(func() { _ = c.DeleteDashboard(created.ID) })

	// GetDashboard
	got, err := c.GetDashboard(created.ID, nil)
	if err != nil {
		t.Fatalf("GetDashboard: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetDashboard: expected id %q, got %+v", created.ID, got)
	}

	// GetDashboards
	list, err := c.GetDashboards(nil)
	if err != nil {
		t.Fatalf("GetDashboards: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("GetDashboards: expected non-empty list")
	}

	// PatchDashboard
	patched, err := c.PatchDashboard(created.ID, &Dashboard{Note: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchDashboard: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchDashboard: expected id %q, got %+v", created.ID, patched)
	}

	// CreateDashboards
	batch, err := c.CreateDashboards([]Dashboard{
		{Name: uniqueName("e2e_dash"), Icon: "space_dashboard"},
		{Name: uniqueName("e2e_dash"), Icon: "space_dashboard"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateDashboards: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreateDashboards: expected 2, got %d", len(batch))
	}
	keys := []string{batch[0].ID, batch[1].ID}
	t.Cleanup(func() { _ = c.DeleteDashboards(keys) })

	// PatchDashboards (keys + shared payload)
	updated, err := c.PatchDashboards(keys, &Dashboard{Note: "e2e_keys"}, nil)
	if err != nil {
		t.Fatalf("PatchDashboards: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchDashboards: expected 2, got %d", len(updated))
	}

	// PatchDashboardsBatch (array of items)
	batched, err := c.PatchDashboardsBatch([]Dashboard{
		{ID: batch[0].ID, Note: "e2e_batch0"},
		{ID: batch[1].ID, Note: "e2e_batch1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchDashboardsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchDashboardsBatch: expected 2, got %d", len(batched))
	}

	// DeleteDashboards
	if err := c.DeleteDashboards(keys); err != nil {
		t.Fatalf("DeleteDashboards: %v", err)
	}

	// DeleteDashboard
	if err := c.DeleteDashboard(created.ID); err != nil {
		t.Fatalf("DeleteDashboard: %v", err)
	}
}

// TestE2EAutomationPanels exercises every panel method against a live Directus
// instance. Panels require a parent dashboard.
func TestE2EAutomationPanels(t *testing.T) {
	c := itestClient(t)
	dashID := newTestDashboard(t, c)

	newPanel := func() *Panel {
		return &Panel{
			Dashboard:  dashID,
			Type:       "label",
			Name:       uniqueName("e2e_panel"),
			Width:      4,
			Height:     4,
			PositionX:  1,
			PositionY:  1,
			ShowHeader: true,
		}
	}

	// CreatePanel
	created, err := c.CreatePanel(newPanel(), nil)
	if err != nil {
		t.Fatalf("CreatePanel: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreatePanel: expected non-empty id")
	}
	t.Cleanup(func() { _ = c.DeletePanel(created.ID) })

	// GetPanel
	got, err := c.GetPanel(created.ID, nil)
	if err != nil {
		t.Fatalf("GetPanel: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetPanel: expected id %q, got %+v", created.ID, got)
	}

	// GetPanels
	list, err := c.GetPanels(nil)
	if err != nil {
		t.Fatalf("GetPanels: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("GetPanels: expected non-empty list")
	}

	// PatchPanel
	patched, err := c.PatchPanel(created.ID, &Panel{Note: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchPanel: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchPanel: expected id %q, got %+v", created.ID, patched)
	}

	// CreatePanels
	batch, err := c.CreatePanels([]Panel{*newPanel(), *newPanel()}, nil)
	if err != nil {
		t.Fatalf("CreatePanels: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreatePanels: expected 2, got %d", len(batch))
	}
	keys := []string{batch[0].ID, batch[1].ID}
	t.Cleanup(func() { _ = c.DeletePanels(keys) })

	// PatchPanels (keys + shared payload)
	updated, err := c.PatchPanels(keys, &Panel{Note: "e2e_keys"}, nil)
	if err != nil {
		t.Fatalf("PatchPanels: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchPanels: expected 2, got %d", len(updated))
	}

	// PatchPanelsBatch (array of items)
	batched, err := c.PatchPanelsBatch([]Panel{
		{ID: batch[0].ID, Note: "e2e_batch0"},
		{ID: batch[1].ID, Note: "e2e_batch1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchPanelsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchPanelsBatch: expected 2, got %d", len(batched))
	}

	// DeletePanels
	if err := c.DeletePanels(keys); err != nil {
		t.Fatalf("DeletePanels: %v", err)
	}

	// DeletePanel
	if err := c.DeletePanel(created.ID); err != nil {
		t.Fatalf("DeletePanel: %v", err)
	}
}

// TestE2EAutomationFlows exercises every flow method against a live Directus
// instance, including triggering a webhook flow.
func TestE2EAutomationFlows(t *testing.T) {
	c := itestClient(t)

	newFlow := func() *Flow {
		return &Flow{
			Name:    uniqueName("e2e_flow"),
			Status:  "active",
			Trigger: "webhook",
			Options: map[string]any{"method": "POST", "async": false},
		}
	}

	// CreateFlow
	created, err := c.CreateFlow(newFlow(), nil)
	if err != nil {
		t.Fatalf("CreateFlow: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateFlow: expected non-empty id")
	}
	t.Cleanup(func() { _ = c.DeleteFlow(created.ID) })

	// GetFlow
	got, err := c.GetFlow(created.ID, nil)
	if err != nil {
		t.Fatalf("GetFlow: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetFlow: expected id %q, got %+v", created.ID, got)
	}

	// GetFlows
	list, err := c.GetFlows(nil)
	if err != nil {
		t.Fatalf("GetFlows: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("GetFlows: expected non-empty list")
	}

	// PatchFlow
	patched, err := c.PatchFlow(created.ID, &Flow{Description: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchFlow: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchFlow: expected id %q, got %+v", created.ID, patched)
	}

	// TriggerFlow (POST body)
	if _, err := c.TriggerFlow("POST", created.ID, map[string]string{"x": "1"}); err != nil {
		t.Fatalf("TriggerFlow: %v", err)
	}

	// CreateFlows
	batch, err := c.CreateFlows([]Flow{*newFlow(), *newFlow()}, nil)
	if err != nil {
		t.Fatalf("CreateFlows: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreateFlows: expected 2, got %d", len(batch))
	}
	keys := []string{batch[0].ID, batch[1].ID}
	t.Cleanup(func() { _ = c.DeleteFlows(keys) })

	// PatchFlows (keys + shared payload)
	updated, err := c.PatchFlows(keys, &Flow{Description: "e2e_keys"}, nil)
	if err != nil {
		t.Fatalf("PatchFlows: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchFlows: expected 2, got %d", len(updated))
	}

	// PatchFlowsBatch (array of items)
	batched, err := c.PatchFlowsBatch([]Flow{
		{ID: batch[0].ID, Description: "e2e_batch0"},
		{ID: batch[1].ID, Description: "e2e_batch1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchFlowsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchFlowsBatch: expected 2, got %d", len(batched))
	}

	// DeleteFlows
	if err := c.DeleteFlows(keys); err != nil {
		t.Fatalf("DeleteFlows: %v", err)
	}

	// DeleteFlow
	if err := c.DeleteFlow(created.ID); err != nil {
		t.Fatalf("DeleteFlow: %v", err)
	}
}

// TestE2EAutomationOperations exercises every operation method against a live
// Directus instance. Operations require a parent flow.
func TestE2EAutomationOperations(t *testing.T) {
	c := itestClient(t)

	flow, err := c.CreateFlow(&Flow{
		Name:    uniqueName("e2e_flow"),
		Status:  "active",
		Trigger: "webhook",
		Options: map[string]any{"method": "POST", "async": false},
	}, nil)
	if err != nil {
		t.Fatalf("CreateFlow (fixture): %v", err)
	}
	// Flow deleted last, after all operations are removed.
	t.Cleanup(func() { _ = c.DeleteFlow(flow.ID) })

	newOp := func() *Operation {
		return &Operation{
			Flow:      flow.ID,
			Key:       "e2e_log",
			Type:      "log",
			Name:      uniqueName("e2e_op"),
			PositionX: 20,
			PositionY: 20,
		}
	}

	// CreateOperation
	created, err := c.CreateOperation(newOp(), nil)
	if err != nil {
		t.Fatalf("CreateOperation: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateOperation: expected non-empty id")
	}
	t.Cleanup(func() { _ = c.DeleteOperation(created.ID) })

	// GetOperation
	got, err := c.GetOperation(created.ID, nil)
	if err != nil {
		t.Fatalf("GetOperation: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetOperation: expected id %q, got %+v", created.ID, got)
	}

	// GetOperations
	list, err := c.GetOperations(nil)
	if err != nil {
		t.Fatalf("GetOperations: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("GetOperations: expected non-empty list")
	}

	// PatchOperation
	patched, err := c.PatchOperation(created.ID, &Operation{Name: uniqueName("e2e_op")}, nil)
	if err != nil {
		t.Fatalf("PatchOperation: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchOperation: expected id %q, got %+v", created.ID, patched)
	}

	// CreateOperations
	batch, err := c.CreateOperations([]Operation{*newOp(), *newOp()}, nil)
	if err != nil {
		t.Fatalf("CreateOperations: %v", err)
	}
	if len(batch) != 2 {
		t.Fatalf("CreateOperations: expected 2, got %d", len(batch))
	}
	keys := []string{batch[0].ID, batch[1].ID}
	t.Cleanup(func() { _ = c.DeleteOperations(keys) })

	// PatchOperations (keys + shared payload)
	updated, err := c.PatchOperations(keys, &Operation{Name: uniqueName("e2e_op")}, nil)
	if err != nil {
		t.Fatalf("PatchOperations: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("PatchOperations: expected 2, got %d", len(updated))
	}

	// PatchOperationsBatch (array of items)
	batched, err := c.PatchOperationsBatch([]Operation{
		{ID: batch[0].ID, Name: uniqueName("e2e_op")},
		{ID: batch[1].ID, Name: uniqueName("e2e_op")},
	}, nil)
	if err != nil {
		t.Fatalf("PatchOperationsBatch: %v", err)
	}
	if len(batched) != 2 {
		t.Fatalf("PatchOperationsBatch: expected 2, got %d", len(batched))
	}

	// DeleteOperations
	if err := c.DeleteOperations(keys); err != nil {
		t.Fatalf("DeleteOperations: %v", err)
	}

	// DeleteOperation
	if err := c.DeleteOperation(created.ID); err != nil {
		t.Fatalf("DeleteOperation: %v", err)
	}
}
