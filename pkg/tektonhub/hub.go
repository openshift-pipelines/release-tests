package tektonhub

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"

	"net/http"
	"os"
	"runtime"
	"text/template"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"gopkg.in/yaml.v2"
)

func CreateHubCR(cs *clients.Clients) {
	var hub = struct {
		TargetNamespace string
		CatalogName     string
		Org             string
		Type            string
		Provider        string
		URL             string
		Revision        string
		RefreshInterval string
	}{
		TargetNamespace: config.TargetNamespace,
		CatalogName:     "manoj",
		Org:             "tektoncd",
		Type:            "github",
		Provider:        "github",
		URL:             "https://github.com/tektoncd/catalog",
		Revision:        "main",
		RefreshInterval: "2m",
	}

	if _, err := config.TempDir(); err != nil {
		assert.FailOnError(err)
	}
	defer config.RemoveTempDir()

	tmpl, err := config.Read("tektonhub.yaml.tmp")
	if err != nil {
		assert.FailOnError(err)
	}

	sub, err := template.New("hub").Parse(string(tmpl))
	if err != nil {
		assert.FailOnError(err)
	}

	var buffer bytes.Buffer
	if err = sub.Execute(&buffer, hub); err != nil {
		assert.FailOnError(err)
	}
	file, err := config.TempFile("hub.yaml")
	assert.FailOnError(err)
	if err = os.WriteFile(file, buffer.Bytes(), 0666); err != nil {
		assert.FailOnError(err)
	}

	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", file).Stdout())
}

func GetTektonHubElements() {
	chromedriverURL := getChromedriverURL()
	downloadFile(chromedriverURL, "chromedriver.zip")
	extractFile("chromedriver.zip", "chromedriver")
	UpdateYAMLAndCommit(
		"https://github.com/manojbison/catalog.git",
		"manojbison",
<<<<<<< HEAD
		"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",
=======
		"",
>>>>>>> origin/master
		"task/git-clone/0.9/git-clone.yaml",
		"git clone test",
		"manojbison",
		"mbison@redhat.com",
	)
	time.Sleep(5 * time.Second)
	TektonHubElements()
}

func getChromedriverURL() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	switch osName {
	case "windows":
		osName = "win"
	case "darwin":
		osName = "mac"
	case "linux":
		osName = "linux"
	}

	switch arch {
	case "amd64":
		arch = "64"
	}

	return fmt.Sprintf("https://chromedriver.storage.googleapis.com/111.0.5563.64/chromedriver_%s%s.zip", osName, arch)
}

func downloadFile(url string, filepath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
}

