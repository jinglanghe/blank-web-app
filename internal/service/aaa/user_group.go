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

type userGroupDetail struct {
	Data itemsUserGroup `json:"data"`
}

type itemsUserGroup struct {
	Items []aaa.UserGroup `json:"items"`
}

func UserGroupDetail(token string, userGroupId int64) (o *aaa.UserGroup, err error) {
	aaaDomain := getAaaDomain()
	url := fmt.Sprintf("http://%s/api/v1/user-groups/%d", aaaDomain, userGroupId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		logging.Error(err).Msg("")
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("get user group %d detail failed, status code: %d", userGroupId, resp.StatusCode)
		logging.Error(err).Msg("")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	userGroupDetail := userGroupDetail{}
	err = json.Unmarshal(body, &userGroupDetail)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	if len(userGroupDetail.Data.Items) == 0 {
		return &aaa.UserGroup{}, nil
	}

	return &userGroupDetail.Data.Items[0], nil
}

func SystemAdminGroup() *aaa.UserGroup {
	systemAdminGroup := aaa.UserGroup{
		Base: model.Base{
			ID: int64(SysAdminUgId),
		},
		Account: SysAdminUgName,
	}
	return &systemAdminGroup
}

func DefaultUserGroup() *aaa.UserGroup {
	defaultOrg := DefaultOrg()
	defaultGroup := aaa.UserGroup{
		Base: model.Base{
			ID: int64(OrgAdminUgId),
		},
		Account:        OrgAdminUgName,
		OrganizationID: defaultOrg.ID,
		Organization:   *defaultOrg,
	}
	return &defaultGroup
}
