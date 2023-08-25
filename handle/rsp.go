package handle

type Rsp struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StatusData struct {
	GlobalMode  bool `json:"global_mode"`
	Inited      bool `json:"inited"`
	Running     bool `json:"running"`
	SystemProxy bool `json:"system_proxy"`
}
