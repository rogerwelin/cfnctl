package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	gocfn "github.com/awslabs/goformation/v4/cloudformation"
	"github.com/rogerwelin/cfnctl/pkg/client"
	"gopkg.in/yaml.v2"
)

type Parameters []struct {
	ParameterKey   string `yaml:"ParameterKey"`
	ParameterValue string `yaml:"ParameterValue"`
}

func mergeParameters(inParams gocfn.Parameters) ([]types.Parameter, error) {
	var params []types.Parameter
	var paramStruct Parameters
	paramFile, err := ioutil.ReadFile("vars.yaml")
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

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	svc := cloudformation.NewFromConfig(cfg)

	ctl := &client.Cfnctl{
		Svc:           svc,
		StackName:     "dynamolambda",
		ChangesetName: "dynamolambda",
	}

	/*
		err = cfnCtlPlan(ctl)
		if err != nil {
			panic(err)
		}
	*/

	err = cfnCtlApply(ctl)
	if err != nil {
		panic(err)
	}

}
