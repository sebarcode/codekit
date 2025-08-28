package codekit

import (
	"gopkg.in/yaml.v3"
)

func YamlToM(text string) (M, error) {
	var res map[string]interface{}
	if err := yaml.Unmarshal([]byte(text), &res); err != nil {
		return nil, err
	}
	return M(res), nil
}

func YamlToObj(text string, obj interface{}) error {
	return yaml.Unmarshal([]byte(text), obj)
}

func ObjToYaml(obj interface{}) (string, error) {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
