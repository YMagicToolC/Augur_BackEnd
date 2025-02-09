package models

type APIRequest struct {
	BirthTime  string `json:"birthTime"`
	BirthPlace string `json:"birthPlace"`
	Contact    string `json:"contact"`
	Way        string `json:"way"`
	Gender     string `json:"gender"`
}

type APIResponseData struct {
	WorkflowRunID string `json:"workflow_run_id"`
	TaskID        string `json:"task_id"`
	Data          struct {
		ID          string    `json:"id"`
		WorkflowID  string    `json:"workflow_id"`
		Status      string    `json:"status"`
		Outputs     APIOutput `json:"outputs"`
		Error       *string   `json:"error"` // 使用指针以处理可能为 null 的情况
		ElapsedTime float64   `json:"elapsed_time"`
		TotalTokens int       `json:"total_tokens"`
		TotalSteps  int       `json:"total_steps"`
		CreatedAt   int64     `json:"created_at"`
		FinishedAt  int64     `json:"finished_at"`
	} `json:"data"`
}

type APIOutput struct {
	Message string `json:"message"` // 算命结果
	Contact string `json:"contact"`
}
