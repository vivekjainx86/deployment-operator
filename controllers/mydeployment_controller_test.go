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
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("MyDeployment controller", func() {
	It("Creates a MyDeployment object", func() {
		ctx := context.Background()
		myDeploy := deploymentv1alpha1.MyDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vivek",
				Namespace: "default",
			},
			Spec: deploymentv1alpha1.MyDeploymentSpec{},
		}

		zapLog, _ := zap.NewDevelopment()
		reconcile := &MyDeploymentReconciler{
			Client: k8sClient,
			Log:    zapr.NewLogger(zapLog),
			Scheme: scheme.Scheme,
		}
		BaseTemplatePath = "./templates"

		Expect(k8sClient.Create(ctx, &myDeploy)).To(Succeed())

		myDeployName := types.NamespacedName{Namespace: myDeploy.Namespace, Name: myDeploy.Name}
		Expect(reconcile.Reconcile(context.TODO(), ctrl.Request{NamespacedName: myDeployName})).To(Equal(ctrl.Result{Requeue: false}))
		myDeployRes := deploymentv1alpha1.MyDeployment{}
		Expect(k8sClient.Get(ctx, myDeployName, &myDeployRes)).To(Succeed())
		Expect(myDeployRes.Status.Image).To(Equal("nginx"))
		deployment := appsv1.Deployment{}
		myDeployName.Name = myDeploy.Name + "-web"
		Expect(k8sClient.Get(ctx, myDeployName, &deployment)).To(Succeed())

	})
})
