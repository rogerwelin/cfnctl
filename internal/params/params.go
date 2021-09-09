package params

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/goformation/v4"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

type Parameters []struct {
	ParameterKey   string `yaml:"ParameterKey"`
	ParameterValue string `yaml:"ParameterValue"`
}

func BuildInputParams(params []string) {
	res := []string{}

	fmt.Printf("Enter parameter value/s\n\n")

	for _, val := range params {
		p := promptui.Prompt{
			Label: val,
		}
		result, err := p.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		res = append(res, result)
	}
	fmt.Println(res)
}

func CheckInputParams(path string) (bool, []string, error) {
	var params []string

	template, err := goformation.Open(path)
	if err != nil {
		return false, nil, err
	}

	if len(template.Parameters) == 0 {
		return false, nil, nil
	}

	for key, val := range template.Parameters {
		if val.Default == nil {
			params = append(params, key)
		}
	}

	return true, params, nil
}

func MergeFileParams(path string) ([]types.Parameter, error) {
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
