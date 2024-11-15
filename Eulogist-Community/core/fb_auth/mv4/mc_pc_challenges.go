package fbauth

import (
	fb_client "Eulogist/core/fb_auth/mv4/client"
	"bytes"
	"encoding/json"
	"fmt"
)

// ...
type FNumRequest struct {
	Data string `json:"data"`
}

// ...
func TransferData(client *fb_client.Client, content string) string {
	r, err := client.HttpClient.Get(fmt.Sprintf("%s/api/phoenix/transfer_start_type?content=%s", client.AuthServer, content))
	if err != nil {
		panic(err)
	}
	resp := fb_client.AssertAndParse[map[string]any](r)
	succ, _ := resp["success"].(bool)
	if !succ {
		err_m, _ := resp["message"].(string)
		panic(fmt.Sprintf("Failed to transfer start type: %s", err_m))
	}
	data, _ := resp["data"].(string)
	return data
}

// ...
func TransferCheckNum(client *fb_client.Client, data string) string {
	rspreq := &FNumRequest{
		Data: data,
	}
	msg, err := json.Marshal(rspreq)
	if err != nil {
		panic("Failed to encode json")
	}
	r, err := client.HttpClient.Post(fmt.Sprintf("%s/api/phoenix/transfer_check_num", client.AuthServer), "application/json", bytes.NewBuffer(msg))
	if err != nil {
		panic(err)
	}
	resp := fb_client.AssertAndParse[map[string]any](r)
	succ, _ := resp["success"].(bool)
	if !succ {
		err_m, _ := resp["message"].(string)
		panic(fmt.Sprintf("Failed to transfer check num: %s", err_m))
	}
	val, _ := resp["value"].(string)
	return val
}
