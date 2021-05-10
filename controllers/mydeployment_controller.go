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
	"fmt"
	//"reflect"

	"github.com/go-logr/logr"
	//"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	//"sigs.k8s.io/controller-runtime/pkg/predicate"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
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

	// build the template
	// td := &TemplateData{
	// 	Instance: instance,
	// 	Extra:    nil,
	// }
	//
	// yamlObject, err := td.buildObjectWithTemplate("./controllers/templates/deployment.yaml")
	// if err != nil {
	// 	log.Info("Unable to build yaml from template")
	// 	return ctrl.Result{}, err
	// }

	// fmt.Println("Converting object")
	// deploy := &appsv1.Deployment{}
	// err = runtime.DefaultUnstructuredConverter.FromUnstructured(yamlObject, deploy)
	// if err != nil {
	// 	log.Info("Unable to convert yaml object from template")
	// 	return ctrl.Result{}, err
	// }

	// Get object
	// found := &appsv1.Deployment{}
	// if err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found); err != nil {
	// 	if errors.IsNotFound(err) {
	// 		fmt.Println("Setting controller reference")
	// 		if err = controllerutil.SetControllerReference(instance, deploy, r.Scheme); err != nil {
	// 			log.Info("Error while setting controller reference")
	// 			return ctrl.Result{}, err
	// 		}
	// 		fmt.Println("Creating object")
	// 		if err = r.Create(ctx, deploy); err != nil {
	// 			log.Info("Error while creating object")
	// 			return ctrl.Result{}, err
	// 		}
	// 		return ctrl.Result{Requeue: true}, nil
	// 	}
	// 	log.Info("Error while getting object")
	// 	return ctrl.Result{}, err
	// }
	//
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	fmt.Println("Updating object")
	// 	if err = r.Update(ctx, found); err != nil {
	// 		log.Info("Error while updating object")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	cr := &ComponentReconciler{
		Context:    ctx,
		Reconciler: r,
		Instance:   instance,
	}

	// Call deploymentComponent method
	requeue, err := cr.deploymentComponent("./controllers/templates/deployment.yaml")
	if err != nil {
		return ctrl.Result{}, err
	} else if requeue {
		return ctrl.Result{Requeue: true}, nil
	}

	// fmt.Println("Setting Status")
	// instance.Status = found.Status
	// if err := r.Status().Update(ctx, instance); err != nil {
	// 	log.Info("Unable to update MyDeployment status")
	// 	return ctrl.Result{}, err
	// }
	fmt.Println("End reconcile")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//pred := predicate.ResourceVersionChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&deploymentv1alpha1.MyDeployment{}).
		Owns(&appsv1.Deployment{}).
		//WithEventFilter(pred).
		Complete(r)
}
