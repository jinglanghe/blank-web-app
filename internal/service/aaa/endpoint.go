package aaa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/model/aaa"
	"github.com/apulis/sdk/go-utils/logging"
	"net/http"
)

func RegisterEndPoints(p *aaa.EndPointsAndPolicies) error {
	aaaDomain := getAaaDomain()
	url := fmt.Sprintf("http://%s/api/v1/endpoints", aaaDomain)

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(p)
	if err != nil {
		logging.Error(err).Msg("register endpoints failed at json encoder")
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logging.Error(err).Msg("register endpoints failed at new request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logging.Error(err).Msg("")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("register endpoints failed, status code: %d", resp.StatusCode)
		logging.Error(err).Msg("")
		return err
	}
	return nil
}
