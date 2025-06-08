package apitestrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	testCfg "github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/config/types"
	apierror "github.com/RyanTokManMokMTM/api-testing-go/utils/error"
	apihelper "github.com/RyanTokManMokMTM/api-testing-go/utils/helper"
	apiutil "github.com/RyanTokManMokMTM/api-testing-go/utils/util"

	"github.com/onsi/ginkgo/v2"

	"github.com/samber/lo"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

//nolint:gochecknoglobals // local golobal template
var (
	uriTemp    = template.New("request")
	headerTemp = template.New("header")
	queryTemp  = template.New("query")
	bodyTemp   = template.New("body")
	respTemp   = template.New("response")
)

type apiTestRunnerParams struct {
	url             string                                      // url
	body            string                                      // body
	headerMap       map[string]string                           // headerMap
	queryMap        map[string]string                           // queryMap
	respAssestFuncs []func(*http.Response, *http.Request) error // respAssestFuncs
}

func (apiTestRunner *APITestRunner) getValueFromResponse(
	respCache *map[string]any,
	name string,
	step string,
	formField string,
) (any, error) {
	var data any
	var err error
	prevResp, ok := (*respCache)[step]
	if !ok {
		return nil, apierror.ErrStepRespNotFound
	}
	data, err = apihelper.GetRespFieldData(formField, prevResp)
	if err != nil {
		return nil, err
	}

	(*respCache)[name] = data
	return data, nil
}

// return (url, body, headerMap, queryMap, respAssestFuncs, err).
func (apiTestRunner *APITestRunner) parseRequestData(
	variables map[string]any,
	workflow testCfg.Workflow,
) (
	*apiTestRunnerParams,
	error, // err
) {
	headerMap, err := apiTestRunner.parseHeader(variables, workflow.Request.Headers)
	if err != nil {
		return nil, err
	}

	url, err := apiTestRunner.parseURI(variables, workflow.Request.URI)
	if err != nil {
		return nil, err
	}

	queryMap, err := apiTestRunner.parseQuery(variables, workflow.Request.Query)
	if err != nil {
		return nil, err
	}

	bodyStr, err := apiTestRunner.parseBody(variables, workflow.Request.Body)
	if err != nil {
		return nil, err
	}

	//nolint:bodyclose // closed by defer func
	respAssestFuncs, err := apiTestRunner.setAssertFunc(variables, workflow)
	if err != nil {
		return nil, err
	}
	return &apiTestRunnerParams{
		url:             url,
		body:            bodyStr,
		headerMap:       headerMap,
		queryMap:        queryMap,
		respAssestFuncs: respAssestFuncs,
	}, nil
}

//nolint:bodyclose // is closed by end()
func (apiTestRunner *APITestRunner) setAssertFunc(
	variables map[string]any,
	workflow testCfg.Workflow,
) ([]func(*http.Response, *http.Request) error, error) {
	// handle expected response
	assertFuncs := make([]func(*http.Response, *http.Request) error, 0)

	expectedChain := jsonpath.Chain()
	expectedChain.Equal("$.code", workflow.ExpectResponse.Code) // Must

	// For all equals data
	for _, equal := range workflow.ExpectResponse.Body.Equals {
		expectedChain.Present("$.data")
		field := fmt.Sprintf("$.%s", equal.Field)

		if err := apiTestRunner.handleEqualValue(field, equal.Type, equal.Value, variables, expectedChain); err != nil {
			return nil, err
		}
	}

	// mathching data with regex
	for _, match := range workflow.ExpectResponse.Body.Matches {
		field := fmt.Sprintf("$.%s", match.Field)
		expectedChain.Matches(field, match.Regex)
	}

	// check if the data is present
	for _, notNull := range workflow.ExpectResponse.Body.Presents {
		field := fmt.Sprintf("$.%s", notNull.Field)
		expectedChain.Present(field)
	}
	//nolint:bodyclose // is closed by end()
	assertFuncs = append(assertFuncs, expectedChain.End())

	// check is not nil
	// Json Path.
	for _, notPresent := range workflow.ExpectResponse.Body.NotPresents {
		field := fmt.Sprintf("$.%s", notPresent.Field)
		assertFuncs = append(assertFuncs, jsonpath.NotPresent(field))
	}

	for _, greaterThan := range workflow.ExpectResponse.Body.GreaterThans {
		field := fmt.Sprintf("$.%s", greaterThan.Field)
		assertFuncs = append(assertFuncs, jsonpath.GreaterThan(field, greaterThan.Value))
	}

	for _, lessThan := range workflow.ExpectResponse.Body.LessThans {
		field := fmt.Sprintf("$.%s", lessThan.Field)
		assertFuncs = append(assertFuncs, jsonpath.LessThan(field, lessThan.Value))
	}

	return assertFuncs, nil
}

