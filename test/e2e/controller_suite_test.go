package e2e

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/combo/api/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	kubeRestClient *rest.RESTClient
	kubeClientSet  *kubernetes.Clientset
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

	// Load the kube config to construct a client
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Build a client that can query the combo API's
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		Fail("Error while building kubeconfig for suite:" + err.Error())
	}

	// Build a client set before editing the config for the rest client
	if kubeClientSet, err = kubernetes.NewForConfig(config); err != nil {
		Fail("Error while building kube clientset for suite: " + err.Error())
	}

	// Add customization to config for interacting with combo APIs
	v1alpha1.AddToScheme(scheme.Scheme)

	config.APIPath = "/apis"
	config.ContentConfig.GroupVersion = &v1alpha1.GroupVersion
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	// Create the kube client
	if kubeRestClient, err = rest.UnversionedRESTClientFor(config); err != nil {
		Fail("Error while building kube rest client for suite:" + err.Error())
	}

	// Confirm that the combination and template resources exist on cluster before starting suite
	var combinationList v1alpha1.CombinationList
	Expect(kubeRestClient.Get().Resource("combinations").Do(ctx).Into(&combinationList)).To(Succeed())
	Expect(combinationList.APIVersion).To(Not(BeZero()))

	var templateList v1alpha1.TemplateList
	Expect(kubeRestClient.Get().Resource("templates").Do(ctx).Into(&templateList)).To(Succeed())
	Expect(templateList.APIVersion).To(Not(BeZero()))

	// Confirm that the combo operator is running and healthy
	Eventually(func() (bool, error) {
		deployment, err := kubeClientSet.AppsV1().Deployments("combo").Get(ctx, "combo-operator", v1.GetOptions{})
		return deployment.Status.AvailableReplicas > 0, err
	}).Should(BeTrue())

})
