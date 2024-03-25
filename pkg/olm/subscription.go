package olm

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	// Interval specifies the time between two polls.
	Interval = 10 * time.Second
	// Timeout specifies the timeout for the function PollImmediate to reach a certain status.
	Timeout            = 8 * time.Minute
	OperatorsNamespace = "openshift-operators"
	OLMNamespace       = "openshift-marketplace"
)

var (
	immediate             = int64(0)
	immediateDeleteOption = &metav1.DeleteOptions{GracePeriodSeconds: &immediate}
)

func SubscribeAndWaitForOperatorToBeReady(cs *clients.Clients, subscriptionName, channel, catalogsource string) (*v1alpha1.Subscription, error) {
	createSubscription(subscriptionName, channel, catalogsource)

	subs, err := WaitForSubscriptionState(cs, subscriptionName, OperatorsNamespace, IsSubscriptionInstalledCSVPresent)
	if err != nil {
		return nil, err
	}

	csvName := subs.Status.InstalledCSV
	_, err = WaitForClusterServiceVersionState(cs, csvName, OperatorsNamespace, IsCSVSucceeded)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func UptadeSubscriptionAndWaitForOperatorToBeReady(cs *clients.Clients, subscriptionName, channel string) (*v1alpha1.Subscription, error) {
	if _, err := UpdateSubscription(cs, subscriptionName, channel); err != nil {
		return nil, err
	}

	subs, err := WaitForSubscriptionState(cs, subscriptionName, OperatorsNamespace, IsSubscriptionInstalledCSVPresent)
	if err != nil {
		return nil, err
	}

	csvName := subs.Status.InstalledCSV

	_, err = WaitForClusterServiceVersionState(cs, csvName, OperatorsNamespace, IsCSVSucceeded)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func getSubcription(cs *clients.Clients, name string) *v1alpha1.Subscription {
	subscription, err := cs.OLM.OperatorsV1alpha1().Subscriptions(OperatorsNamespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get subscription %s in namespace %s \n %v", name, OperatorsNamespace, err)
	}
	return subscription
}

func createSubscription(name, channel, catalogsource string) {
	var subscription = struct {
		OperatorNamespace string
		SourceNamespace   string
		Channel           string
		SubscriptionName  string
		CatalogSource     string
	}{
		OperatorNamespace: OperatorsNamespace,
		SourceNamespace:   OLMNamespace,
		Channel:           channel,
		SubscriptionName:  name,
		CatalogSource:     catalogsource,
	}

	if _, err := config.TempDir(); err != nil {
		testsuit.T.Fail(err)
	}
	defer config.RemoveTempDir()

	tmpl, err := config.Read("subscription.yaml.tmp")
	if err != nil {
		testsuit.T.Fail(err)
	}

	sub, err := template.New("subscription").Parse(string(tmpl))
	if err != nil {
		testsuit.T.Fail(err)
	}

	var buffer bytes.Buffer
	if err = sub.Execute(&buffer, subscription); err != nil {
		testsuit.T.Fail(err)
	}
	file, err := config.TempFile("subscription.yaml")
	if err != nil {
		testsuit.T.Fail(err)
	}
	if err = os.WriteFile(file, buffer.Bytes(), 0666); err != nil {
		testsuit.T.Fail(err)
	}

	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", file).Stdout())
}

// OperatorCleanup deletes All related CSVs, subscription & installplan from cluster
func OperatorCleanup(cs *clients.Clients, name string) {
	sub := getSubcription(cs, name)

	//Delete CSV
	err := cs.OLM.OperatorsV1alpha1().ClusterServiceVersions(OperatorsNamespace).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to delete CSVs in namespace %s \n %v", OperatorsNamespace, err)
	}

	log.Printf("Output %s \n", cmd.MustSucceed(
		"oc", "delete", "--ignore-not-found", "-n", OperatorsNamespace,
		"subscription", sub.Name,
	).Stdout())
}

func UpdateSubscription(cs *clients.Clients, name, channel string) (*v1alpha1.Subscription, error) {
	subscription := getSubcription(cs, name)
	subscription.Spec.Channel = channel
	subs, err := cs.OLM.OperatorsV1alpha1().Subscriptions(OperatorsNamespace).Update(context.Background(), subscription, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func WaitForSubscriptionState(cs *clients.Clients, name, namespace string, inState func(s *v1alpha1.Subscription, err error) (bool, error)) (*v1alpha1.Subscription, error) {
	var lastState *v1alpha1.Subscription
	var err error
	waitErr := wait.PollImmediate(Interval, Timeout, func() (bool, error) {
		lastState, err = cs.OLM.OperatorsV1alpha1().Subscriptions(namespace).Get(context.Background(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, errors.Wrapf(waitErr, "subscription %s is not in desired state, got: %+v", name, lastState)
	}
	return lastState, nil
}

func WaitForClusterServiceVersionState(cs *clients.Clients, name, namespace string, inState func(s *v1alpha1.ClusterServiceVersion, err error) (bool, error)) (*v1alpha1.ClusterServiceVersion, error) {
	var lastState *v1alpha1.ClusterServiceVersion
	var err error
	waitErr := wait.PollImmediate(Interval, Timeout, func() (bool, error) {
		lastState, err = cs.OLM.OperatorsV1alpha1().ClusterServiceVersions(namespace).Get(context.Background(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, errors.Wrapf(waitErr, "clusterserviceversion %s is not in desired state, got: %+v", name, lastState)
	}
	return lastState, nil
}

func IsCSVSucceeded(c *v1alpha1.ClusterServiceVersion, err error) (bool, error) {
	return c.Status.Phase == "Succeeded", err
}

func IsSubscriptionInstalledCSVPresent(s *v1alpha1.Subscription, err error) (bool, error) {
	return s.Status.InstalledCSV != "" && s.Status.InstalledCSV != "<none>", err
}