func (apiTestRunner *APITestRunner) handleEqualValue(
	field string,
	equalType types.FieldType,
	value any,
	variables map[string]any,
	expectedChain *jsonpath.AssertionChain,
) error {
	if equalType == "" {
		switch v := value.(type) {
		case int:
			expectedChain.Equal(field, float64(v))
		case int8:
			expectedChain.Equal(field, float64(v))
		case int16:
			expectedChain.Equal(field, float64(v))
		case int32:
			expectedChain.Equal(field, float64(v))
		case int64:
			expectedChain.Equal(field, float64(v))
		case uint:
			expectedChain.Equal(field, float64(v))
		case uint8:
			expectedChain.Equal(field, float64(v))
		case uint16:
			expectedChain.Equal(field, float64(v))
		case uint32:
			expectedChain.Equal(field, float64(v))
		case uint64:
			expectedChain.Equal(field, float64(v))
		default:
			expectedChain.Equal(field, value)
		}
		return nil
	}

	// MARK: parse string with template and cast to the type of field
	// eg: "{{.mid}}" -> "123"
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("equal value '%+v', must be a string type, got %T, %w", value, value, apierror.ErrTypeNotSupported)
	}

	// the string should be include '{{' , '}}' and .
	pattern := `\{\{\s*\.\w+(?:\s+\w+)*\s*\}\}`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return fmt.Errorf("invalid template syntax: must be in format '{{.variableName}}' or '{{ .variableName }}', got: %s, %w", str, apierror.ErrTypeNotSupported)
	}

	if !strings.Contains(str, ".") {
		return fmt.Errorf("template variable must contain a dot (.) to reference a variable, got: %s, %w", str, apierror.ErrTypeNotSupported)
	}

	if matchedErr := apiTestRunner.isExistVars(str, variables); matchedErr != nil {
		return fmt.Errorf("template variable '%s' references an undefined variable, %w", str, matchedErr)
	}

	// parse template
	str, err := apiTestRunner.parseTemplate(respTemp, str, variables)
	if err != nil {
		return err
	}

	// the convert to the type of field
	// eg: "number" -> int
	convertedValue, err := apiTestRunner.castStrToSpecificType(str, equalType)
	if err != nil {
		return err
	}
	expectedChain.Equal(field, convertedValue)
	return nil
}

func (apiTestRunner *APITestRunner) castStrToSpecificType(str string, typeStr types.FieldType) (any, error) {
	value, err := apihelper.ConvertStrToType(str, typeStr)
	if err != nil {
		return nil, err
	}
	if v, converted := value.(int64); converted {
		return float64(v), nil
	}
	return value, nil
}

func (apiTestRunner *APITestRunner) parseHeader(variables map[string]any, reqHeader string) (map[string]string, error) {
	if reqHeader == "" {
		return map[string]string{}, nil
	}

	if matchedErr := apiTestRunner.isExistVars(reqHeader, variables); matchedErr != nil {
		return nil, matchedErr
	}

	str, err := apiTestRunner.parseTemplate(headerTemp, reqHeader, variables)
	if err != nil {
		return nil, err
	}
	headerMap := make(map[string]string)

	err = json.Unmarshal([]byte(str), &headerMap)
	if err != nil {
		return nil, err
	}

	return headerMap, nil
}

