package modules

import (
	"GoScheduler/internal/models"
	rpcClient "GoScheduler/internal/modules/rpc/client"
	pb "GoScheduler/internal/modules/rpc/proto"
	"GoScheduler/internal/task"
	"fmt"
)

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel models.Task, taskUniqueId uint) (result string, err error) {
	taskRequest := new(pb.TaskRequest)
	taskRequest.Timeout = int32(taskModel.Timeout)
	taskRequest.Command = taskModel.Command
	taskRequest.Id = int64(taskUniqueId)
	resultChan := make(chan task.TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		go func(th models.TaskHostDetail) {
			output, err := rpcClient.Exec(th.Name, th.Port, taskRequest)
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s:%d]\n%s\n%s\n\n",
				th.Alias, th.Name, th.Port, errorMessage, output,
			)
			resultChan <- task.TaskResult{Err: err, Result: outputMessage}
		}(taskHost)
	}

	var aggregationErr error = nil
	aggregationResult := ""
	for i := 0; i < len(taskModel.Hosts); i++ {
		taskResult := <-resultChan
		aggregationResult += taskResult.Result
		if taskResult.Err != nil {
			aggregationErr = taskResult.Err
		}
	}

	return aggregationResult, aggregationErr
}
