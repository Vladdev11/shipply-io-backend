package tasks

import (
	"fmt"
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

// SystemTaskChannel is a channel that is used to send tasks to the system task worker
var TaskChannel = make(chan *models.SystemTask, 100)

// InitalizeTaskProcesser is a function that is used to initialize the system task worker
func InitalizeTaskProcesser(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				task := <-TaskChannel
				ExecuteTask(task)
				//add a delay for cpu usage
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func StartAllTasks() {
	//get all tasks that are not completed
	tasks, err := models.GetAllTasksForProcessing()
	if err != nil {
		panic(err)
	}

	//send all tasks to the task channel
	for _, task := range tasks {
		//create a copy of the task so that the pointer is not the same as the one in the array
		//if we don't do this, the task passed to the channel will not be correct
		taskToRun := task
		TaskChannel <- &taskToRun
	}
}

// ExecuteTask is a function that is used to execute a system task
func ExecuteTask(st *models.SystemTask) {
	parameters, err := st.GetParameters()
	if err != nil {
		st.Errors = err.Error()
		st.HasError = true
		st.Update()
		return
	}

	st.StartedAt = time.Now()
	st.Processing = true
	err = st.Update()

	switch st.TaskType {
	case util.PullShopifyOrdersTaskType:
		err = PullShopifyOrdersTask(parameters["store_id"].(int), parameters["shop_name"].(string), parameters["access_token"].(string))
		if err != nil {
			st.Errors = err.Error()
			st.HasError = true
		}
		st.Update()
	case util.PullShopifyProductsTaskType:
		err = PullShopifyProductsTask(parameters["store_id"].(int), parameters["shop_name"].(string), parameters["access_token"].(string))
		if err != nil {
			st.Errors = err.Error()
			st.HasError = true
		}
		st.Update()
	default:
		fmt.Println("Unknown Task")
	}

	st.CompletedAt = time.Now()
	st.Processing = false
	st.Update()
}
