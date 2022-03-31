package aaa

import (
	"encoding/json"
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model/aaa"
	"github.com/apulis/sdk/go-utils/logging"
	"io"
	"net/http"
)

type userCurrent struct {
	Data aaa.User `json:"data"`
}

// 获取当前用户的信息
func UserCurrent(token string) (u *aaa.User, err error) {
	aaaDomain := getAaaDomain()
	url := fmt.Sprintf("http://%s/api/v1/users/current", aaaDomain)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		logging.Error(err).Msg("")
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("get current user info failed, status code: %d", resp.StatusCode)
		logging.Error(err).Msg("")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	userCurrent := userCurrent{}
	err = json.Unmarshal(body, &userCurrent)
	if err != nil {
		logging.Error(err).Msg("")
		return
	}

	return &userCurrent.Data, nil
}
