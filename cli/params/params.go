package params

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/goformation/v4"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

type parameters []struct {
	ParameterKey   string `yaml:"ParameterKey"`
	ParameterValue string `yaml:"ParameterValue"`
}

// BuildInputParams builds CF parameter struct from given input
func BuildInputParams(params []string) ([]types.Parameter, error) {
	res := make(map[string]string)
	var cfParams []types.Parameter
	whiteBold := color.New(color.Bold).SprintfFunc()

	if len(params) > 1 {
		fmt.Printf("%s\n\n", whiteBold("Enter parameter values:"))
	} else {
		fmt.Printf("%s\n\n", whiteBold("Enter parameter value:"))
	}

	for _, val := range params {
		p := promptui.Prompt{
			Label: val,
		}
		result, err := p.Run()
		if err != nil {
			return nil, err
		}
		res[val] = result
	}

	for key, val := range res {
		param := types.Parameter{
			ParameterKey:   &key,
			ParameterValue: &val,
		}
		cfParams = append(cfParams, param)
	}

	return cfParams, nil
}

// CheckInputParams checks if the given CF template contains parameters or not
func CheckInputParams(path string) (bool, []string, error) {
	var params []string

	template, err := goformation.Open(path)
	if err != nil {
		return false, nil, fmt.Errorf("could not open template, %w", err)
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

// MergeFileParams reads parameters from a separate file and returns CF API parameter type
func MergeFileParams(path string) ([]types.Parameter, error) {
	var params []types.Parameter
	var paramStruct parameters
	paramFile, err := os.ReadFile(path)
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
