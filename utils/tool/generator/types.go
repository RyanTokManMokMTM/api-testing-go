package generator

// TestCase 表示一個測試用例
type TestCase struct {
	Name        string
	Description string
	Steps       []TestStep
}

// TestStep 表示測試步驟
type TestStep struct {
	Name           string
	Method         string
	URI            string
	Headers        map[string]string
	Body           interface{}
	ExpectedCode   string
	ExpectedStatus int
	ExpectedBody   map[string]interface{}
	Variables      []string
	FromResponses  []FromResponse
	ResponseChecks []ResponseCheck
}

// FromResponse 表示從響應中提取的變量
type FromResponse struct {
	Step      string
	Name      string
	FromField string
}

// ResponseCheck 表示響應檢查
type ResponseCheck struct {
	Type      string
	Field     string
	Value     interface{}
	Regex     string
	IsPresent bool
}
