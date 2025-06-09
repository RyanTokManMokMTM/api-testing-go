package generator

// TestCase represents a test case
type TestCase struct {
	Name        string
	Description string
	Steps       []TestStep
}

// TestStep represents a test step
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

// FromResponse represents a variable extracted from a response
type FromResponse struct {
	Step      string
	Name      string
	FromField string
}

// ResponseCheck represents a response check
type ResponseCheck struct {
	Type      string
	Field     string
	Value     interface{}
	Regex     string
	IsPresent bool
}
