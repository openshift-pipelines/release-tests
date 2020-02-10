package pipelines

//var _ = gauge.Step("Create pipeline from <path_to_pipeline_yaml>", func(path_to_pipeline_yaml string) {
//// TODO: Fix Tkn().Path instaed Tkn.Cmd()
//cmd.AssertOutput(&cmd.Cmd{
//Args: []string{
//store.Tkn().Path, "pipeline", "create",
//"--from", filepath.Join(helper.RootDir(), path_to_pipeline_yaml),
//"-n", store.Namespace()},
//Expected: icmd.Expected{
//ExitCode: 0,
//Err:      icmd.None,
//Out:      "Pipeline created: test-pipeline\n",
//},
//})
//})

//var _ = gauge.Step("Create pipeline file <path_to_pipeline_yaml> - In Non-existance namespace", func(path_to_pipeline_yaml string) {
//cmd.AssertOutput(&cmd.Cmd{
//Args: []string{store.Tkn().Path, "pipeline", "create", "--from", filepath.Join(helper.RootDir(), path_to_pipeline_yaml), "-n", "non-existance"},
//Expected: icmd.Expected{
//ExitCode: 1,
//Err:      "namespaces \"non-existance\" not found",
//},
//})
//})

//var _ = gauge.Step("Create pipeline file <path_to_pipeline_yaml> - with unsupported file format", func(path_to_pipeline_yaml string) {
//cmd.AssertOutput(&cmd.Cmd{
//Args: []string{store.Tkn().Path, "pipeline", "create", "--from", filepath.Join(helper.RootDir(), path_to_pipeline_yaml), "-n", store.Namespace()},
//Expected: icmd.Expected{
//ExitCode: 1,
//Err:      "inavlid file format for " + filepath.Join(helper.RootDir(), path_to_pipeline_yaml) + ": .yaml or .yml file extension and format required",
//},
//})
//})

//var _ = gauge.Step("Create pipeline from file <path_to_pipeline_yaml> - with mismatched Resource kind", func(path_to_pipeline_yaml string) {
//cmd.AssertOutput(&cmd.Cmd{
//Args: []string{store.Tkn().Path, "pipeline", "create", "--from", filepath.Join(helper.RootDir(), path_to_pipeline_yaml), "-n", store.Namespace()},
//Expected: icmd.Expected{
//ExitCode: 1,
//Err:      "provided kind PipelineRun instead of kind Pipeline",
//},
//})
//})

//var _ = gauge.Step("Existing pipeline", func() {
//cmd.AssertOutput(&cmd.Cmd{
//Args: []string{store.Tkn().Path, "pipeline", "create", "--from", filepath.Join(helper.RootDir(), "../testdata") + "/pipeline.yaml", "-n", store.Namespace()},
//Expected: icmd.Expected{
//ExitCode: 1,
//Err:      "failed to create pipeline \"test-pipeline\": pipelines.tekton.dev \"test-pipeline\" already exists\n",
//},
//})
//})
