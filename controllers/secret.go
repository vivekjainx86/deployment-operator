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
	"encoding/json"
	"github.com/pkg/errors"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *MyDeploymentReconciler) ensureSecret(instance *deploymentv1alpha1.MyDeployment, templateFile string) error {
	// create template object
	td := &TemplateData{
		Instance: instance,
		Extra: map[string]interface{}{
			"Password": "123456789",
			"Name":     "Shaktimaan",
		},
	}

	secret := &corev1.Secret{}
	err := td.buildObjectWithTemplate(templateFile, secret)
	if err != nil {
		return errors.Wrap(err, "secret.go: Unable to build yaml from template")
	}

	// Get object
	found := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secret.Name, Namespace: secret.Namespace}}

	_, err = controllerutil.CreateOrUpdate(context.TODO(), r.Client, found, func() error {
		// update the secret
		found.Labels = secret.Labels
		found.Annotations = secret.Annotations
		found.Type = secret.Type
		found.StringData = secret.StringData
		return controllerutil.SetControllerReference(instance, found, r.Scheme)
	})
	if err != nil {
		return errors.Wrap(err, "secret.go: Secret CreateOrUpdate failed")
	}

	secretBytes, err := json.Marshal(found.Data)
	if err != nil {
		return errors.Wrap(err, "secret.go: Unable to serialize secret data")
	}

	instance.Status.SecretHash = HashItem(secretBytes)
	if err := r.Status().Update(context.TODO(), instance); err != nil {
		return errors.Wrap(err, "secret.go: Unable to update status")
	}

	return nil
}
