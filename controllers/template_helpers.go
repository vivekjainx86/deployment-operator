/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"text/template"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
)

// MyDeploymentReconciler reconciles a MyDeployment object
type TemplateData struct {
	Instance *deploymentv1alpha1.MyDeployment
	Extra    map[string]interface{}
}

func (td *TemplateData) buildObjectWithTemplate(templateFile string) (map[string]interface{}, error) {
	// "Parse" parses a string into a template
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	// variable to store executed template
	var yamlDocument bytes.Buffer
	err = t.Execute(&yamlDocument, td)
	if err != nil {
		return nil, err
	}

	// convert yaml into Go object
	obj := make(map[string]interface{})
	err = yaml.Unmarshal(yamlDocument.Bytes(), &obj)
	if err != nil {
		return nil, err
	}

	// convert all nested objects to map[string]interface{}
	for key, value := range obj {
		obj[key] = cleanUpMapValue(value)
	}

	return obj, nil
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	default:
		return v
	}
}
