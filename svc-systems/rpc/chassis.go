//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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

// Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/chassis"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
)

var podName = os.Getenv("POD_NAME")

var (
	// JSONMarshalFunc function pointer for the json.Marshal
	JSONMarshalFunc = json.Marshal
)

// NewChassisRPC returns an instance of ChassisRPC
func NewChassisRPC(
	authWrapper func(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error),
	createHandler *chassis.Create,
	getCollectionHandler *chassis.GetCollection,
	deleteHandler *chassis.Delete,
	getHandler *chassis.Get,
	updateHandler *chassis.Update) *ChassisRPC {

	return &ChassisRPC{
		IsAuthorizedRPC:      authWrapper,
		GetCollectionHandler: getCollectionHandler,
		GetHandler:           getHandler,
		DeleteHandler:        deleteHandler,
		UpdateHandler:        updateHandler,
		CreateHandler:        createHandler,
	}
}

// ChassisRPC struct helps to register service
type ChassisRPC struct {
	IsAuthorizedRPC      func(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error)
	GetCollectionHandler *chassis.GetCollection
	GetHandler           *chassis.Get
	DeleteHandler        *chassis.Delete
	UpdateHandler        *chassis.Update
	CreateHandler        *chassis.Create
}

// UpdateChassis defines the operations which handles the RPC request response
// for updating the system resource of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (cha *ChassisRPC) UpdateChassis(ctx context.Context, req *chassisproto.UpdateChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming chassis update request with %s", req.URL)
	var resp chassisproto.GetChassisResponse
	r := auth(ctx, cha.IsAuthorizedRPC, req.SessionToken, []string{common.PrivilegeConfigureComponents}, func() response.RPC {
		return cha.UpdateHandler.Handle(ctx, req)
	})

	rewrite(ctx, r, &resp)
	l.LogWithFields(ctx).Debugf("outgoing response from update chassis request %s", string(resp.Body))
	return &resp, nil
}

// DeleteChassis defines the operations which handles the RPC request response
// for deleting the system resource of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (cha *ChassisRPC) DeleteChassis(ctx context.Context, req *chassisproto.DeleteChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming chassis Delete request with %s", req.URL)
	var resp chassisproto.GetChassisResponse
	r := auth(ctx, cha.IsAuthorizedRPC, req.SessionToken, []string{common.PrivilegeConfigureComponents}, func() response.RPC {
		return cha.DeleteHandler.Handle(ctx, req)
	})

	rewrite(ctx, r, &resp)
	l.LogWithFields(ctx).Debugf("outgoing response Delete Chassis request: %s", string(resp.Body))
	return &resp, nil
}

// CreateChassis defines the operations which handles the RPC request response
// for creating the system resource of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (cha *ChassisRPC) CreateChassis(ctx context.Context, req *chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugln("incoming chassis create request")
	var resp chassisproto.GetChassisResponse
	r := auth(ctx, cha.IsAuthorizedRPC, req.SessionToken, []string{common.PrivilegeConfigureComponents}, func() response.RPC {
		return cha.CreateHandler.Handle(ctx, req)
	})

	rewrite(ctx, r, &resp)
	l.LogWithFields(ctx).Debugf("outgoing response for Create Chassis: %s", string(resp.Body))
	return &resp, nil
}

// GetChassisResource defines the operations which handles the RPC request response
// for the getting the system resource of systems micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (cha *ChassisRPC) GetChassisResource(ctx context.Context, req *chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming getchassisResource request with %s", req.URL)
	var resp chassisproto.GetChassisResponse
	sessionToken := req.SessionToken
	authResp, err := cha.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		rewrite(ctx, authResp, &resp)
		return &resp, nil
	}
	var pc = chassis.PluginContact{
		ContactClient:   pmbhandle.ContactPlugin,
		DecryptPassword: common.DecryptWithPrivateKey,
		GetPluginStatus: scommon.GetPluginStatus,
	}
	data, _ := pc.GetChassisResource(ctx, req)
	rewrite(ctx, data, &resp)
	l.LogWithFields(ctx).Debugf("outgoing response for get chassisResource : %s", string(resp.Body))
	return &resp, nil
}

// GetChassisCollection defines the operation which handles the RPC request response
// for getting all the server chassis added.
// Retrieves all the keys with table name ChassisCollection and create the response
// to send back to requested user.
func (cha *ChassisRPC) GetChassisCollection(ctx context.Context, req *chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming GetChassisCollection request with %s", req.URL)
	var resp chassisproto.GetChassisResponse
	r := auth(ctx, cha.IsAuthorizedRPC, req.SessionToken, []string{common.PrivilegeLogin}, func() response.RPC {
		return cha.GetCollectionHandler.Handle(ctx)
	})
	rewrite(ctx, r, &resp)
	l.LogWithFields(ctx).Debugf("outgoing response Get ChassisCollection : %s", string(resp.Body))
	return &resp, nil
}

// GetChassisInfo defines the operations which handles the RPC request response
// for the getting the system resource of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (cha *ChassisRPC) GetChassisInfo(ctx context.Context, req *chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming GetChassisInfo request with %s", req.URL)
	var resp chassisproto.GetChassisResponse
	r := auth(ctx, cha.IsAuthorizedRPC, req.SessionToken, []string{common.PrivilegeLogin}, func() response.RPC {
		return cha.GetHandler.Handle(ctx, req)
	})

	rewrite(ctx, r, &resp)
	return &resp, nil
}

func rewrite(ctx context.Context, source response.RPC, target *chassisproto.GetChassisResponse) *chassisproto.GetChassisResponse {
	target.Header = source.Header
	target.StatusCode = source.StatusCode
	target.StatusMessage = source.StatusMessage
	target.Body = jsonMarshal(ctx, source.Body)
	return target
}

func jsonMarshal(ctx context.Context, input interface{}) []byte {
	if bytes, alreadyBytes := input.([]byte); alreadyBytes {
		return bytes
	}
	bytes, err := JSONMarshalFunc(input)
	if err != nil {
		l.LogWithFields(ctx).Println("error in unmarshalling response object from util-libs", err.Error())
	}
	return bytes
}

func generateResponse(ctx context.Context, input interface{}) []byte {
	if bytes, alreadyBytes := input.([]byte); alreadyBytes {
		return bytes
	}
	bytes, err := JSONMarshalFunc(input)
	if err != nil {
		l.LogWithFields(ctx).Error("error in unmarshalling response object from util-libs" + err.Error())
	}
	return bytes
}
