package onecloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-onecloud/internal/ocp_client"
	"time"
)

const OperationStatusInProgress = "in progress"
const OperationStatusFailed = "failed"
const OperationStatusAborted = "aborted"
const OperationStatusSucceeded = "succeeded"

func waitForOperationSuccess(ctx context.Context, client ocp_client.Client, operationID string, timeout time.Duration) (interface{}, error) {
	pending := []string{
		OperationStatusInProgress,
	}
	target := []string{
		OperationStatusSucceeded,
	}

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    OperationStateRefresh(ctx, client, operationID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"operation failed, operation_id: %s, err: %s",
			operationID, err)
	}

	return result, nil
}

func OperationStateRefresh(ctx context.Context, client ocp_client.Client, operationID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		operation, code, err := client.GetOperation(ctx, operationID)
		if err != nil {
			if code == 0 {
				return nil, OperationStatusInProgress, nil
			}
			return nil, "", err
		}
		var createErr error
		operationStatus := operation["status"].(string)
		if operationStatus == OperationStatusFailed || operationStatus == OperationStatusAborted {
			operationType := operation["operation_type"].(string)
			progress := operation["progress"].(map[string]interface{})
			steps := progress["steps_details"].([]interface{})
			latestStep := steps[int(progress["completed_steps"].(float64))]
			stepName := latestStep.(map[string]interface{})["name"].(string)
			createErr = fmt.Errorf("%s, at the step: %s", operationType, stepName)
		}
		return operation, operationStatus, createErr
	}
}
