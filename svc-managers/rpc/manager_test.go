//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil), nil
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil), nil
}

func mockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetManagerByURL(url string) (string, *errors.Error) {
	managerData := make(map[string]interface{})
	managerData["ManagerType"] = "BMC"
	managerData["Status"] = `{"State":"Enabled"}}`
	managerData["Name"] = "somePlugin"
	switch url {
	case "/redfish/v1/Managers/nonExistingUUID", "/redfish/v1/Managers/invalidURL.1", "/redfish/v1/Managers/invalidURL", "/redfish/v1/Managers/invalidID":
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	case "/redfish/v1/Managers/noDevice":
		managerData["ManagerType"] = "Service"
		managerData["Status"] = `{"State":"Absent"}}`
		managerData["Name"] = "odimra"
		managerData["ID"] = "noDevice"
		managerData["UUID"] = "noDevice"
		managerData["FirmwareVersion"] = "1.0"
	case "/redfish/v1/Managers/uuid":
		managerData["Name"] = "someOtherID"
	case "/redfish/v1/Managers/noPlugin":
		managerData["Name"] = "noPlugin"
	case "/redfish/v1/Managers/noToken":
		managerData["Name"] = "noToken"
	case "/redfish/v1/Managers/" + config.Data.RootServiceUUID:
		managerData["ManagerType"] = "Service"
		managerData["Status"] = `{"State":"Enabled"}}`
		managerData["Name"] = "odimra"
		managerData["ManagerID"] = config.Data.RootServiceUUID
		managerData["UUID"] = config.Data.RootServiceUUID
		managerData["FirmwareVersion"] = "1.0"
	}
	data, _ := json.Marshal(managerData)
	return string(data), nil
}

