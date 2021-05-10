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
	"fmt"
	"reflect"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (cr *ComponentReconciler) deploymentComponent(templateFile string) (bool, error) {
	instance := cr.Instance.(*deploymentv1alpha1.MyDeployment)
	log := cr.Reconciler.Log.WithValues("namespace", instance.Namespace, "name", instance.Name)

	// create template object
	td := &TemplateData{
		Instance: instance,
		Extra:    nil,
	}

	yamlObject, err := td.buildObjectWithTemplate(templateFile)
	if err != nil {
		log.Info("Unable to build yaml from template")
		return false, err
	}

	// Get, create or update object
	fmt.Println("Converting object")
	deploy := &appsv1.Deployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(yamlObject, deploy)
	if err != nil {
		log.Info("Unable to convert yaml object from template")
		return false, err
	}

	// Get object
	found := &appsv1.Deployment{}
	if err := cr.Reconciler.Get(cr.Context, types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Setting controller reference")
			if err = controllerutil.SetControllerReference(instance, deploy, cr.Reconciler.Scheme); err != nil {
				log.Info("Error while setting controller reference")
				return false, err
			}
			fmt.Println("Creating object")
			if err = cr.Reconciler.Create(cr.Context, deploy); err != nil {
				log.Info("Error while creating object")
				return false, err
			}
			return true, nil
		}
		log.Info("Error while getting object")
		return false, err
	}

	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		found.Spec = deploy.Spec
		fmt.Println("Updating object")
		if err = cr.Reconciler.Update(cr.Context, found); err != nil {
			log.Info("Error while updating object")
			return false, err
		}
	}

	fmt.Println("Setting Status")
	instance.Status = found.Status
	if err := cr.Reconciler.Status().Update(cr.Context, instance); err != nil {
		log.Info("Unable to update MyDeployment status")
		return false, err
	}

	return false, nil
}
