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

package system

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	"github.com/bharath-b-hpe/odimra/lib-utilities/errors"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"github.com/bharath-b-hpe/odimra/svc-aggregation/agmodel"
	"github.com/bharath-b-hpe/odimra/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

func (e *ExternalInterface) addPluginData(req AddResourceRequest, taskID, targetURI string, pluginContactRequest getResourceRequest) (response.RPC, string, []byte) {
	var resp response.RPC
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	if !(req.Oem.PreferredAuthType == "BasicAuth" || req.Oem.PreferredAuthType == "XAuthToken") {
		errMsg := "error: incorrect request property value for PreferredAuthType"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, errMsg, []interface{}{"PreferredAuthType", "[BasicAuth, XAuthToken]"}, taskInfo), "", nil
	}

	// checking the plugin type
	if !isPluginTypeSupported(req.Oem.PluginType) {
		errMsg := "error: incorrect request property value for PluginType"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, errMsg, []interface{}{"PluginType", config.Data.SupportedPluginTypes}, taskInfo), "", nil
	}

	// checking whether the Plugin already exists
	// If GetPluginData was successful, it indicates plugin already exists,
	// but it could also return errors, for below reasons, and has to be considered
	// as successful fetch of plugin data
	// error is nil - Plugin was successfully fetched.
	// error is not nil, Plugin data read but JSON unmarshaling failed
	// error is not nil, Plugin data read but decryption of plugin password failed
	// error is not nil, DB query failed, can't say for sure if queried plugin exists,
	// except when read fails with plugin data not found, and will continue with add process,
	// and any other errors, will fail add plugin operation.
	_, errs := agmodel.GetPluginData(req.Oem.PluginID)
	if errs == nil || (errs != nil && (errs.ErrNo() == errors.JSONUnmarshalFailed || errs.ErrNo() == errors.DecryptionFailed)) {
		errMsg := "error:plugin with name " + req.Oem.PluginID + " already exists"
		log.Println(errMsg)
		return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, []interface{}{"Plugin", "PluginID", req.Oem.PluginID}, taskInfo), "", nil
	}
	if errs != nil && errs.ErrNo() != errors.DBKeyNotFound {
		errMsg := "error: DB lookup failed for " + req.Oem.PluginID + " plugin: " + errs.Error()
		log.Println(errMsg)
		if errs.ErrNo() == errors.DBConnFailed {
			return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errMsg,
				[]interface{}{"Backend", config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, taskInfo), "", nil
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, []interface{}{}, taskInfo), "", nil
	}

	// encrypt plugin password
	ciphertext, err := e.EncryptPassword([]byte(req.Password))
	if err != nil {
		errMsg := "error: encryption failed: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}

	ipData := strings.Split(req.ManagerAddress, ":")
	var plugin = agmodel.Plugin{
		IP:                ipData[0],
		Port:              ipData[1],
		Username:          req.UserName,
		Password:          []byte(req.Password),
		ID:                req.Oem.PluginID,
		PluginType:        req.Oem.PluginType,
		PreferredAuthType: req.Oem.PreferredAuthType,
	}
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while creating the session: ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	// Verfiying the plugin Status
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.OID = "/ODIM/v1/Status"
	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}
	// extracting the EMB Type and EMB Queue name
	var statusResponse common.StatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	var queueList = make([]string, 0)
	if statusResponse.EventMessageBus != nil {
		for i := 0; i < len(statusResponse.EventMessageBus.EmbQueue); i++ {
			queueList = append(queueList, statusResponse.EventMessageBus.EmbQueue[i].QueueName)
		}
	}

	// Getting all managers info from plugin
	pluginContactRequest.OID = "/ODIM/v1/Managers"
	body, _, getResponse, err = contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}
	//  Extract all managers info and loop  over each members
	managersMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &managersMap)
	if err != nil {
		errMsg := "unable to parse the managers resposne" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	var managersData = make(map[string][]byte)
	managerMembers := managersMap["Members"]

	// Getting the indivitual managers response
	for _, object := range managerMembers.([]interface{}) {
		pluginContactRequest.OID = object.(map[string]interface{})["@odata.id"].(string)
		body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
		}
		managersData[pluginContactRequest.OID] = body
	}
	e.SubscribeToEMB(plugin.ID, queueList)
	// saving all plugin manager data
	var listMembers = make([]agresponse.ListMember, 0)
	for oid, data := range managersData {
		dbErr := agmodel.GenericSave(updateManagerName(data, plugin.ID), "Managers", oid)
		if err != nil {
			errMsg := "error: while saving the plugin data with generic save: " + dbErr.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
		}
		listMembers = append(listMembers, agresponse.ListMember{
			OdataID: oid,
		})
	}

	// store encrypted password
	plugin.Password = ciphertext
	// saving the pluginData
	dbErr := agmodel.SavePluginData(plugin)
	if dbErr != nil {
		errMsg := "error: while saving the plugin data: " + dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		"Location":     listMembers[0].OdataID,
	}
	var managersList = make([]string, 0)
	for i := 0; i < len(listMembers); i++ {
		managersList = append(managersList, listMembers[i].OdataID)
	}
	e.PublishEvent(managersList, "ManagerCollection")
	log.Println("sucessfully added  plugin with the id ", req.Oem.PluginID)
	return resp, uuid.NewV4().String(), ciphertext
}