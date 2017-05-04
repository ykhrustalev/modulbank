package modulbank

import (
	"fmt"
)

// Requires
// account-info
func (client *Client) V1AccountInfo() ([]AccountInfo, error) {
	logger := client.logger.WithField("method", "V1AccountInfo")

	var result []AccountInfo
	_, err := client.handleRequest(logger, "POST", "/v1/account-info", nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Requires
// operation-history
func (client *Client) V1OperationHistory(account string, search *OperationHistorySearch) ([]Operation, error) {
	logger := client.logger.WithField("method", "V1AccountInfo")

	var result []Operation

	path := fmt.Sprintf("/v1/operation-history/%s", account)
	if search == nil {
		search = &OperationHistorySearch{}
	}

	_, err := client.handleRequest(logger, "POST", path, search, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
