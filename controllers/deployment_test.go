package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Deployment routine", func() {

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

	It("Creates a deployment", func() {
		ctx := context.Background()
		instance.Name = "deployment-test-1"

		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		Expect(reconcile.ensureDeployment(instance, "./templates/deployment.yaml")).To(BeNil())
		deployment := &appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-web"}, deployment)).To(Succeed())
		//Expect(instanceRes.Status.Image).To(Equal("nginx"))

	})

	It("Creates a deployment second", func() {
		ctx := context.Background()
		instance.Name = "deployment-test-2"

		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		Expect(reconcile.ensureDeployment(instance, "./templates/deployment.yaml")).To(BeNil())
		deployment := &appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-web"}, deployment)).To(Succeed())
		//Expect(instanceRes.Status.Image).To(Equal("nginx"))

	})
})
