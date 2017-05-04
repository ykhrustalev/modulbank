package modulbank

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

// Requires
// account-info
func (client *Client) V1_AccountInfo() ([]AccountInfo, error) {
	logger := client.logger.WithField("method", "V1_AccountInfo")

	var result []AccountInfo
	_, err := client.handleRequest(logger, "POST", "/v1/account-info", nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Requires
// operation-history
func (client *Client) V1_OperationHistory(bankAccountId string, search *OperationHistorySearch) ([]Operation, error) {
	logger := client.logger.WithField("method", "V1_OperationHistory")

	var result []Operation

	path := fmt.Sprintf("/v1/operation-history/%s", bankAccountId)
	if search == nil {
		search = &OperationHistorySearch{}
	}

	_, err := client.handleRequest(logger, "POST", path, search, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Requires
// account-info
func (client *Client) V1_AccountBalance(bankAccountId string) (result float64, err error) {
	logger := client.logger.WithField("method", "V1_AccountBalance")

	path := fmt.Sprintf("/v1/account-info/balance/%s", bankAccountId)

	req, err := client.buildRequest(logger, "POST", path, nil)
	if err != nil {
		return
	}

	resp, err := client.doRequest(logger, req, 200)
	defer client.closeBody(resp)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithError(err).Error("failed to read body")
		return
	}

	result, err = strconv.ParseFloat(string(body), 64)
	if err != nil {
		logger.WithError(err).Error("failed to parse value from body")
		return
	}

	return
}
