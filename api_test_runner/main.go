package apitestrunner

import (
	"errors"
	"fmt"

	testCfg "github.com/RyanTokManMokMTM/api-testing-go/config"
	apihelper "github.com/RyanTokManMokMTM/api-testing-go/utils/helper"
	apiutil "github.com/RyanTokManMokMTM/api-testing-go/utils/util"

	"github.com/onsi/ginkgo/v2"

	"github.com/steinfletcher/apitest"
)

type APITestRunner struct {
	config *testCfg.APITest
}

func NewAPITestRuuner(config *testCfg.APITest) *APITestRunner {
	return &APITestRunner{
		config: config,
	}
}

func (apiTestRunner *APITestRunner) Start(
	t ginkgo.GinkgoTInterface,
	cfg *testCfg.APITest,
	workflow testCfg.Workflow,
	respData any,
	respCache *map[string]interface{},
) (*apitest.Result, error) {
	variables := make(map[string]any)

	// Set variables
	for _, v := range workflow.Request.Vars {
		variables[v.Name] = cfg.Config[v.Name] // get value from config
	}

	// If there is variable in request from the privous response.
	for _, v := range workflow.Request.FromResponse {
		if respCache == nil {
			return nil, errors.New("response cache is nil")
		}
		var data any
		var err error

		if len(v.Step) == 0 {
			data, err = apihelper.GetRespFieldData(v.FromField, respData)
		} else {
			data, err = apiTestRunner.getValueFromResponse(
				respCache,
				v.Name,
				v.Step,
				v.FromField)
		}

		if err != nil {
			return nil, err
		}

		variables[v.Name] = data // get value from response
	}

	params, err := apiTestRunner.parseRequestData(variables, workflow)
	if err != nil {
		return nil, err
	}

	apitest := apitest.New(fmt.Sprintf("Step - %s", workflow.Step))
	if apiTestRunner.config.NetworkEnable {
		apitest = apitest.EnableNetworking()
		params.url = fmt.Sprintf("%s%s", apiTestRunner.config.Host, params.url)
	} else {
		apitest = apitest.Handler(apiutil.GetHandler())
	}

	return apiTestRunner.apiTestWithMethod(
		t,
		apitest,
		workflow.Request.Method,
		params,
		workflow.ExpectResponse.StatusCode,
	)
}
