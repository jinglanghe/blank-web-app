package aaa

import (
	"encoding/json"
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model/aaa"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	"io"
	"net/http"
)

type orgDetail struct {
	Data itemsOrganization `json:"data"`
}

type itemsOrganization struct {
	Items []aaa.Organization `json:"items"`
}

func OrgDetail(token string, orgId int64) (o *aaa.Organization, err error) {
	aaaDomain := getAaaDomain()
	url := fmt.Sprintf("http://%s/api/v1/orgs/%d", aaaDomain, orgId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		logging.Error(err).Msg("")
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("get org %d detail failed, status code: %d", orgId, resp.StatusCode)
		logging.Error(err).Msg("")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	orgDetail := orgDetail{}
	err = json.Unmarshal(body, &orgDetail)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	if len(orgDetail.Data.Items) == 0 {
		return &aaa.Organization{}, nil
	}

	return &orgDetail.Data.Items[0], nil
}

func DefaultOrg() *aaa.Organization {
	defaultOrg := aaa.Organization{
		Base: model.Base{
			ID: int64(DefaultOrgId),
		},
		Account: DefaultOrgName,
	}
	return &defaultOrg
}
