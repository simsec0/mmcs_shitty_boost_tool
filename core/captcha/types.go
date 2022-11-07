package captcha

type CaptchaClient struct {
	Host   string
	Apikey string
}

type Task struct {
	Type    string `json:"type"`
	Host    string `json:"websiteURL"`
	Sitekey string `json:"websiteKey"`

	Data      string `json:"data"`
	UserAgent string `json:"userAgent"`

	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     string `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin"`
	ProxyPassword string `json:"proxyPassword"`
}

type CreateTaskRequest struct {
	Apikey string      `json:"clientKey"`
	Task   interface{} `json:"task"`
}

type CreateTaskResponse struct {
	TaskID int `json:"taskId"`
}

type TaskResultRequest struct {
	Apikey string `json:"clientKey"`
	TaskID int    `json:"taskId"`
}

type TaskResultResponse struct {
	Status   string `json:"status"`
	Solution struct {
		GeneratedPassUUID string `json:"gRecaptchaResponse"`
	}
}

type GetBalanceRequest struct {
	Apikey string `json:"clientKey"`
}

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
}
