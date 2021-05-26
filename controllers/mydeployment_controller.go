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
	"crypto/sha1"
	"encoding/hex"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/predicate"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	// This variable is overridden by controller tests.
	// This is required because template folder path is relatively different for main.go and test suite.
	BaseTemplatePath = "./controllers/templates"
)

// MyDeploymentReconciler reconciles a MyDeployment object
type MyDeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=deployment.ridecell.io,resources=mydeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=deployment.ridecell.io,resources=mydeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=deployment.ridecell.io,resources=mydeployments/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MyDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mydeployment", req.NamespacedName)
	log.Info("Reconciling object")
	// Get MyDeployment Object
	instance := &deploymentv1alpha1.MyDeployment{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Info("Unable to fetch MyDeployment object from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// initialize instance with default values
	r.ensureDefaults(instance)

	// Send start notification
	err := r.ensureStartNotification(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Call deployment sub-routine
	err = r.ensureDeployment(instance, BaseTemplatePath+"/deployment.yaml")
	if err != nil {
		return r.manageError(instance, err, false)
	}
	// Call secret sub-routine
	err = r.ensureSecret(instance, BaseTemplatePath+"/secret.yaml")
	if err != nil {
		return r.manageError(instance, err, false)
	}
	return r.manageSuccess(instance, false)
}

func (r *MyDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//pred := predicate.ResourceVersionChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&deploymentv1alpha1.MyDeployment{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Secret{}).
		//WithEventFilter(pred).
		Complete(r)
}

func (r *MyDeploymentReconciler) manageError(instance *deploymentv1alpha1.MyDeployment, err error, requeue bool) (ctrl.Result, error) {
	// Log the error
	log := r.Log.WithValues("mydeployment", instance.Namespace+"/"+instance.Name)
	log.Error(err, "Updating error status")

	_ = r.ensureErrorNotification(instance)

	// Set error status
	instance.Status.Status = "Error"
	instance.Status.Message = err.Error()
	er := r.Status().Update(context.TODO(), instance)
	if er != nil {
		log.Error(er, "Unable to update error status")
		return ctrl.Result{Requeue: requeue}, er
	}
	return ctrl.Result{Requeue: requeue}, err
}

func (r *MyDeploymentReconciler) manageSuccess(instance *deploymentv1alpha1.MyDeployment, requeue bool) (ctrl.Result, error) {
	log := r.Log.WithValues("mydeployment", instance.Namespace+"/"+instance.Name)

	err := r.ensureSuccessNotification(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Set success status
	instance.Status.Status = "Success"
	instance.Status.Message = "Reconcile completed"
	er := r.Status().Update(context.TODO(), instance)
	if er != nil {
		log.Error(er, "Unable to update success status")
		return ctrl.Result{Requeue: requeue}, er
	} else {
		log.Info("Update status completed")
	}
	return ctrl.Result{Requeue: requeue}, nil
}

func HashItem(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
