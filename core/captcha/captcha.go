package captcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (c *CaptchaClient) GetBalance() (*GetBalanceResponse, error) {
	body, err := json.Marshal(GetBalanceRequest{
		Apikey: c.Apikey,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Host+"/getBalance", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var reply GetBalanceResponse

	err = json.NewDecoder(res.Body).Decode(&reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *CaptchaClient) CreateTask(task interface{}) (*CreateTaskResponse, error) {
	body, err := json.Marshal(CreateTaskRequest{
		Apikey: c.Apikey,
		Task:   task,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Host+"/createTask", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var reply CreateTaskResponse

	err = json.NewDecoder(res.Body).Decode(&reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *CaptchaClient) GetResult(taskID int) (*TaskResultResponse, error) {
	body, err := json.Marshal(TaskResultRequest{
		Apikey: c.Apikey,
		TaskID: taskID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Host+"/getTaskResult", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var reply TaskResultResponse

	err = json.NewDecoder(res.Body).Decode(&reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *CaptchaClient) JoinTaskResult(taskID int) (*TaskResultResponse, error) {
	for i := 1; i <= 25; i++ {
		result, err := c.GetResult(taskID)
		if err != nil {
			return nil, err
		}

		if result.Status == "ready" {
			return result, nil
		} else if result.Status == "processing" {
			time.Sleep(1 * time.Second)
		} else {
			return nil, errors.New("failed to solve task: " + result.Status)
		}
	}

	return nil, errors.New("failed to solve task: timeout")
}
