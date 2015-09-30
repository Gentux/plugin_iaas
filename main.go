/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dullgiulio/pingo"

	// vendor this dependency
	nan "nanocloud.com/plugins/iaas/libnan"
)

type IaasConfig struct {
	Url  string
	Port string
}

type Iaas struct{}

type VmInfo struct {
	Ico         string
	Name        string
	DisplayName string
	Status      string
	Locked      bool
}

var (
	g_IaasConfig IaasConfig
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (p *Iaas) Configure(jsonConfig string, _outMsg *string) error {
	var iaasConfig map[string]string

	err := json.Unmarshal([]byte(jsonConfig), &iaasConfig)
	if err != nil {
		r := fmt.Sprintf("ERROR: failed to unmarshal Iaas Plugin configuration : %s", err.Error())
		log.Printf(r)
		os.Exit(0)
		*_outMsg = r
		return nil
	}

	g_IaasConfig.Url = iaasConfig["Url"]
	g_IaasConfig.Port = iaasConfig["Port"]

	return nil
}

func (p *Iaas) ListRunningVm(jsonParams string, _outMsg *string) error {

	var (
		response struct {
			Result struct {
				DownloadingVmNames []string
				AvailableVMNames   []string
				BootingVmNames     []string
				RunningVmNames     []string
			}
			Error string
			Id    int
		}
		vmList      []VmInfo
		icon        string
		locked      bool
		status      string
		displayName string
	)
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", g_IaasConfig.Url, g_IaasConfig.Port),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		r := nan.NewExitCode(1, "ERROR: failed to contact Iaas API : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas API response : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	// TODO: Lots of Data aren't from iaas API
	for _, vmName := range response.Result.AvailableVMNames {
		if strings.Contains(vmName, "windows") {
			icon = "windows"
			locked = false
			displayName = "Windows Applications"
		} else {
			icon = "apps"
			locked = true
			displayName = "Haptic"
		}

		if stringInSlice(vmName, response.Result.RunningVmNames) {
			status = "running"
		} else if stringInSlice(vmName, response.Result.BootingVmNames) {
			status = "booting"
		} else if stringInSlice(vmName, response.Result.DownloadingVmNames) {
			status = "download"
		} else if stringInSlice(vmName, response.Result.AvailableVMNames) {
			status = "available"
		}
		vmList = append(vmList, VmInfo{
			Ico:         icon,
			Name:        vmName,
			DisplayName: displayName,
			Status:      status,
			Locked:      locked,
		})
	}

	jsonOuput, _ := json.Marshal(vmList)
	*_outMsg = string(jsonOuput)
	return err
}

func (p *Iaas) DownloadVm(jsonParams string, _outMsg *string) error {

	var (
		params = map[string]string{
			"vmname": jsonParams,
		}
		response struct {
			Result struct {
				Success bool
			}
		}
	)

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", g_IaasConfig.Url, g_IaasConfig.Port),
		"Iaas.Download",
		params,
	)
	if err != nil {
		r := nan.NewExitCode(1, "ERROR: failed to contact Iaas API : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas API response : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	if response.Result.Success == true {
		*_outMsg = "true"
	} else {
		*_outMsg = "false"
	}
	return nil
}

func (p *Iaas) DownloadStatus(jsonParams string, _outMsg *string) error {
	var (
		response struct {
			Result struct {
				AvailableVMNames   []string
				RunningVmNames     []string
				DownloadInProgress bool
			}
			Error string
			Id    int
		}
	)
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", g_IaasConfig.Url, g_IaasConfig.Port),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		r := nan.NewExitCode(1, "ERROR: failed to contact Iaas API : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas API response : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	if response.Result.DownloadInProgress {
		*_outMsg = "true"
	} else {
		*_outMsg = "false"
	}
	return err
}

func (p *Iaas) StartVm(jsonParams string, _outMsg *string) error {

	var (
		params = map[string]string{
			"name": jsonParams,
		}
		response struct {
			Result struct {
				Success bool
			}
		}
	)

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", g_IaasConfig.Url, g_IaasConfig.Port),
		"Iaas.Start",
		params,
	)
	if err != nil {
		r := nan.NewExitCode(1, "ERROR: failed to contact Iaas API : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas API response : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	if response.Result.Success == true {
		*_outMsg = "true"
	} else {
		*_outMsg = "false"
	}
	return nil
}

func (p *Iaas) StopVm(jsonParams string, _outMsg *string) error {

	var (
		params map[string]string
		vmName string
	)

	err := json.Unmarshal([]byte(jsonParams), &vmName)
	if err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Iaas.AccountParams : "+err.Error())
		log.Printf(r.Message) // for on-screen debug output
		*_outMsg = r.ToJson() // return codes for IPC should use JSON as much as possible
		return nil
	}

	params["vmname"] = vmName

	*_outMsg, err = jsonRpcRequest(
		g_IaasConfig.Url,
		"Iaas.Stop",
		params,
	)

	return err
}

func jsonRpcRequest(url string, method string, param map[string]string) (string, error) {

	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params":  []map[string]string{0: param},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return "", err
	}

	return string(body), nil
}

func main() {

	plugin := &Iaas{}

	pingo.Register(plugin)

	pingo.Run()
}
