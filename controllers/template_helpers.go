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
	"sigs.k8s.io/yaml"
	"text/template"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
)

type TemplateData struct {
	Instance *deploymentv1alpha1.MyDeployment
	Extra    map[string]interface{}
}

func (td *TemplateData) buildObjectWithTemplate(templateFile string, object interface{}) error {
	// "Parse" parses a string into a template
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	// create yaml using TemplateData
	var yamlDocument bytes.Buffer
	err = t.Execute(&yamlDocument, td)
	if err != nil {
		return err
	}

	// convert yaml into Go object
	err = yaml.Unmarshal(yamlDocument.Bytes(), object)
	if err != nil {
		return err
	}

	return nil
}
