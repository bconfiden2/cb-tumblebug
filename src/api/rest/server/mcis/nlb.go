/*
Copyright 2019 The Cloud-Barista Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package mcis is to handle REST API for mcis
package mcis

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-tumblebug/src/core/common"
	"github.com/cloud-barista/cb-tumblebug/src/core/mcis"
	"github.com/labstack/echo/v4"
)

// RestPostNLB godoc
// @Summary Create NLB
// @Description Create NLB
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nsId path string true "Namespace ID" default(ns01)
// @Param option query string false "Option: [required params for register] connectionName, name, cspNLBId" Enums(register)
// @Param nlbReq body mcis.TbNLBReq true "Details of the NLB object"
// @Success 200 {object} mcis.TbNLBInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb [post]
func RestPostNLB(c echo.Context) error {

	nsId := c.Param("nsId")

	optionFlag := c.QueryParam("option")

	u := &mcis.TbNLBReq{}
	if err := c.Bind(u); err != nil {
		return err
	}

	fmt.Println("[POST NLB]")

	content, err := mcis.CreateNLB(nsId, u, optionFlag)

	if err != nil {
		common.CBLog.Error(err)
		mapA := map[string]string{"message": err.Error()}
		return c.JSON(http.StatusInternalServerError, &mapA)
	}
	return c.JSON(http.StatusCreated, content)
}

/* function RestPutNLB not yet implemented
// RestPutNLB godoc
// @Summary Update NLB
// @Description Update NLB
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nlbInfo body mcis.TbNLBInfo true "Details of the NLB object"
// @Success 200 {object} mcis.TbNLBInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb/{nlbId} [put]
*/
func RestPutNLB(c echo.Context) error {
	//nsId := c.Param("nsId")

	return nil
}

// RestGetNLB godoc
// @Summary Get NLB
// @Description Get NLB
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nsId path string true "Namespace ID" default(ns01)
// @Param nlbId path string true "NLB ID"
// @Success 200 {object} mcis.TbNLBInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb/{nlbId} [get]
func RestGetNLB(c echo.Context) error {

	nsId := c.Param("nsId")

	resourceId := c.Param("resourceId")

	res, err := mcis.GetNLB(nsId, resourceId)
	if err != nil {
		mapA := map[string]string{"message": "Failed to find the NLB " + resourceId}
		return c.JSON(http.StatusNotFound, &mapA)
	} else {
		return c.JSON(http.StatusOK, &res)
	}
}

// Response structure for RestGetAllNLB
type RestGetAllNLBResponse struct {
	NLB []mcis.TbNLBInfo `json:"nlb"`
}

// RestGetAllNLB godoc
// @Summary List all NLBs or NLBs' ID
// @Description List all NLBs or NLBs' ID
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nsId path string true "Namespace ID" default(ns01)
// @Param option query string false "Option" Enums(id)
// @Param filterKey query string false "Field key for filtering (ex: cspNLBName)"
// @Param filterVal query string false "Field value for filtering (ex: ns01-alibaba-ap-northeast-1-vpc)"
// @Success 200 {object} JSONResult{[DEFAULT]=RestGetAllNLBResponse,[ID]=common.IdList} "Different return structures by the given option param"
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb [get]
func RestGetAllNLB(c echo.Context) error {

	nsId := c.Param("nsId")

	optionFlag := c.QueryParam("option")
	filterKey := c.QueryParam("filterKey")
	filterVal := c.QueryParam("filterVal")

	if optionFlag == "id" {
		content := common.IdList{}
		var err error
		content.IdList, err = mcis.ListNLBId(nsId)
		if err != nil {
			mapA := map[string]string{"message": "Failed to list NLBs' ID; " + err.Error()}
			return c.JSON(http.StatusNotFound, &mapA)
		}

		return c.JSON(http.StatusOK, &content)
	} else {

		resourceList, err := mcis.ListNLB(nsId, filterKey, filterVal)
		if err != nil {
			mapA := map[string]string{"message": "Failed to list NLBs; " + err.Error()}
			return c.JSON(http.StatusNotFound, &mapA)
		}

		var content struct {
			NLB []mcis.TbNLBInfo `json:"nlb"`
		}

		content.NLB = resourceList.([]mcis.TbNLBInfo) // type assertion (interface{} -> array)
		return c.JSON(http.StatusOK, &content)
		// return c.JSON(http.StatusBadRequest, nil)
	}
}

// RestDelNLB godoc
// @Summary Delete NLB
// @Description Delete NLB
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nsId path string true "Namespace ID" default(ns01)
// @Param nlbId path string true "NLB ID"
// @Success 200 {object} common.SimpleMsg
// @Failure 404 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb/{nlbId} [delete]
func RestDelNLB(c echo.Context) error {

	nsId := c.Param("nsId")

	resourceId := c.Param("resourceId")

	forceFlag := c.QueryParam("force")

	err := mcis.DelNLB(nsId, resourceId, forceFlag)
	if err != nil {
		common.CBLog.Error(err)
		mapA := map[string]string{"message": err.Error()}
		return c.JSON(http.StatusInternalServerError, &mapA)
	}

	mapA := map[string]string{"message": "The NLB " + resourceId + " has been deleted"}
	return c.JSON(http.StatusOK, &mapA)
}

// RestDelAllNLB godoc
// @Summary Delete all NLBs
// @Description Delete all NLBs
// @Tags [Infra resource] NLB management
// @Accept  json
// @Produce  json
// @Param nsId path string true "Namespace ID" default(ns01)
// @Param match query string false "Delete resources containing matched ID-substring only" default()
// @Success 200 {object} common.IdList
// @Failure 404 {object} common.SimpleMsg
// @Router /ns/{nsId}/nlb [delete]
func RestDelAllNLB(c echo.Context) error {

	nsId := c.Param("nsId")

	forceFlag := c.QueryParam("force")
	subString := c.QueryParam("match")

	output, err := mcis.DelAllNLB(nsId, subString, forceFlag)
	if err != nil {
		common.CBLog.Error(err)
		mapA := map[string]string{"message": err.Error()}
		return c.JSON(http.StatusConflict, &mapA)
	}

	return c.JSON(http.StatusOK, output)
}
