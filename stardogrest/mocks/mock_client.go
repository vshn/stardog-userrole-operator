// Code generated by MockGen. DO NOT EDIT.
// Source: stardogrest/client/stardog_client_test.go

// Package stardogmock is a generated GoMock package.
package stardogmock

import (
	reflect "reflect"

	runtime "github.com/go-openapi/runtime"
	gomock "github.com/golang/mock/gomock"
	db "github.com/vshn/stardog-userrole-operator/stardogrest/client/db"
	roles "github.com/vshn/stardog-userrole-operator/stardogrest/client/roles"
	roles_permissions "github.com/vshn/stardog-userrole-operator/stardogrest/client/roles_permissions"
	users "github.com/vshn/stardog-userrole-operator/stardogrest/client/users"
	users_permissions "github.com/vshn/stardog-userrole-operator/stardogrest/client/users_permissions"
	users_roles "github.com/vshn/stardog-userrole-operator/stardogrest/client/users_roles"
)

// MockStardogTestClient is a mock of StardogTestClient interface.
type MockStardogTestClient struct {
	ctrl     *gomock.Controller
	recorder *MockStardogTestClientMockRecorder
}

// MockStardogTestClientMockRecorder is the mock recorder for MockStardogTestClient.
type MockStardogTestClientMockRecorder struct {
	mock *MockStardogTestClient
}

// NewMockStardogTestClient creates a new mock instance.
func NewMockStardogTestClient(ctrl *gomock.Controller) *MockStardogTestClient {
	mock := &MockStardogTestClient{ctrl: ctrl}
	mock.recorder = &MockStardogTestClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStardogTestClient) EXPECT() *MockStardogTestClientMockRecorder {
	return m.recorder
}

