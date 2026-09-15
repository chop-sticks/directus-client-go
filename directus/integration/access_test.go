//go:build integration

package integration

import (
	. "github.com/chop-sticks/directus-client-go/directus"
	"testing"
)

func TestE2EAccessUsers(t *testing.T) {
	c := itestClient(t)

	// GetUsersMe
	me, err := c.GetUsersMe(nil)
	if err != nil {
		t.Fatalf("GetUsersMe: %v", err)
	}
	if me == nil || me.ID == "" {
		t.Fatalf("GetUsersMe returned empty user")
	}

	// GetUsers
	users, err := c.GetUsers(nil)
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) == 0 {
		t.Fatalf("GetUsers returned no users")
	}

	// CreateUser
	created, err := c.CreateUser(&User{
		Email:     uniqueName("e2e") + "@example.com",
		Password:  "Testpassw0rd!",
		FirstName: "E2E",
	}, nil)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateUser returned empty user")
	}
	t.Cleanup(func() { _ = c.DeleteUser(created.ID) })

	// GetUser
	got, err := c.GetUser(created.ID, nil)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetUser mismatch: %+v", got)
	}

	// PatchUser
	patched, err := c.PatchUser(created.ID, &User{Title: "e2e-title"}, nil)
	if err != nil {
		t.Fatalf("PatchUser: %v", err)
	}
	if patched == nil || patched.Title != "e2e-title" {
		t.Fatalf("PatchUser did not persist title: %+v", patched)
	}

	// CreateUsers (2)
	multi, err := c.CreateUsers([]User{
		{Email: uniqueName("e2e") + "@example.com", Password: "Testpassw0rd!"},
		{Email: uniqueName("e2e") + "@example.com", Password: "Testpassw0rd!"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateUsers: %v", err)
	}
	if len(multi) != 2 {
		t.Fatalf("CreateUsers expected 2, got %d", len(multi))
	}
	keys := []string{multi[0].ID, multi[1].ID}
	t.Cleanup(func() { _ = c.DeleteUsers(keys) })

	// PatchUsers (keys)
	pk, err := c.PatchUsers(keys, &User{Title: "e2e-bulk"}, nil)
	if err != nil {
		t.Fatalf("PatchUsers: %v", err)
	}
	if len(pk) != 2 {
		t.Fatalf("PatchUsers expected 2, got %d", len(pk))
	}

	// PatchUsersBatch
	pb, err := c.PatchUsersBatch([]User{
		{ID: multi[0].ID, Title: "e2e-batch-0"},
		{ID: multi[1].ID, Title: "e2e-batch-1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchUsersBatch: %v", err)
	}
	if len(pb) != 2 {
		t.Fatalf("PatchUsersBatch expected 2, got %d", len(pb))
	}

	// PatchUsersMe
	updatedMe, err := c.PatchUsersMe(&User{Description: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchUsersMe: %v", err)
	}
	if updatedMe == nil {
		t.Fatalf("PatchUsersMe returned nil")
	}

	// DeleteUser (single) — create one dedicated user
	single, err := c.CreateUser(&User{
		Email:    uniqueName("e2e") + "@example.com",
		Password: "Testpassw0rd!",
	}, nil)
	if err != nil {
		t.Fatalf("CreateUser (for delete): %v", err)
	}
	if err := c.DeleteUser(single.ID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	// DeleteUsers (many) — create two dedicated users
	del, err := c.CreateUsers([]User{
		{Email: uniqueName("e2e") + "@example.com", Password: "Testpassw0rd!"},
		{Email: uniqueName("e2e") + "@example.com", Password: "Testpassw0rd!"},
	}, nil)
	if err != nil {
		t.Fatalf("CreateUsers (for delete): %v", err)
	}
	if err := c.DeleteUsers([]string{del[0].ID, del[1].ID}); err != nil {
		t.Fatalf("DeleteUsers: %v", err)
	}
}

func TestE2EAccessUsersInviteAndRegister(t *testing.T) {
	c := itestClient(t)
	roleID := newTestRole(t, c)

	// InviteUser — environment-dependent (email transport), log-only.
	if err := c.InviteUser(uniqueName("e2e")+"@example.com", roleID, ""); err != nil {
		t.Logf("InviteUser: %v", err)
	}

	// AcceptUserInvite — bad token must fail.
	if err := c.AcceptUserInvite("bad-token", "Testpassw0rd!"); err == nil {
		t.Fatalf("AcceptUserInvite(bad) expected error, got nil")
	}

	// RegisterUser — public registration disabled; expect error.
	if err := c.RegisterUser(uniqueName("e2e")+"@example.com", "Testpassw0rd!", nil); err == nil {
		t.Fatalf("RegisterUser expected error, got nil")
	}

	// RegisterUserVerify — bad token must fail.
	if err := c.RegisterUserVerify("bad-token"); err == nil {
		t.Fatalf("RegisterUserVerify(bad) expected error, got nil")
	}
}

func TestE2EAccessUsersTwoFactor(t *testing.T) {
	c := itestClient(t)

	// GenerateTwoFactorSecret — expect success, returns map with secret.
	secretMap, err := c.GenerateTwoFactorSecret(defaultPassword)
	if err != nil {
		t.Fatalf("GenerateTwoFactorSecret: %v", err)
	}
	secret, _ := secretMap["secret"].(string)
	if secret == "" {
		t.Fatalf("GenerateTwoFactorSecret returned no secret: %+v", secretMap)
	}

	// EnableTwoFactor — wrong otp must fail.
	if err := c.EnableTwoFactor(secret, "000000"); err == nil {
		t.Fatalf("EnableTwoFactor(bad otp) expected error, got nil")
	}

	// DisableTwoFactor — wrong otp must fail.
	if err := c.DisableTwoFactor("000000"); err == nil {
		t.Fatalf("DisableTwoFactor(bad otp) expected error, got nil")
	}
}

func TestE2EAccessRoles(t *testing.T) {
	c := itestClient(t)

	// GetRoles
	roles, err := c.GetRoles(nil)
	if err != nil {
		t.Fatalf("GetRoles: %v", err)
	}
	if len(roles) == 0 {
		t.Fatalf("GetRoles returned no roles")
	}

	// GetRolesMe
	if _, err := c.GetRolesMe(nil); err != nil {
		t.Fatalf("GetRolesMe: %v", err)
	}

	// CreateRole
	created, err := c.CreateRole(&Role{Name: uniqueName("e2e_role")}, nil)
	if err != nil {
		t.Fatalf("CreateRole: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateRole returned empty role")
	}
	t.Cleanup(func() { _ = c.DeleteRole(created.ID) })

	// GetRole
	got, err := c.GetRole(created.ID, nil)
	if err != nil {
		t.Fatalf("GetRole: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetRole mismatch: %+v", got)
	}

	// PatchRole
	patched, err := c.PatchRole(created.ID, &Role{Description: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchRole: %v", err)
	}
	if patched == nil || patched.Description != "e2e" {
		t.Fatalf("PatchRole did not persist: %+v", patched)
	}

	// CreateRoles (2)
	multi, err := c.CreateRoles([]Role{
		{Name: uniqueName("e2e_role")},
		{Name: uniqueName("e2e_role")},
	}, nil)
	if err != nil {
		t.Fatalf("CreateRoles: %v", err)
	}
	if len(multi) != 2 {
		t.Fatalf("CreateRoles expected 2, got %d", len(multi))
	}
	keys := []string{multi[0].ID, multi[1].ID}
	t.Cleanup(func() { _ = c.DeleteRoles(keys) })

	// PatchRoles (keys)
	pk, err := c.PatchRoles(keys, &Role{Description: "e2e-bulk"}, nil)
	if err != nil {
		t.Fatalf("PatchRoles: %v", err)
	}
	if len(pk) != 2 {
		t.Fatalf("PatchRoles expected 2, got %d", len(pk))
	}

	// PatchRolesBatch
	pb, err := c.PatchRolesBatch([]Role{
		{ID: multi[0].ID, Description: "e2e-batch-0"},
		{ID: multi[1].ID, Description: "e2e-batch-1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchRolesBatch: %v", err)
	}
	if len(pb) != 2 {
		t.Fatalf("PatchRolesBatch expected 2, got %d", len(pb))
	}

	// DeleteRole (single)
	single, err := c.CreateRole(&Role{Name: uniqueName("e2e_role")}, nil)
	if err != nil {
		t.Fatalf("CreateRole (for delete): %v", err)
	}
	if err := c.DeleteRole(single.ID); err != nil {
		t.Fatalf("DeleteRole: %v", err)
	}

	// DeleteRoles (many)
	del, err := c.CreateRoles([]Role{
		{Name: uniqueName("e2e_role")},
		{Name: uniqueName("e2e_role")},
	}, nil)
	if err != nil {
		t.Fatalf("CreateRoles (for delete): %v", err)
	}
	if err := c.DeleteRoles([]string{del[0].ID, del[1].ID}); err != nil {
		t.Fatalf("DeleteRoles: %v", err)
	}
}

func TestE2EAccessPolicies(t *testing.T) {
	c := itestClient(t)

	// GetPolicies
	policies, err := c.GetPolicies(nil)
	if err != nil {
		t.Fatalf("GetPolicies: %v", err)
	}
	if len(policies) == 0 {
		t.Fatalf("GetPolicies returned no policies")
	}

	// GetPolicyGlobals
	globals, err := c.GetPolicyGlobals()
	if err != nil {
		t.Fatalf("GetPolicyGlobals: %v", err)
	}
	if globals == nil {
		t.Fatalf("GetPolicyGlobals returned nil")
	}

	// CreatePolicy
	created, err := c.CreatePolicy(&Policy{Name: uniqueName("e2e_policy")}, nil)
	if err != nil {
		t.Fatalf("CreatePolicy: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreatePolicy returned empty policy")
	}
	t.Cleanup(func() { _ = c.DeletePolicy(created.ID) })

	// GetPolicy
	got, err := c.GetPolicy(created.ID, nil)
	if err != nil {
		t.Fatalf("GetPolicy: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetPolicy mismatch: %+v", got)
	}

	// PatchPolicy
	patched, err := c.PatchPolicy(created.ID, &Policy{Description: "e2e"}, nil)
	if err != nil {
		t.Fatalf("PatchPolicy: %v", err)
	}
	if patched == nil || patched.Description != "e2e" {
		t.Fatalf("PatchPolicy did not persist: %+v", patched)
	}

	// CreatePolicies (2)
	multi, err := c.CreatePolicies([]Policy{
		{Name: uniqueName("e2e_policy")},
		{Name: uniqueName("e2e_policy")},
	}, nil)
	if err != nil {
		t.Fatalf("CreatePolicies: %v", err)
	}
	if len(multi) != 2 {
		t.Fatalf("CreatePolicies expected 2, got %d", len(multi))
	}
	keys := []string{multi[0].ID, multi[1].ID}
	t.Cleanup(func() { _ = c.DeletePolicies(keys) })

	// PatchPolicies (keys)
	pk, err := c.PatchPolicies(keys, &Policy{Description: "e2e-bulk"}, nil)
	if err != nil {
		t.Fatalf("PatchPolicies: %v", err)
	}
	if len(pk) != 2 {
		t.Fatalf("PatchPolicies expected 2, got %d", len(pk))
	}

	// PatchPoliciesBatch
	pb, err := c.PatchPoliciesBatch([]Policy{
		{ID: multi[0].ID, Description: "e2e-batch-0"},
		{ID: multi[1].ID, Description: "e2e-batch-1"},
	}, nil)
	if err != nil {
		t.Fatalf("PatchPoliciesBatch: %v", err)
	}
	if len(pb) != 2 {
		t.Fatalf("PatchPoliciesBatch expected 2, got %d", len(pb))
	}

	// DeletePolicy (single)
	single, err := c.CreatePolicy(&Policy{Name: uniqueName("e2e_policy")}, nil)
	if err != nil {
		t.Fatalf("CreatePolicy (for delete): %v", err)
	}
	if err := c.DeletePolicy(single.ID); err != nil {
		t.Fatalf("DeletePolicy: %v", err)
	}

	// DeletePolicies (many)
	del, err := c.CreatePolicies([]Policy{
		{Name: uniqueName("e2e_policy")},
		{Name: uniqueName("e2e_policy")},
	}, nil)
	if err != nil {
		t.Fatalf("CreatePolicies (for delete): %v", err)
	}
	if err := c.DeletePolicies([]string{del[0].ID, del[1].ID}); err != nil {
		t.Fatalf("DeletePolicies: %v", err)
	}
}

func TestE2EAccessPermissions(t *testing.T) {
	c := itestClient(t)
	policyID := newTestPolicy(t, c)
	coll := newTestCollection(t, c)

	// GetPermissions
	perms, err := c.GetPermissions(nil)
	if err != nil {
		t.Fatalf("GetPermissions: %v", err)
	}
	_ = perms

	// GetUserPermissions
	up, err := c.GetUserPermissions()
	if err != nil {
		t.Fatalf("GetUserPermissions: %v", err)
	}
	if up == nil {
		t.Fatalf("GetUserPermissions returned nil")
	}

	// GetItemPermissions
	ip, err := c.GetItemPermissions(coll, "")
	if err != nil {
		t.Fatalf("GetItemPermissions: %v", err)
	}
	if ip == nil {
		t.Fatalf("GetItemPermissions returned nil")
	}

	// CreatePermission
	created, err := c.CreatePermission(&Permission{
		Policy:     policyID,
		Collection: coll,
		Action:     "read",
		Fields:     []string{"*"},
	}, nil)
	if err != nil {
		t.Fatalf("CreatePermission: %v", err)
	}
	if created == nil || created.ID == 0 {
		t.Fatalf("CreatePermission returned empty permission")
	}
	t.Cleanup(func() { _ = c.DeletePermission(created.ID) })

	// GetPermission
	got, err := c.GetPermission(created.ID, nil)
	if err != nil {
		t.Fatalf("GetPermission: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetPermission mismatch: %+v", got)
	}

	// PatchPermission — custom permission rules are a licensed feature; on an
	// unlicensed instance this is restricted, so treat an error as acceptable.
	patched, err := c.PatchPermission(created.ID, &Permission{Action: "read", Fields: []string{"id"}}, nil)
	if err != nil {
		t.Logf("PatchPermission (edition-restricted): %v", err)
	} else if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchPermission mismatch: %+v", patched)
	}

	// CreatePermissions (batch). NOTE: Directus returns ALL permissions on a
	// batch create (not just the created rows), so derive the created ids by
	// querying this policy instead of trusting the response cardinality.
	if _, err := c.CreatePermissions([]Permission{
		{Policy: policyID, Collection: coll, Action: "create", Fields: []string{"*"}},
		{Policy: policyID, Collection: coll, Action: "update", Fields: []string{"*"}},
	}, nil); err != nil {
		t.Fatalf("CreatePermissions: %v", err)
	}
	list, err := c.GetPermissions(&Query{Filter: map[string]any{
		"policy": map[string]any{"_eq": policyID},
		"action": map[string]any{"_in": []string{"create", "update"}},
	}})
	if err != nil {
		t.Fatalf("GetPermissions(filter): %v", err)
	}
	var keys []int
	for _, p := range list {
		keys = append(keys, p.ID)
	}
	if len(keys) < 2 {
		t.Fatalf("expected >=2 permissions for policy, got %d", len(keys))
	}
	t.Cleanup(func() { _ = c.DeletePermissions(keys) })

	// PatchPermissions (keys) — edition-restricted; log-only.
	if _, err := c.PatchPermissions(keys, &Permission{Fields: []string{"id"}}, nil); err != nil {
		t.Logf("PatchPermissions (edition-restricted): %v", err)
	}

	// PatchPermissionsBatch — edition-restricted; log-only.
	if _, err := c.PatchPermissionsBatch([]Permission{
		{ID: keys[0], Fields: []string{"title"}},
		{ID: keys[1], Fields: []string{"title"}},
	}, nil); err != nil {
		t.Logf("PatchPermissionsBatch (edition-restricted): %v", err)
	}

	// DeletePermission (single). Single creates reliably return one row.
	single, err := c.CreatePermission(&Permission{Policy: policyID, Collection: coll, Action: "share", Fields: []string{"*"}}, nil)
	if err != nil {
		t.Fatalf("CreatePermission (for delete): %v", err)
	}
	if err := c.DeletePermission(single.ID); err != nil {
		t.Fatalf("DeletePermission: %v", err)
	}

	// DeletePermissions (many). Use a second collection so (policy,collection,
	// action) stays unique, and single creates for reliable ids.
	coll2 := newTestCollection(t, c)
	d1, err := c.CreatePermission(&Permission{Policy: policyID, Collection: coll2, Action: "read", Fields: []string{"*"}}, nil)
	if err != nil {
		t.Fatalf("CreatePermission d1: %v", err)
	}
	d2, err := c.CreatePermission(&Permission{Policy: policyID, Collection: coll2, Action: "create", Fields: []string{"*"}}, nil)
	if err != nil {
		t.Fatalf("CreatePermission d2: %v", err)
	}
	if err := c.DeletePermissions([]int{d1.ID, d2.ID}); err != nil {
		t.Fatalf("DeletePermissions: %v", err)
	}
}

func TestE2EAccessRecords(t *testing.T) {
	c := itestClient(t)
	policyID := newTestPolicy(t, c)
	userID := newTestUser(t, c)
	roleID := newTestRole(t, c)

	// GetAccesses
	accesses, err := c.GetAccesses(nil)
	if err != nil {
		t.Fatalf("GetAccesses: %v", err)
	}
	if len(accesses) == 0 {
		t.Fatalf("GetAccesses returned no records")
	}

	// CreateAccess — bind policy to a user.
	created, err := c.CreateAccess(&Access{Policy: policyID, User: userID}, nil)
	if err != nil {
		t.Fatalf("CreateAccess: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateAccess returned empty record")
	}
	t.Cleanup(func() { _ = c.DeleteAccess(created.ID) })

	// GetAccess
	got, err := c.GetAccess(created.ID, nil)
	if err != nil {
		t.Fatalf("GetAccess: %v", err)
	}
	if got == nil || got.ID != created.ID {
		t.Fatalf("GetAccess mismatch: %+v", got)
	}

	// PatchAccess — move the binding to a role instead.
	patched, err := c.PatchAccess(created.ID, &Access{User: nil, Role: roleID}, nil)
	if err != nil {
		t.Fatalf("PatchAccess: %v", err)
	}
	if patched == nil || patched.ID != created.ID {
		t.Fatalf("PatchAccess mismatch: %+v", patched)
	}

	// CreateAccesses (2) — two distinct policies bound to the same role.
	p2 := newTestPolicy(t, c)
	p3 := newTestPolicy(t, c)
	multi, err := c.CreateAccesses([]Access{
		{Policy: p2, Role: roleID},
		{Policy: p3, Role: roleID},
	}, nil)
	if err != nil {
		t.Fatalf("CreateAccesses: %v", err)
	}
	if len(multi) != 2 {
		t.Fatalf("CreateAccesses expected 2, got %d", len(multi))
	}
	keys := []string{multi[0].ID, multi[1].ID}
	t.Cleanup(func() { _ = c.DeleteAccesses(keys) })

	// PatchAccesses (keys)
	pk, err := c.PatchAccesses(keys, &Access{Sort: new(1)}, nil)
	if err != nil {
		t.Fatalf("PatchAccesses: %v", err)
	}
	if len(pk) != 2 {
		t.Fatalf("PatchAccesses expected 2, got %d", len(pk))
	}

	// PatchAccessesBatch
	pb, err := c.PatchAccessesBatch([]Access{
		{ID: multi[0].ID, Sort: new(2)},
		{ID: multi[1].ID, Sort: new(3)},
	}, nil)
	if err != nil {
		t.Fatalf("PatchAccessesBatch: %v", err)
	}
	if len(pb) != 2 {
		t.Fatalf("PatchAccessesBatch expected 2, got %d", len(pb))
	}

	// DeleteAccess (single)
	single, err := c.CreateAccess(&Access{Policy: newTestPolicy(t, c), Role: roleID}, nil)
	if err != nil {
		t.Fatalf("CreateAccess (for delete): %v", err)
	}
	if err := c.DeleteAccess(single.ID); err != nil {
		t.Fatalf("DeleteAccess: %v", err)
	}

	// DeleteAccesses (many)
	del, err := c.CreateAccesses([]Access{
		{Policy: newTestPolicy(t, c), Role: roleID},
		{Policy: newTestPolicy(t, c), Role: roleID},
	}, nil)
	if err != nil {
		t.Fatalf("CreateAccesses (for delete): %v", err)
	}
	if err := c.DeleteAccesses([]string{del[0].ID, del[1].ID}); err != nil {
		t.Fatalf("DeleteAccesses: %v", err)
	}
}
