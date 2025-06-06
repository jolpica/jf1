package uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type RequestResult struct {
	StatusCode   int
	RequestData  JolpicaUploadRequestPayload
	ResponseData JolpicaUploadResponsePayload
	ProcessedDir ProcessedDirectory

	Err error
}

func sendDataLoadRequests(ctx context.Context, dirLoadResultc <-chan DirectoryLoadResult, config UploadConfig, token string) <-chan RequestResult {
	requestResultc := make(chan RequestResult, 5)
	client := &http.Client{}

	maxConcurrentRequests := config.MaxConcurrentRequests
	reqSem := make(chan struct{}, maxConcurrentRequests)

	go func() {
		var wg sync.WaitGroup
		defer func() { wg.Wait(); close(requestResultc); close(reqSem) }()
		for result := range dirLoadResultc {
			wg.Add(1)
			go func(result DirectoryLoadResult) {
				defer wg.Done()
				if result.Err != nil {
					select {
					case requestResultc <- RequestResult{Err: result.Err}:
					case <-ctx.Done():
					}
					return
				}
				if config.OnlyUpdateUploaded {
					// Skip making requests and pass on the result directly
					requestResultc <- RequestResult{
						StatusCode:   0,
						RequestData:  JolpicaUploadRequestPayload{},
						ResponseData: JolpicaUploadResponsePayload{},
						ProcessedDir: *result.Result,
						Err:          nil,
					}
					return
				}
				makeDataLoadRequest(ctx, requestResultc, reqSem, client, result.Result, config, token)
			}(result)
		}
	}()

	return requestResultc
}

type JolpicaUploadRequestPayload struct {
	DryRun      bool             `json:"dry_run"`
	Description string           `json:"description"`
	Data        []map[string]any `json:"data"`
}
type JolpicaUploadResponsePayload struct {
	UpdatedCount int                                             `json:"updated_count"`
	CreatedCount int                                             `json:"created_count"`
	TotalCount   int                                             `json:"total_count"`
	Models       map[string]JolpicaUploadResponsePerModelPayload `json:"models"`

	Errors []map[string]any
}
type JolpicaUploadResponsePerModelPayload struct {
	UpdatedCount int   `json:"updated_count"`
	CreatedCount int   `json:"created_count"`
	Created      []int `json:"created"`
	Updated      []int `json:"updated"`
}

func makeDataLoadRequest(ctx context.Context, requestResultc chan RequestResult, reqSem chan struct{}, client *http.Client, processedDir *ProcessedDirectory, config UploadConfig, token string) {
	var requestResult RequestResult
	defer func() {
		select {
		case requestResultc <- requestResult:
		case <-ctx.Done():
		}
	}()

	payload := JolpicaUploadRequestPayload{
		DryRun:      config.DryRun,
		Description: processedDir.Description(),
		Data:        processedDir.Data,
	}

	request, err := createJolpicaHttpRequest(payload, config, token)
	if err != nil {
		requestResult = RequestResult{Err: fmt.Errorf("error generating request (%s): %v", processedDir.Description(), err)}
		return
	}

	reqSem <- struct{}{}
	resp, err := client.Do(request)
	<-reqSem

	if err != nil {
		requestResult = RequestResult{Err: fmt.Errorf("error sending request (%s): %v", processedDir.Description(), err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		requestResult = RequestResult{Err: fmt.Errorf("error decoding response (%s): %v", processedDir.Description(), err)}
		return
	}

	var respData JolpicaUploadResponsePayload
	if err := json.Unmarshal(body, &respData); err != nil {
		requestResult = RequestResult{Err: fmt.Errorf("error unmarshalling response (%s): %v", processedDir.Description(), err)}
		return
	}

	requestResult = RequestResult{
		ProcessedDir: *processedDir,
		StatusCode:   resp.StatusCode,
		RequestData:  payload,
		ResponseData: respData,
	}
}

func createJolpicaHttpRequest(payload JolpicaUploadRequestPayload, config UploadConfig, token string) (*http.Request, error) {
	url := fmt.Sprintf("%s/data/import/", config.BaseUrl)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	request.Header = http.Header{
		// TODO: Pass in token as parameter
		"Authorization": []string{fmt.Sprintf("Token %s", token)},
		"Content-Type":  []string{"application/json"},
	}
	return request, nil
}