// AddRole mocks base method.
func (m *MockStardogTestClient) AddRole(params *users_roles.AddRoleParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*users_roles.AddRoleNoContent, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddRole", varargs...)
	ret0, _ := ret[0].(*users_roles.AddRoleNoContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRole indicates an expected call of AddRole.
func (mr *MockStardogTestClientMockRecorder) AddRole(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRole", reflect.TypeOf((*MockStardogTestClient)(nil).AddRole), varargs...)
}

// AddRolePermission mocks base method.
func (m *MockStardogTestClient) AddRolePermission(params *roles_permissions.AddRolePermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles_permissions.ClientOption) (*roles_permissions.AddRolePermissionCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddRolePermission", varargs...)
	ret0, _ := ret[0].(*roles_permissions.AddRolePermissionCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRolePermission indicates an expected call of AddRolePermission.
func (mr *MockStardogTestClientMockRecorder) AddRolePermission(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRolePermission", reflect.TypeOf((*MockStardogTestClient)(nil).AddRolePermission), varargs...)
}

// AddUserPermission mocks base method.
func (m *MockStardogTestClient) AddUserPermission(params *users_permissions.AddUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_permissions.ClientOption) (*users_permissions.AddUserPermissionCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddUserPermission", varargs...)
	ret0, _ := ret[0].(*users_permissions.AddUserPermissionCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUserPermission indicates an expected call of AddUserPermission.
func (mr *MockStardogTestClientMockRecorder) AddUserPermission(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserPermission", reflect.TypeOf((*MockStardogTestClient)(nil).AddUserPermission), varargs...)
}

// ChangePassword mocks base method.
func (m *MockStardogTestClient) ChangePassword(params *users.ChangePasswordParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.ChangePasswordOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangePassword", varargs...)
	ret0, _ := ret[0].(*users.ChangePasswordOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockStardogTestClientMockRecorder) ChangePassword(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockStardogTestClient)(nil).ChangePassword), varargs...)
}

// CreateNewDatabase mocks base method.
func (m *MockStardogTestClient) CreateNewDatabase(params *db.CreateNewDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...db.ClientOption) (*db.CreateNewDatabaseCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateNewDatabase", varargs...)
	ret0, _ := ret[0].(*db.CreateNewDatabaseCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewDatabase indicates an expected call of CreateNewDatabase.
func (mr *MockStardogTestClientMockRecorder) CreateNewDatabase(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewDatabase", reflect.TypeOf((*MockStardogTestClient)(nil).CreateNewDatabase), varargs...)
}

// CreateRole mocks base method.
func (m *MockStardogTestClient) CreateRole(params *roles.CreateRoleParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles.ClientOption) (*roles.CreateRoleCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateRole", varargs...)
	ret0, _ := ret[0].(*roles.CreateRoleCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockStardogTestClientMockRecorder) CreateRole(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockStardogTestClient)(nil).CreateRole), varargs...)
}

// CreateUser mocks base method.
func (m *MockStardogTestClient) CreateUser(params *users.CreateUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.CreateUserCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateUser", varargs...)
	ret0, _ := ret[0].(*users.CreateUserCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStardogTestClientMockRecorder) CreateUser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStardogTestClient)(nil).CreateUser), varargs...)
}

// DropDatabase mocks base method.
func (m *MockStardogTestClient) DropDatabase(params *db.DropDatabaseParams, authInfo runtime.ClientAuthInfoWriter, opts ...db.ClientOption) (*db.DropDatabaseOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DropDatabase", varargs...)
	ret0, _ := ret[0].(*db.DropDatabaseOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DropDatabase indicates an expected call of DropDatabase.
func (mr *MockStardogTestClientMockRecorder) DropDatabase(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropDatabase", reflect.TypeOf((*MockStardogTestClient)(nil).DropDatabase), varargs...)
}

// GetDBSize mocks base method.
func (m *MockStardogTestClient) GetDBSize(params *db.GetDBSizeParams, authInfo runtime.ClientAuthInfoWriter, opts ...db.ClientOption) (*db.GetDBSizeOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDBSize", varargs...)
	ret0, _ := ret[0].(*db.GetDBSizeOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDBSize indicates an expected call of GetDBSize.
func (mr *MockStardogTestClientMockRecorder) GetDBSize(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDBSize", reflect.TypeOf((*MockStardogTestClient)(nil).GetDBSize), varargs...)
}

// GetUser mocks base method.
func (m *MockStardogTestClient) GetUser(params *users.GetUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.GetUserOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUser", varargs...)
	ret0, _ := ret[0].(*users.GetUserOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStardogTestClientMockRecorder) GetUser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStardogTestClient)(nil).GetUser), varargs...)
}

// IsEnabled mocks base method.
func (m *MockStardogTestClient) IsEnabled(params *users.IsEnabledParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.IsEnabledOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsEnabled", varargs...)
	ret0, _ := ret[0].(*users.IsEnabledOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEnabled indicates an expected call of IsEnabled.
func (mr *MockStardogTestClientMockRecorder) IsEnabled(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEnabled", reflect.TypeOf((*MockStardogTestClient)(nil).IsEnabled), varargs...)
}

// IsSuperuser mocks base method.
func (m *MockStardogTestClient) IsSuperuser(params *users.IsSuperuserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.IsSuperuserOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsSuperuser", varargs...)
	ret0, _ := ret[0].(*users.IsSuperuserOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSuperuser indicates an expected call of IsSuperuser.
func (mr *MockStardogTestClientMockRecorder) IsSuperuser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSuperuser", reflect.TypeOf((*MockStardogTestClient)(nil).IsSuperuser), varargs...)
}

// ListDatabases mocks base method.
func (m *MockStardogTestClient) ListDatabases(params *db.ListDatabasesParams, authInfo runtime.ClientAuthInfoWriter, opts ...db.ClientOption) (*db.ListDatabasesOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListDatabases", varargs...)
	ret0, _ := ret[0].(*db.ListDatabasesOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDatabases indicates an expected call of ListDatabases.
func (mr *MockStardogTestClientMockRecorder) ListDatabases(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDatabases", reflect.TypeOf((*MockStardogTestClient)(nil).ListDatabases), varargs...)
}

// ListEffectivePermissions mocks base method.
func (m *MockStardogTestClient) ListEffectivePermissions(params *users_permissions.ListEffectivePermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_permissions.ClientOption) (*users_permissions.ListEffectivePermissionsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListEffectivePermissions", varargs...)
	ret0, _ := ret[0].(*users_permissions.ListEffectivePermissionsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEffectivePermissions indicates an expected call of ListEffectivePermissions.
func (mr *MockStardogTestClientMockRecorder) ListEffectivePermissions(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEffectivePermissions", reflect.TypeOf((*MockStardogTestClient)(nil).ListEffectivePermissions), varargs...)
}

// ListRolePermissions mocks base method.
func (m *MockStardogTestClient) ListRolePermissions(params *roles_permissions.ListRolePermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles_permissions.ClientOption) (*roles_permissions.ListRolePermissionsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListRolePermissions", varargs...)
	ret0, _ := ret[0].(*roles_permissions.ListRolePermissionsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRolePermissions indicates an expected call of ListRolePermissions.
func (mr *MockStardogTestClientMockRecorder) ListRolePermissions(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRolePermissions", reflect.TypeOf((*MockStardogTestClient)(nil).ListRolePermissions), varargs...)
}

// ListRoleUsers mocks base method.
func (m *MockStardogTestClient) ListRoleUsers(params *users_roles.ListRoleUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*users_roles.ListRoleUsersOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListRoleUsers", varargs...)
	ret0, _ := ret[0].(*users_roles.ListRoleUsersOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoleUsers indicates an expected call of ListRoleUsers.
func (mr *MockStardogTestClientMockRecorder) ListRoleUsers(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoleUsers", reflect.TypeOf((*MockStardogTestClient)(nil).ListRoleUsers), varargs...)
}

// ListRoles mocks base method.
func (m *MockStardogTestClient) ListRoles(params *roles.ListRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles.ClientOption) (*roles.ListRolesOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListRoles", varargs...)
	ret0, _ := ret[0].(*roles.ListRolesOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoles indicates an expected call of ListRoles.
func (mr *MockStardogTestClientMockRecorder) ListRoles(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoles", reflect.TypeOf((*MockStardogTestClient)(nil).ListRoles), varargs...)
}

// ListUserPermissions mocks base method.
func (m *MockStardogTestClient) ListUserPermissions(params *users_permissions.ListUserPermissionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_permissions.ClientOption) (*users_permissions.ListUserPermissionsOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUserPermissions", varargs...)
	ret0, _ := ret[0].(*users_permissions.ListUserPermissionsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserPermissions indicates an expected call of ListUserPermissions.
func (mr *MockStardogTestClientMockRecorder) ListUserPermissions(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserPermissions", reflect.TypeOf((*MockStardogTestClient)(nil).ListUserPermissions), varargs...)
}

// ListUserRoles mocks base method.
func (m *MockStardogTestClient) ListUserRoles(params *users_roles.ListUserRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*users_roles.ListUserRolesOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUserRoles", varargs...)
	ret0, _ := ret[0].(*users_roles.ListUserRolesOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserRoles indicates an expected call of ListUserRoles.
func (mr *MockStardogTestClientMockRecorder) ListUserRoles(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserRoles", reflect.TypeOf((*MockStardogTestClient)(nil).ListUserRoles), varargs...)
}

// ListUsers mocks base method.
func (m *MockStardogTestClient) ListUsers(params *users.ListUsersParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.ListUsersOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUsers", varargs...)
	ret0, _ := ret[0].(*users.ListUsersOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStardogTestClientMockRecorder) ListUsers(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStardogTestClient)(nil).ListUsers), varargs...)
}

// PutRoles mocks base method.
func (m *MockStardogTestClient) PutRoles(params *users_roles.PutRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*users_roles.PutRolesOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutRoles", varargs...)
	ret0, _ := ret[0].(*users_roles.PutRolesOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutRoles indicates an expected call of PutRoles.
func (mr *MockStardogTestClientMockRecorder) PutRoles(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutRoles", reflect.TypeOf((*MockStardogTestClient)(nil).PutRoles), varargs...)
}

// RemoveRole mocks base method.
func (m *MockStardogTestClient) RemoveRole(params *roles.RemoveRoleParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles.ClientOption) (*roles.RemoveRoleNoContent, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveRole", varargs...)
	ret0, _ := ret[0].(*roles.RemoveRoleNoContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRole indicates an expected call of RemoveRole.
func (mr *MockStardogTestClientMockRecorder) RemoveRole(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRole", reflect.TypeOf((*MockStardogTestClient)(nil).RemoveRole), varargs...)
}

// RemoveRoleOfUser mocks base method.
func (m *MockStardogTestClient) RemoveRoleOfUser(params *users_roles.RemoveRoleOfUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_roles.ClientOption) (*users_roles.RemoveRoleOfUserNoContent, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveRoleOfUser", varargs...)
	ret0, _ := ret[0].(*users_roles.RemoveRoleOfUserNoContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRoleOfUser indicates an expected call of RemoveRoleOfUser.
func (mr *MockStardogTestClientMockRecorder) RemoveRoleOfUser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRoleOfUser", reflect.TypeOf((*MockStardogTestClient)(nil).RemoveRoleOfUser), varargs...)
}

// RemoveRolePermission mocks base method.
func (m *MockStardogTestClient) RemoveRolePermission(params *roles_permissions.RemoveRolePermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...roles_permissions.ClientOption) (*roles_permissions.RemoveRolePermissionCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveRolePermission", varargs...)
	ret0, _ := ret[0].(*roles_permissions.RemoveRolePermissionCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRolePermission indicates an expected call of RemoveRolePermission.
func (mr *MockStardogTestClientMockRecorder) RemoveRolePermission(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRolePermission", reflect.TypeOf((*MockStardogTestClient)(nil).RemoveRolePermission), varargs...)
}

// RemoveUser mocks base method.
func (m *MockStardogTestClient) RemoveUser(params *users.RemoveUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.RemoveUserNoContent, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveUser", varargs...)
	ret0, _ := ret[0].(*users.RemoveUserNoContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockStardogTestClientMockRecorder) RemoveUser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockStardogTestClient)(nil).RemoveUser), varargs...)
}

// RemoveUserPermission mocks base method.
func (m *MockStardogTestClient) RemoveUserPermission(params *users_permissions.RemoveUserPermissionParams, authInfo runtime.ClientAuthInfoWriter, opts ...users_permissions.ClientOption) (*users_permissions.RemoveUserPermissionCreated, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveUserPermission", varargs...)
	ret0, _ := ret[0].(*users_permissions.RemoveUserPermissionCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveUserPermission indicates an expected call of RemoveUserPermission.
func (mr *MockStardogTestClientMockRecorder) RemoveUserPermission(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserPermission", reflect.TypeOf((*MockStardogTestClient)(nil).RemoveUserPermission), varargs...)
}

// SetEnabled mocks base method.
func (m *MockStardogTestClient) SetEnabled(params *users.SetEnabledParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.SetEnabledOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetEnabled", varargs...)
	ret0, _ := ret[0].(*users.SetEnabledOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetEnabled indicates an expected call of SetEnabled.
func (mr *MockStardogTestClientMockRecorder) SetEnabled(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEnabled", reflect.TypeOf((*MockStardogTestClient)(nil).SetEnabled), varargs...)
}

// SetTransport mocks base method.
func (m *MockStardogTestClient) SetTransport(transport runtime.ClientTransport) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTransport", transport)
}

// SetTransport indicates an expected call of SetTransport.
func (mr *MockStardogTestClientMockRecorder) SetTransport(transport interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTransport", reflect.TypeOf((*MockStardogTestClient)(nil).SetTransport), transport)
}

// Submit mocks base method.
func (m *MockStardogTestClient) Submit(arg0 *runtime.ClientOperation) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Submit", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Submit indicates an expected call of Submit.
func (mr *MockStardogTestClientMockRecorder) Submit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Submit", reflect.TypeOf((*MockStardogTestClient)(nil).Submit), arg0)
}

// ValidateUser mocks base method.
func (m *MockStardogTestClient) ValidateUser(params *users.ValidateUserParams, authInfo runtime.ClientAuthInfoWriter, opts ...users.ClientOption) (*users.ValidateUserOK, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{params, authInfo}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateUser", varargs...)
	ret0, _ := ret[0].(*users.ValidateUserOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockStardogTestClientMockRecorder) ValidateUser(params, authInfo interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params, authInfo}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockStardogTestClient)(nil).ValidateUser), varargs...)
}
