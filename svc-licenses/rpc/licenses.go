//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
	"net/http"
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	lgr "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
)

var podName = os.Getenv("POD_NAME")

// GetLicenseService to get license service details
func (l *Licenses) GetLicenseService(ctx context.Context, req *licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.LicenseService, podName)
	resp := &licenseproto.GetLicenseResponse{}
	authResp, err := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			lgr.Log.Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, l.connector.GetLicenseService(req))
	return resp, nil
}

// GetLicenseCollection to get license collection
func (l *Licenses) GetLicenseCollection(ctx context.Context, req *licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.LicenseService, podName)
	resp := &licenseproto.GetLicenseResponse{}
	authResp, err := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			lgr.Log.Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, l.connector.GetLicenseCollection(ctx, req))
	return resp, nil
}

// GetLicenseResource to get license resource
func (l *Licenses) GetLicenseResource(ctx context.Context, req *licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.LicenseService, podName)
	resp := &licenseproto.GetLicenseResponse{}
	authResp, err := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			lgr.Log.Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, l.connector.GetLicenseResource(ctx, req))
	return resp, nil
}

// InstallLicenseService to install license
func (l *Licenses) InstallLicenseService(ctx context.Context, req *licenseproto.InstallLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.LicenseService, podName)
	resp := &licenseproto.GetLicenseResponse{}
	authResp, err := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			lgr.Log.Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, l.connector.InstallLicenseService(ctx, req))
	return resp, nil
}