func mockGetPluginData(pluginID string) (mgrmodel.Plugin, *errors.Error) {
	if pluginID == "someOtherID" {
		return mgrmodel.Plugin{
			IP:                "localhost",
			Port:              "9091",
			Username:          "admin",
			Password:          []byte("password"),
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
		}, nil
	} else if pluginID == "noToken" {
		return mgrmodel.Plugin{
			IP:                "localhost",
			Port:              "9092",
			Username:          "admin",
			Password:          []byte("password"),
			ID:                "noToken",
			PreferredAuthType: "XAuthToken",
		}, nil
	} else if pluginID == "noPlugin" {
		return mgrmodel.Plugin{}, errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return mgrmodel.Plugin{
		IP:                "localhost",
		Port:              "9093",
		Username:          "admin",
		Password:          []byte("password"),
		ID:                "somePlugin",
		PreferredAuthType: "BasicAuth",
	}, nil
}

func mockUpdateData(key string, updateData map[string]interface{}, table string) error {
	if key == "/redfish/v1/Managers/uuid.1/VirtualMedia/1" {
		return nil
	} else if key == "/redfish/v1/Managers/uuid1.1/VirtualMedia/4" {
		return errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return nil
}

func mockGetResource(table, key string) (string, *errors.Error) {
	if key == "/redfish/v1/Managers/uuid1.1/Ethernet" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	} else if key == "/redfish/v1/Managers/uuid1.1/Virtual" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	} else if key == "/redfish/v1/Managers/uuid1.1/VirtualMedia/4" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return "body", nil
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

func mockGetDeviceInfo(ctx context.Context, req mgrcommon.ResourceInfoRequest) (string, error) {
	if req.URL == "/redfish/v1/Managers/deviceAbsent.1" || req.URL == "/redfish/v1/Managers/uuid1.1/Ethernet" {
		return "", fmt.Errorf("error")
	} else if req.URL == "/redfish/v1/Managers/uuid1.1/Virtual" {
		return "", fmt.Errorf("error")
	} else if req.URL == "/redfish/v1/Managers/uuid1.1/VirtualMedia/4" {
		return "", fmt.Errorf("error")
	}
	manager := mgrmodel.Manager{
		Status: &mgrmodel.Status{
			State: "Enabled",
		},
	}
	dataByte, err := json.Marshal(manager)
	return string(dataByte), err
}

func mockDeviceRequest(ctx context.Context, req mgrcommon.ResourceInfoRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
	if req.URL == "/redfish/v1/Managers/deviceAbsent.1" || req.URL == "/redfish/v1/Managers/uuid1.1/Virtual" {
		resp.StatusCode = http.StatusNotFound
		resp.StatusMessage = response.ResourceNotFound
		return resp
	}
	manager := mgrmodel.Manager{
		Status: &mgrmodel.Status{
			State: "Enabled",
		},
	}
	dataByte, err := json.Marshal(manager)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(dataByte, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	return resp
}

func mockGetExternalInterface() *managers.ExternalInterface {
	return &managers.ExternalInterface{
		Device: managers.Device{
			GetDeviceInfo: mockGetDeviceInfo,
			ContactClient: mockContactClient,
			DeviceRequest: mockDeviceRequest,
		},
		DB: managers.DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetManagerByURL:     mockGetManagerByURL,
			GetPluginData:       mockGetPluginData,
			UpdateData:          mockUpdateData,
			GetResource:         mockGetResource,
		},
	}
}

func mockGetAllKeysFromTable(table string) ([]string, error) {
	return []string{"/redfish/v1/Managers/uuid.1"}, nil
}

func TestGetManagerCollection(t *testing.T) {
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	ctx := mockContext()
	type args struct {
		req  *managersproto.ManagerRequest
		resp *managersproto.ManagerResponse
	}
	tests := []struct {
		name       string
		mgr        *Managers
		args       args
		StatusCode int32
	}{
		{
			name: "Request with valid token",
			mgr:  mgr,
			args: args{
				req: &managersproto.ManagerRequest{
					SessionToken: "validToken",
				},
				resp: &managersproto.ManagerResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			mgr:  mgr,
			args: args{
				req: &managersproto.ManagerRequest{
					SessionToken: "invalidToken",
				},
				resp: &managersproto.ManagerResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.mgr.GetManagersCollection(ctx, tt.args.req)
			if err != nil {
				t.Errorf("Manager.GetManagersCollection() got = %v, want %v", err, nil)
			}
			if resp.StatusCode != tt.StatusCode {
				t.Errorf("Manager.GetManagersCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetManagerwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		SessionToken: "InvalidToken",
	}
	resp, _ := mgr.GetManager(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusOK.")
}
func TestGetManagerwithValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    config.Data.RootServiceUUID,
		SessionToken: "validToken",
	}
	resp, err := mgr.GetManager(ctx, req)
	assert.Nil(t, err, "There should be no error")

	var manager mgrmodel.Manager
	json.Unmarshal(resp.Body, &manager)

	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
	assert.Equal(t, "odimra", manager.Name, "incorrect name")
	assert.Equal(t, "Service", manager.ManagerType, "incorrect type")
	assert.Equal(t, req.ManagerID, manager.ID, "incorrect id")
	assert.Equal(t, "1.0", manager.FirmwareVersion, "incorrect firmware version")
}

func TestGetManagerResourcewithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
	}
	resp, _ := mgr.GetManagersResource(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}
func TestGetManagerResourcewithValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/EthernetInterfaces/1",
		ResourceID:   "1",
	}
	resp, err := mgr.GetManagersResource(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")
}
func TestVirtualMediaEject(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia",
		ResourceID:   "1",
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.VirtualMediaEject(ctx, req)
	fmt.Println(resp)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "1",
		URL:          "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.VirtualMediaEject(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")

}
func TestVirtualMediaInsert(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		ResourceID:   "1",
		RequestBody: []byte(`{"Image":"http://10.1.0.1/ISO/ubuntu-18.04.4-server-amd64.iso",
										"Inserted":true,
										"WriteProtected":true
										}`),
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.VirtualMediaInsert(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "1",
		URL:          "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody:  []byte(`{"Image":"http://10.1.0.1/ISO/ubuntu-18.04.4-server-amd64.iso"}`),
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.VirtualMediaInsert(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestGetRemoteAccountService(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1",
		ResourceID:   "1",
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.GetRemoteAccountService(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "1",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1",
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.GetRemoteAccountService(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestCreateRemoteAccountService(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		ResourceID:   "1",
		RequestBody: []byte(`{"UserName":"Username",
										"Password":"Password",
										"RoleId":"Administrator"
										}`),
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.CreateRemoteAccountService(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusCreated, "Status code should be StatusCreated.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "1",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{"UserName":"Username",
										"Password":"Password",
										"RoleId":"Administrator"}`),
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.CreateRemoteAccountService(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestDeleteRemoteAccountService(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/5",
		ResourceID:   "5",
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.DeleteRemoteAccountService(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusNoContent, "Status code should be StatusNoContent.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "5",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/5",
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.DeleteRemoteAccountService(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestUpdateRemoteAccountService(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1",
		ResourceID:   "1",
		RequestBody: []byte(`{
								"Password":"Password",
								"RoleId":"Administrator"
							}`),
	}
	var resp = &managersproto.ManagerResponse{}
	resp, err := mgr.UpdateRemoteAccountService(ctx, req)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	// Invalid
	req = &managersproto.ManagerRequest{
		ManagerID:    "uuid.1",
		SessionToken: "InvalidToken",
		ResourceID:   "1",
		URL:          "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1",
		RequestBody: []byte(`{
										"Password":"Password",
										"RoleId":"Administrator"}`),
	}
	resp = &managersproto.ManagerResponse{}
	resp, _ = mgr.UpdateRemoteAccountService(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}
