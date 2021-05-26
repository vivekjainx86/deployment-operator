package controllers

import (
	"context"
	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	deploymentv1alpha1 "github.com/Ridecell/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Deployment routine", func() {
	It("Creates a deployment", func() {
		ctx := context.Background()
		myDeploy := deploymentv1alpha1.MyDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vivek-returns",
				Namespace: "default",
			},
			Spec: deploymentv1alpha1.MyDeploymentSpec{
				Image: "nginx",
			},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &MyDeploymentReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}

		Expect(k8sClient.Create(ctx, &myDeploy)).To(Succeed())

		Expect(reconcile.ensureDeployment(&myDeploy, "./templates/deployment.yaml")).To(BeNil())
		deployment := appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: myDeploy.Namespace, Name: myDeploy.Name + "-web"}, &deployment)).To(Succeed())
		//Expect(myDeployRes.Status.Image).To(Equal("nginx"))

	})
})
