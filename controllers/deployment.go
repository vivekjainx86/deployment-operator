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
	"context"
	"github.com/pkg/errors"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *MyDeploymentReconciler) ensureDeployment(instance *deploymentv1alpha1.MyDeployment, templateFile string) error {
	// create template object
	td := &TemplateData{
		Instance: instance,
		Extra: map[string]interface{}{
			"HTTPS": "true",
		},
	}

	deploy := &appsv1.Deployment{}
	err := td.buildObjectWithTemplate(templateFile, deploy)
	if err != nil {
		return errors.Wrap(err, "deployment.go: Unable to build yaml from template")
	}

	// Get object
	found := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: deploy.Name, Namespace: deploy.Namespace}}

	_, err = controllerutil.CreateOrPatch(context.TODO(), r.Client, found, func() error {
		// update the Deployment object
		found.Labels = deploy.Labels
		if found.Annotations == nil {
			found.Annotations = map[string]string{}
		}
		for k, v := range deploy.Annotations {
			found.Annotations[k] = v
		}
		found.Spec = deploy.Spec
		return controllerutil.SetControllerReference(instance, found, r.Scheme)
	})
	if err != nil {
		return errors.Wrap(err, "deployment.go: Deployment CreateOrUpdate failed")
	}

	instance.Status.Image = instance.Spec.Image
	if err := r.Status().Update(context.TODO(), instance); err != nil {
		return errors.Wrap(err, "deployment.go: Unable to update status")
	}

	return nil
}