func extractFile(zipFileName string, extractToDir string) error {
	// Open the zip file for reading.
	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	// Create the extractToDir directory if it doesn't exist.
	if _, err := os.Stat(extractToDir); os.IsNotExist(err) {
		os.Mkdir(extractToDir, os.ModePerm)
	}

	// Iterate through each file in the zip archive.
	for _, f := range r.File {
		// Open the file inside the zip archive.
		rc, err := f.Open()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer rc.Close()

		// Create the file in the extractToDir directory.
		path := fmt.Sprintf("%s/%s", extractToDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			outFile, err := os.Create(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer outFile.Close()

			// Copy the contents of the file from the zip archive to the output file.
			_, err = io.Copy(outFile, rc)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	return nil
}

func TektonHubElements() (map[string]selenium.WebElement, error) {
	// Start a Selenium WebDriver serverx
	apiURL, uiURL, _ := operator.VerifyTektonHubURLs(store.Clients())
	fmt.Println(uiURL, apiURL)
	driverPath := "chromedriver/chromedriver"
	err := os.Chmod(driverPath, 0755)
	if err != nil {
		return nil, err
	}
	service, err := selenium.NewChromeDriverService(driverPath, 9515)
	if err != nil {
		return nil, fmt.Errorf("failed to start the WebDriver server: %w", err)
	}
	defer service.Stop()

	// Set up Chrome options
	opts := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--disable-gpu",
			"--ignore-certificate-errors",
			"--no-sandbox",
			"--allow-insecure-localhost",
			"--allow-running-insecure-content",
			"--ignore-urlfetcher-cert-requests",
			"--ssl-version-min=tls1.2",
		},
	}

	// Connect to Chrome WebDriver instance
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(opts)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the WebDriver: %w", err)
	}

	// Send HTTPS request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send HTTPS request with custom client
	resp, err := client.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print response status and body
	fmt.Println("Response status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response body:", string(body))

	// Wait for the page to load
	time.Sleep(5 * time.Second)

	if err := wd.Get(uiURL); err != nil {
		return nil, fmt.Errorf("error navigating to Tekton Hub website: %w", err)
	}

	// Wait for the page to load
	time.Sleep(5 * time.Second)

	// Find the "kind", "platform", "catalog", "taskname" and "category" elements on the page
	kindElement, err := wd.FindElement(selenium.ByXPATH, "//h1[@class='hub-filter-header' and text()='Kind']")
	if err != nil {
		return nil, fmt.Errorf("error finding 'kind' element: %w", err)
	}
	platformElement, err := wd.FindElement(selenium.ByXPATH, "//h1[@class='hub-filter-header' and text()='Platform']")
	if err != nil {
		return nil, fmt.Errorf("error finding 'platform' element: %w", err)
	}
	catalogElement, err := wd.FindElement(selenium.ByXPATH, "//h1[@class='hub-filter-header' and text()='Catalog']")
	if err != nil {
		return nil, fmt.Errorf("error finding 'catalog' element: %w", err)
	}
	categoryElement, err := wd.FindElement(selenium.ByXPATH, "//h1[@class='hub-filter-header' and text()='Category']")
	if err != nil {
		return nil, fmt.Errorf("error finding 'category' element: %w", err)
	}

	err = wd.Get(uiURL + "/tekton/task/git-clone")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Wait for the page to load
	time.Sleep(5 * time.Second)

	taskname, err := wd.FindElement(selenium.ByXPATH, "//h1[@class='hub-details-resource-name' and text()='git clone test']")
	if err != nil {
		return nil, fmt.Errorf("error finding 'kind' element: %w", err)
	}

	fmt.Println(taskname)
	// Return a map of the elements
	return map[string]selenium.WebElement{
		"kind":     kindElement,
		"platform": platformElement,
		"catalog":  catalogElement,
		"category": categoryElement,
		"taskname": taskname,
	}, nil

}

func UpdateYAMLAndCommit(repoURL, username, password, Yamlpath, displayName, authorName, authorEmail string) error {
	// Set authentication credentials
	auth := &ghttp.BasicAuth{
		Username: username,
		Password: password,
	}

	// Clone the repository
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := git.PlainClone(cwd, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     auth,
	})

	fmt.Println("cloned")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to clone repository: %s", err)
	}

	// Read YAML file into a byte array
	yamlFile, err := ioutil.ReadFile(path.Join(cwd, Yamlpath))
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %s", err)
	}

	// Create a map to hold the YAML data
	data := make(map[interface{}]interface{})

	// Unmarshal the YAML data into the map
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML data: %s", err)
	}

	// Modify the desired string value
	data["metadata"].(map[interface{}]interface{})["annotations"].(map[interface{}]interface{})["tekton.dev/displayName"] = displayName

	// Marshal the updated YAML data back to a byte array
	updatedYaml, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal updated YAML data: %s", err)
	}

	// Write the updated YAML back to the file
	err = ioutil.WriteFile(path.Join(cwd, Yamlpath), updatedYaml, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated YAML data to file: %s", err)
	}

	fmt.Println("YAML file updated successfully!")

	// Open the repository
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return fmt.Errorf("failed to open repository: %s", err)
	}

	// Add the changes
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %s", err)
	}

	_, err = w.Add(Yamlpath)
	if err != nil {
		return fmt.Errorf("failed to add changes: %s", err)
	}

	status, _ := w.Status()

	fmt.Println(status)

	// Commit the changes
	commitMsg := fmt.Sprintf("Update %s", Yamlpath)
	commitDate := time.Now()

	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  commitDate,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %s", err)
	}

	// Push the changes
	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("failed to push changes: %s", err)
	}

	// Print the commit hash
	fmt.Println(commit)

	return nil
}
