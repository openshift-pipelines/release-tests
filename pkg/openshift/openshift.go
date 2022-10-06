package openshift

import (
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	imageStream "github.com/openshift/client-go/image/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetImageStreamTags(c *clients.Clients, namespace, name string) []string {
	fmt.Printf("Getting imagestream from the namespace %s", namespace)
	is := imageStream.NewForConfigOrDie(c.KubeConfig)
	isRequired, err := is.ImageV1().ImageStreams(namespace).Get(c.Ctx, name, v1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	tags := isRequired.Spec.Tags
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}