func (apiTestRunner *APITestRunner) parseURI(variables map[string]any, uri string) (string, error) {
	if uri == "" {
		return "", nil
	}

	if matchedErr := apiTestRunner.isExistVars(uri, variables); matchedErr != nil {
		return "", matchedErr
	}

	str, err := apiTestRunner.parseTemplate(uriTemp, uri, variables)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (apiTestRunner *APITestRunner) parseQuery(variables map[string]any, reqQuery string) (map[string]string, error) {
	if reqQuery == "" {
		return map[string]string{}, nil
	}

	if matchedErr := apiTestRunner.isExistVars(reqQuery, variables); matchedErr != nil {
		return nil, matchedErr
	}
	str, err := apiTestRunner.parseTemplate(queryTemp, reqQuery, variables)
	if err != nil {
		return nil, err
	}
	queryMap := make(map[string]string)
	err = json.Unmarshal([]byte(str), &queryMap)
	if err != nil {
		return nil, err
	}

	return queryMap, nil
}

func (apiTestRunner *APITestRunner) parseBody(variables map[string]any, reqBody string) (string, error) {
	if matchedErr := apiTestRunner.isExistVars(reqBody, variables); matchedErr != nil {
		return "", matchedErr
	}

	str, err := apiTestRunner.parseTemplate(bodyTemp, reqBody, variables)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (apiTestRunner *APITestRunner) parseTemplate(
	temp *template.Template,
	str string,
	variables map[string]any,
) (string, error) {
	if str == "" {
		return "", nil
	}
	tempParsed, parseErr := temp.Parse(str)
	if parseErr != nil {
		return "", parseErr
	}

	var buf bytes.Buffer
	if err := tempParsed.Execute(&buf, variables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (apiTestRunner *APITestRunner) apiTestWithMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	method string,
	params *apiTestRunnerParams,
	statusCode int,
) (*apitest.Result, error) {
	switch method {
	case http.MethodGet:
		return lo.ToPtr(apiTestRunner.apiTestWithGetMethod(
			t,
			apiTest,
			params,
			statusCode,
		)), nil
	case http.MethodPost:
		return lo.ToPtr(apiTestRunner.apiTestWithPostMethod(
			t,
			apiTest,
			params,
			statusCode,
		)), nil
	case http.MethodPatch:
		return lo.ToPtr(apiTestRunner.apiTestWithPatchMethod(
			t,
			apiTest,
			params,
			statusCode,
		)), nil
	case http.MethodPut:
		return lo.ToPtr(apiTestRunner.apiTestWithPutMethod(
			t,
			apiTest,
			params,
			statusCode,
		)), nil
	case http.MethodDelete:
		return lo.ToPtr(apiTestRunner.apiTestWithDeleteMethod(
			t,
			apiTest,
			params,
			statusCode,
		)), nil
	default:
		return nil, fmt.Errorf("unsupported method %s", method)
	}
}

func (apiTestRunner *APITestRunner) apiTestWithGetMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	params *apiTestRunnerParams,
	statusCode int,
) apitest.Result {
	resp := apiTest.
		Debug().
		Handler(apiutil.GetHandler()).
		Get(params.url).
		Headers(params.headerMap).
		QueryParams(params.queryMap).
		Expect(t).
		Status(statusCode)

	for _, f := range params.respAssestFuncs {
		resp.Assert(f)
	}
	return resp.End()
}

func (apiTestRunner *APITestRunner) apiTestWithPostMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	params *apiTestRunnerParams,
	statusCode int,
) apitest.Result {
	resp := apiTest.
		Debug().
		Handler(apiutil.GetHandler()).
		Post(params.url).
		Headers(params.headerMap).
		QueryParams(params.queryMap).
		JSON(params.body).
		Expect(t).
		Status(statusCode)

	for _, f := range params.respAssestFuncs {
		resp.Assert(f)
	}

	return resp.End()
}

func (apiTestRunner *APITestRunner) apiTestWithPatchMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	params *apiTestRunnerParams,
	statusCode int,
) apitest.Result {
	resp := apiTest.
		Debug().
		Handler(apiutil.GetHandler()).
		Patch(params.url).
		Headers(params.headerMap).
		QueryParams(params.queryMap).
		JSON(params.body).
		Expect(t).
		Status(statusCode)

	for _, f := range params.respAssestFuncs {
		resp.Assert(f)
	}
	return resp.End()
}

func (apiTestRunner *APITestRunner) apiTestWithPutMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	params *apiTestRunnerParams,
	statusCode int,
) apitest.Result {
	resp := apiTest.
		Debug().
		Handler(apiutil.GetHandler()).
		Put(params.url).
		Headers(params.headerMap).
		QueryParams(params.queryMap).
		JSON(params.body).
		Expect(t).
		Status(statusCode)

	for _, f := range params.respAssestFuncs {
		resp.Assert(f)
	}
	return resp.End()
}

func (apiTestRunner *APITestRunner) apiTestWithDeleteMethod(
	t ginkgo.GinkgoTInterface,
	apiTest *apitest.APITest,
	params *apiTestRunnerParams,
	statusCode int,
) apitest.Result {
	resp := apiTest.
		Debug().
		Handler(apiutil.GetHandler()).
		Delete(params.url).
		Headers(params.headerMap).
		QueryParams(params.queryMap).
		JSON(params.body).
		Expect(t).
		Status(statusCode)

	for _, f := range params.respAssestFuncs {
		resp.Assert(f)
	}
	return resp.End()
}

func (apiTestRunner *APITestRunner) isExistVars(
	str string,
	varsMap map[string]any,
) error {
	pattern := `\{\{\.(.*?)\}\}`
	re := regexp.MustCompile(pattern)

	matched := re.FindAllStringSubmatch(str, -1)
	for _, match := range matched {
		if len(match) < 1 {
			continue
		}

		_, ok := varsMap[match[1]]
		if !ok {
			return fmt.Errorf("%s is not defined, err: %w", match[1], apierror.ErrVarsNotFound)
		}
	}
	return nil
}
