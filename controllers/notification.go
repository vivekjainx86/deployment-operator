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
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func (r *MyDeploymentReconciler) ensureStartNotification(instance *deploymentv1alpha1.MyDeployment) error {

	// Check for image tag
	if instance.Spec.Image != instance.Status.Image {
		fmt.Println("Start notification sent.")
		return nil
	}

	secret := &corev1.Secret{}
	// Get object
	if err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, secret); err != nil {
		if client.IgnoreNotFound(err) == nil {
			fmt.Println("Start notification sent.")
			return nil
		}
		return errors.Wrap(err, "notification.go: Unable to get secret object")
	}

	secretBytes, err := json.Marshal(secret.Data)
	if err != nil {
		return errors.Wrap(err, "notification.go: Unable to serialize secret data")
	}

	if instance.Status.SecretHash != HashItem(secretBytes) {
		fmt.Println("Start notification sent.")
		return nil
	}

	return nil
}

func (r *MyDeploymentReconciler) ensureErrorNotification(instance *deploymentv1alpha1.MyDeployment) error {
	//- if instance.Status.Status is not error, then send error notification
	// this way we can send only 1 notification for error
	if instance.Status.Status != "Error" {
		fmt.Println("Error notification sent.")
	}
	return nil
}

func (r *MyDeploymentReconciler) ensureSuccessNotification(instance *deploymentv1alpha1.MyDeployment) error {
	// TODO:
	// Check for Instance.Status.Status
	//if instance.Status.Status != "Success" {
	//	fmt.Println("Success notification sent.")
	//}
	return nil
}
