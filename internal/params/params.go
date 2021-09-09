package params

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/goformation/v4"
	"gopkg.in/yaml.v2"
)

type Parameters []struct {
	ParameterKey   string `yaml:"ParameterKey"`
	ParameterValue string `yaml:"ParameterValue"`
}

func CheckParams(path string) (bool, error) {
	template, err := goformation.Open(path)
	if err != nil {
		return false, err
	}

	if len(template.Parameters) == 0 {
		return false, nil
	}

	for key, _ := range template.Parameters {
		fmt.Println(key)
	}

	return false, nil
}

func MergeParameters(path string) ([]types.Parameter, error) {
	var params []types.Parameter
	var paramStruct Parameters
	paramFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = yaml.Unmarshal(paramFile, &paramStruct)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, val := range paramStruct {
		param := types.Parameter{
			ParameterKey:   &val.ParameterKey,
			ParameterValue: &val.ParameterValue,
		}
		params = append(params, param)
	}

	return params, nil
}

// 1. check if vars file is there if so run mergeParams
// 2. if no vars file, check if template contains params else run plan/apply
// 2.1. if contains params check if value/Default is present in all params
// 2.2 when no value is present ask user for input
