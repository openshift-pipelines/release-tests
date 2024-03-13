/*
Copyright 2020 The Tekton Authors

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

package operator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
    resource "github.com/openshift-pipelines/release-tests/pkg/config"
	chainv1alpha "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

//"quay.io/openshift-pipeline/chainstest"
var repo string = os.Getenv("CHAINS_REPOSITORY")
var tag string = time.Now().Format("010206150405")
var public_key_path = resource.Path("testdata/chains/key/cosign.pub")

func EnsureTektonChainsExists(clients chainv1alpha.TektonChainInterface, names utils.ResourceNames) (*v1alpha1.TektonChain, error) {
    // If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
    ks, err := clients.Get(context.TODO(), names.TektonChain, metav1.GetOptions{})
    err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
        ks, err = clients.Get(context.TODO(), names.TektonChain, metav1.GetOptions{})
        if err != nil {
            if apierrs.IsNotFound(err) {
                log.Printf("Waiting for availability of chains cr [%s]\n", names.TektonChain)
                return false, nil
            }
            return false, err
        }
        return true, nil
    })
    return ks, err
}

func VerifySignature(resourceType string){
    //Get a signature of taskrun payload
    resourceUID := cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", "jsonpath='{.metadata.uid}'").Stdout()
    resourceUID = strings.Trim(resourceUID, "'")
    jsonpath := fmt.Sprintf("jsonpath=\"{.metadata.annotations.chains\\.tekton\\.dev/signature-%s-%s}\"", resourceType, resourceUID)
    fmt.Println("Waiting 30 seconds")
	cmd.MustSuccedIncreasedTimeout(time.Second*45 ,"sleep", "30")
    signature := cmd.MustSucceed("tkn", resourceType, "describe", "--last", "-o", jsonpath).Stdout()
    signature = strings.Trim(signature, "\"")
    //Decode the signature
    decodedSignature, err := base64.StdEncoding.DecodeString(signature)
        if err != nil {
            testsuit.T.Errorf("Error decoding base64")
        }
    //Create file with signature 
    file, err := os.Create("sign")
    if err != nil {
        testsuit.T.Errorf("Error creating file")
    }
    defer file.Close()
    _, err = file.WriteString(string(decodedSignature))
    if err != nil {
        testsuit.T.Errorf("Error writing to file")
    }
    //Verify signature with signing-secrets
    cmd.MustSucceed("cosign", "verify-blob-attestation", "--insecure-ignore-tlog", "--key", public_key_path, "--signature", "sign", "--type", "slsaprovenance", "--check-claims=false", "/dev/null")
}

func StartKanikoTask() {
    cmd.MustSucceed("oc", "secrets", "link", "pipeline", "quay", "--for=pull,mount")
    image := fmt.Sprintf("IMAGE=%s:%s", repo, tag)
    cmd.MustSucceed("tkn", "task", "start", "--param", image, "--use-param-defaults", "--workspace", "name=source,claimName=chains-pvc", "--workspace", "name=dockerconfig,secret=quay", "kaniko-chains")
    fmt.Println("Waiting 2 minutes for images to appear in image registry")
    cmd.MustSuccedIncreasedTimeout(time.Second*130 ,"sleep", "120")
}

func GetImageDigestedUrl() (string, string) {
    // Get Image digest
    var imageDigest string
    jsonOutput := cmd.MustSucceed("tkn", "tr", "describe", "--last", "-o", "json").Stdout()
    // Parse Json Output
    type TaskRun struct {
        Status struct {
            Results []struct {
                Name  string `json:"name"`
                Value string `json:"value"`
            } `json:"results"`
        } `json:"status"`
    }
    var taskrun TaskRun
    err := json.Unmarshal([]byte(jsonOutput), &taskrun)
    if err != nil {
        testsuit.T.Errorf("Error parsing Json output")
    }

    // Get IMAGE_DIGEST value
    for _, result := range taskrun.Status.Results {
        if strings.Contains(result.Name, "IMAGE_DIGEST"){
            imageDigest = strings.Split(result.Value, ":")[1]
        }
    }

    // Return image url with digest
    url := fmt.Sprintf("%s@sha256:%s", repo, imageDigest)
    return url, imageDigest
}

func VerifyImageSignature() {
    url, _ := GetImageDigestedUrl()
    cmd.MustSucceed("cosign", "verify", "--key", public_key_path, url)
}

func VerifyAttestation() {
    url, _ := GetImageDigestedUrl()
    cmd.MustSucceed("cosign", "verify-attestation", "--key", public_key_path, "--type", "slsaprovenance", url )
}

func CheckAttestation() {
    // Get UUID
    _, imageDigest := GetImageDigestedUrl()
    jsonOutput := cmd.MustSucceed("rekor-cli", "search", "--format", "json", "--sha", imageDigest ).Stdout()

    // Parse Json output to find UUID
    type UUID struct{
        UUIDs []string `json:"UUIDs"`
    }
    var uuid UUID 
    err := json.Unmarshal([]byte(jsonOutput), &uuid)
    if err != nil {
        testsuit.T.Errorf("Error parsing Json output")
    }
    rekor_uuid := uuid.UUIDs[0] 

    // Check the Attestation
    if strings.Contains(cmd.Run("rekor-cli", "get", "--uuid", rekor_uuid).Stdout(), "getLogEntryByUuidNotFound"){
        testsuit.T.Errorf("Failed to find Attestation")
    }
}

func CreateSigningSecretForTektonChains() {
	chainsPublicKey := os.Getenv("CHAINS_COSIGN_PUBLIC")
	chainsPrivateKey := os.Getenv("CHAINS_COSIGN_PRIVATE")
    os.Setenv("COSIGN_PASSWORD", "chainstest")
	chainsPassword := os.Getenv("COSIGN_PASSWORD")
	if chainsPublicKey != "" || chainsPrivateKey != "" {
		cmd.MustSucceed("oc", "create", "secret", "generic", "signing-secrets", "--from-literal=cosign.key="+chainsPrivateKey, "--from-literal=cosign.password="+chainsPassword, "--from-literal=cosign.pub="+chainsPublicKey, "--namespace", "openshift-pipelines")
	} else {
		cmd.MustSucceed("cosign", "generate-key-pair", "k8s://openshift-pipelines/signing-secrets")
	}
    cmd.MustSucceed("mv", "cosign.pub", public_key_path)
}
