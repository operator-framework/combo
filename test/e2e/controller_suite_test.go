package e2e

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/combo/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	kubeclient client.Client
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	SetDefaultEventuallyTimeout(1 * time.Minute)
	SetDefaultEventuallyPollingInterval(1 * time.Second)

	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheme := runtime.NewSchemeBuilder(
		kscheme.AddToScheme,
		v1alpha1.AddToScheme,
	)

	config := ctrl.GetConfigOrDie()

	var err error
	kubeclient, err = client.New(config, client.Options{})
	if err != nil {
		Fail("Error while building client: " + err.Error())
	}

	err = scheme.AddToScheme(kubeclient.Scheme())
	if err != nil {
		Fail("Error while add schemes to client: " + err.Error())
	}

	Context("should already have the combination CRD defined", func() {
		Eventually(func() (bool, error) {
			combinationCRD, err := kubeclient.RESTMapper().ResourceFor(v1alpha1.GroupVersion.WithResource("combination"))
			return combinationCRD.Empty(), err
		}).ShouldNot(BeTrue())
	})

	Context("should already have the template CRD defined", func() {
		Eventually(func() (bool, error) {
			templateCRD, err := kubeclient.RESTMapper().ResourceFor(v1alpha1.GroupVersion.WithResource("template"))
			return templateCRD.Empty(), err
		}).ShouldNot(BeTrue())
	})

	Context("should already have combo operator running and healthy", func() {
		Eventually(func() (int, error) {
			deployment := appsv1.Deployment{}
			deploymentNamespace := types.NamespacedName{
				Name:      "combo-operator",
				Namespace: "combo",
			}

			err = kubeclient.Get(ctx, deploymentNamespace, &deployment)

			return int(deployment.Status.AvailableReplicas), err
		}).Should(BeNumerically(">", 0))
	})
})
