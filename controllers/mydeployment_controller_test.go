package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("MyDeployment controller", func() {

	BeforeEach(func() {
		instance = &deploymentv1alpha1.MyDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "",
				Namespace: "default",
			},
			Spec: deploymentv1alpha1.MyDeploymentSpec{
				Image: "myimagetag",
			},
		}
	})

	It("Creates a MyDeployment object", func() {
		ctx := context.Background()
		instance.Name = "mydeployment-test-1"

		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		instanceName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}
		Expect(reconcile.Reconcile(ctx, ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{Requeue: false}))

		instanceRes := &deploymentv1alpha1.MyDeployment{}
		Expect(k8sClient.Get(ctx, instanceName, instanceRes)).To(Succeed())
		Expect(instanceRes.Status.Image).To(Equal("myimagetag"))

		deployment := &appsv1.Deployment{}
		instanceName.Name = instance.Name + "-web"
		Expect(k8sClient.Get(ctx, instanceName, deployment)).To(Succeed())

	})

	It("Creates a MyDeployment object with default values", func() {
		ctx := context.Background()
		instance.Name = "mydeployment-test-2"
		instance.Spec = deploymentv1alpha1.MyDeploymentSpec{}

		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		instanceName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name}
		Expect(reconcile.Reconcile(ctx, ctrl.Request{NamespacedName: instanceName})).To(Equal(ctrl.Result{Requeue: false}))

		instanceRes := &deploymentv1alpha1.MyDeployment{}
		Expect(k8sClient.Get(ctx, instanceName, instanceRes)).To(Succeed())
		Expect(instanceRes.Status.Image).To(Equal("nginx"))

	})

})
