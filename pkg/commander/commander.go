package commander

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/BlueSageSolutions/courier/pkg/transform"
	"github.com/BlueSageSolutions/courier/pkg/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

const (
	OP_CREATE                       string = "create"
	OP_DESTROY                      string = "destroy"
	OP_BUILD                        string = "build"
	OP_START                        string = "start"
	OP_INSTALL                      string = "install"
	OP_CONFIGURE                    string = "configure"
	CLIENT_KEY                      string = "__CLIENT__"
	ENVIRONMENT_KEY                 string = "__ENVIRONMENT__"
	ARG_TYPE_JQ                     string = "jq"
	ARG_TYPE_LITERAL                string = "literal"
	EXECUTABLE_AWS                  string = "aws"
	EXECUTABLE_ECHO                 string = "echo"
	EXECUTABLE_PERSIST_ARGUMENTS    string = "persist-arguments"
	EXECUTABLE_INTERNAL_CURL        string = "internal-curl"
	EXECUTABLE_CAT                  string = "cat"
	EXECUTABLE_AZURE                string = "az"
	COMMAND_ARG_SOURCE_TYPE_HTTPS   string = "https"
	COMMAND_ARG_SOURCE_TYPE_HTTP    string = "http"
	COMMAND_ARG_SOURCE_TYPE_FILE    string = "file"
	COMMAND_ARG_SOURCE_TYPE_JSON    string = "json"
	COMMAND_ARG_SOURCE_TYPE_TEXT    string = "text"
	COMMAND_ARG_STYLE_PLAIN         string = "plain"
	COMMAND_ARG_STYLE_LONG          string = "long"
	COMMAND_ARG_STYLE_SHORT         string = "short"
	COMMAND_QUOTE_TYPE_NONE         string = ""
	COMMAND_QUOTE_TYPE_SINGLE_QUOTE string = "single"
	COMMAND_QUOTE_TYPE_DOUBLE_QUOTE string = "double"
	SLEEP_BEFORE                    bool   = true
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type ClientProfile struct {
	Name               string   `yaml:"name"`
	InfrastructureType string   `yaml:"infrastructure-type"`
	Contacts           []string `yaml:"contacts"`
	EmailDomain        string   `yaml:"email-domain"`
}

type EnvironmentDesc struct {
	Name       string   `yaml:"name"`
	StartDate  string   `yaml:"start-date"`
	Operations []string `yaml:"operations"`
}

type DatabaseDesc struct {
	Name       string   `yaml:"name"`
	Operations []string `yaml:"operations"`
}

type ServiceDesc struct {
	Name       string   `yaml:"name"`
	URL        string   `yaml:"url"`
	Operations []string `yaml:"operations"`
}

type ManifestResults struct {
}

type Manifest struct {
	ClientProfile ClientProfile `yaml:"client-profile"`
	Environments  []EnvironmentDesc
	Services      []ServiceDesc
	Databases     []DatabaseDesc
}

type OutputChannel struct {
	Output []byte
	Error  error
}

type ExecutionContext struct {
	Client      string `yaml:"client"`
	Environment string `yaml:"environment"`
}

type ResourceList []Resource

type ScriptError struct {
	ErrorMessage string `yaml:"error-message"`
}

type Script []Command
type DeploymentScriptResult []Result
type DeploymentScriptResults map[string]DeploymentScriptResult
type DeploymentScriptSuiteResults struct {
	Results   map[string]DeploymentScriptResults `yaml:"results"`
	Directory string                             `yaml:"directory"`
}
type DeploymentScripts map[string]*DeploymentScriptList

type DeploymentScriptList struct {
	DeploymentScripts []DeploymentScript `yaml:"scripts"`
}

type Sleep struct {
	Timeout       int64  `yaml:"timeout"`
	Before        int64  `yaml:"before"`
	After         int64  `yaml:"after"`
	AfterMessage  string `yaml:"after-message"`
	BeforeMessage string `yaml:"before-message"`
}

type Command struct {
	ScriptReference string        `yaml:"script-reference"`
	Executable      string        `yaml:"executable"`
	Name            string        `yaml:"command"`
	Description     string        `yaml:"description"`
	Sensitive       bool          `yaml:"sensitive"`
	Source          string        `yaml:"source"`
	Replacements    []Replacement `yaml:"replacements"`
	Environment     []string      `yaml:"environment"`
	Directory       string        `yaml:"directory"`
	SubCommand      string        `yaml:"sub-command"`
	Arguments       []Argument    `yaml:"arguments"`
	Sleep           Sleep         `yaml:"sleep"`
}

type Argument struct {
	Name          string                    `yaml:"name"`
	Description   string                    `yaml:"description"`
	Randomize     int                       `yaml:"randomize"`
	Value         string                    `yaml:"value"`
	Style         string                    `yaml:"style"`
	QuoteType     string                    `yaml:"quote-type"`
	SourceType    string                    `yaml:"source-type"`
	Source        string                    `yaml:"source"`
	Interpolation *transform.Transformation `yaml:"interpolation"`
}

type Replacement struct {
	Match             string `yaml:"match"`
	ReplaceWith       string `yaml:"replace-with"`
	ReplaceWithRandom int    `yaml:"replace-with-random"`
}

type Source struct {
	Transformations []transform.Transformation `yaml:"transformations"`
	Data            string                     `yaml:"data"`
}

type DeploymentScript struct {
	Name             string            `yaml:"script"`
	Description      string            `yaml:"description"`
	ExecutionContext *ExecutionContext `yaml:"execution-context"`
	Path             string            `yaml:"path"`
	Sources          map[string]Source `yaml:"sources"`
	SetupScript      Script            `yaml:"setup"`
	DryRun           bool              `yaml:"run-main"`
	MainScript       Script            `yaml:"main"`
	Destroy          bool              `yaml:"run-cleanup"`
	CleanupScript    Script            `yaml:"cleanup"`
}

type Result struct {
	Script      string          `yaml:"script"`
	Sensitive   bool            `yaml:"sensitive"`
	Command     string          `yaml:"command"`
	Description string          `yaml:"description"`
	Output      json.RawMessage `yaml:"output"`
	Error       string          `yaml:"error"`
}

type Resource struct {
	Resource string   `yaml:"resource"`
	Package  string   `yaml:"package"`
	Actions  []string `yaml:"actions"`
}

func Sluggify(plissken string) string {
	snake := matchFirstCap.ReplaceAllString(plissken, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func removeWhiteSpace(jsonStr string) string {
	var sb strings.Builder
	inQuotes := false

	for _, r := range jsonStr {
		if r == '"' {
			inQuotes = !inQuotes
		}

		if unicode.IsSpace(r) && !inQuotes {
			continue
		}
		sb.WriteRune(r)
	}

	return sb.String()
}

func (result Result) WriteOutput() (string, error) {
	filename := fmt.Sprintf("/tmp/%s.json", strings.ReplaceAll(result.Script, ":", "_"))
	err := os.WriteFile(filename, result.Output, 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func WriteCommandOutput(filename string, output []byte) error {
	err := os.WriteFile(filename, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

func EmptyResult() json.RawMessage {
	return json.RawMessage("{}")
}

func LoadResources(file string) (ResourceList, error) {
	var list ResourceList
	yfile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &list)

	if err != nil {
		return nil, err
	}
	return list, nil
}

func IsYamlFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fileInfo.IsDir() {
		return false
	}
	return filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml"
}

func LoadResourceLists(path string) (ResourceList, error) {
	resourceList := make(ResourceList, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsYamlFile(path) {
			list, err := LoadResources(path)
			if err != nil {
				util.GetLogger().Error(fmt.Sprintf("initialization: resources file in error (%s)", path), zap.Error(err))
			}
			resourceList = append(resourceList, list...)
		}
		return nil
	})

	return resourceList, err
}

func LoadEnvironmentVariables(path, filename string) ([]string, error) {
	fqPath := filename
	if len(path) > 0 {
		fqPath = path + "/" + filename
	}
	file, err := os.Open(fqPath)
	if err != nil {
		util.GetLogger().Error("failed to open file", zap.Error(err))
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip lines that start with '#'
		if strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		util.GetLogger().Error("error reading lines", zap.Error(err))
		return nil, err
	}
	return lines, nil
}

func LoadDeploymentScript(file string) (*DeploymentScriptList, error) {
	var list []DeploymentScript
	yfile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &list)

	if err != nil {
		return nil, err
	}
	for index, deploymentScript := range list {
		for key, value := range deploymentScript.Sources {
			deploymentScript.Sources[key] = value
		}
		deploymentScript.Path = file
		list[index] = deploymentScript
	}
	deploymentScriptList := &DeploymentScriptList{DeploymentScripts: list}
	return deploymentScriptList, nil
}

func LoadDeploymentScripts(path string) (*DeploymentScripts, error) {
	deploymentScripts := make(DeploymentScripts)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsYamlFile(path) {
			deploymentScript, err := LoadDeploymentScript(path)
			if err != nil {
				util.GetLogger().Error(fmt.Sprintf("initialization: deployment scripts file in error (%s)", path), zap.Error(err))
			}
			deploymentScripts[path] = deploymentScript
		}
		return nil
	})

	return &deploymentScripts, err
}

// TODO fix names here
func (source Source) Resolve(sourceName string) string {
	resolved := string(EmptyResult())
	if len(source.Data) > 0 {
		if source.Transformations != nil {
			var err error
			eventData := []byte(source.Data)
			dataPayload := &transform.DataPayload{Input: eventData, Outputs: make(transform.Transformed)}
			for _, transformation := range source.Transformations {
				transformations, err := transformation.MungeAndTransform(dataPayload)
				if err != nil {
					return resolved
				}
				for key, value := range transformations {
					if len(value) > 0 {
						dataPayload.Outputs[key] = value
					}
				}
			}
			rawData, err := json.Marshal(dataPayload)
			if err != nil {
				return resolved
			}
			resolved = string(rawData)
		} else {
			resolved = source.Data
		}
	}
	return resolved
}

func (command Command) WriteInterpolatedSource(source, interpolated string) (string, error) {
	filename := fmt.Sprintf("/tmp/_%s.json", strings.ReplaceAll(source, ":", "_"))
	err := os.WriteFile(filename, []byte(interpolated), 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (command Command) InterpolateSource(deploymentScript DeploymentScript, source string, outputs DeploymentScriptResult) string {
	raw := ResolveSourceFromFile(deploymentScript, source, outputs)
	var replaceWith string

	for _, replacement := range command.Replacements {
		if replacement.ReplaceWithRandom > 0 {
			var myErr error
			replaceWith, myErr = transform.GenerateRandomString(replacement.ReplaceWithRandom)
			if myErr != nil {
				replaceWith = ""
			}
		} else {
			switch replacement.ReplaceWith {
			case CLIENT_KEY:
				replaceWith = deploymentScript.ExecutionContext.Client
			case ENVIRONMENT_KEY:
				replaceWith = deploymentScript.ExecutionContext.Environment
			default:
				replaceWith = ResolveSourceFromFile(deploymentScript, replacement.ReplaceWith, outputs)
			}
			replaceWith = strings.TrimSpace(replaceWith)
		}
		raw = strings.ReplaceAll(raw, replacement.Match, replaceWith)
	}
	return raw
}

func ResolveSourceFromFile(deploymentScript DeploymentScript, source string, outputs DeploymentScriptResult) string {
	filename := deploymentScript.SourceAsFileName(source, outputs)
	yfile, err := os.ReadFile(filename)
	if err != nil {
		return string(EmptyResult())
	}
	return string(yfile)

}

func (argument Argument) ResolveSource(deploymentScript DeploymentScript, outputs DeploymentScriptResult) string {
	_, ok := deploymentScript.Sources[argument.Source]
	if !ok {
		for _, result := range outputs {
			if result.Script == argument.Source {
				return string(result.Output)
			}
		}
		return string(EmptyResult())
	}

	return ResolveSourceFromFile(deploymentScript, argument.Source, nil)
}

func (argument Argument) ResolveFileName(executable string, deploymentScript DeploymentScript, outputs DeploymentScriptResult) string {
	if len(argument.Source) == 0 {
		return argument.Source
	}
	filename := deploymentScript.SourceAsFileName(argument.Source, outputs)
	if argument.SourceType == COMMAND_ARG_SOURCE_TYPE_FILE {
		filename = fmt.Sprintf("file://%s", filename)
	}
	return filename
}

func (argument Argument) Enquote(value string) string {
	argumentValue := value
	switch argument.QuoteType {
	case COMMAND_QUOTE_TYPE_NONE:
	case COMMAND_QUOTE_TYPE_SINGLE_QUOTE:
		argumentValue = fmt.Sprintf("'%s'", argumentValue)
	case COMMAND_QUOTE_TYPE_DOUBLE_QUOTE:
		argumentValue = fmt.Sprintf("\"%s\"", argumentValue)
	}
	return argumentValue
}

func (argument Argument) Resolve(executable string, deploymentScript DeploymentScript, outputs DeploymentScriptResult) string {
	if len(argument.Value) > 0 {
		return argument.Enquote(argument.Value)
	} else if argument.Randomize > 0 {
		random, err := transform.GenerateRandomString(argument.Randomize)
		if err != nil {
			return ""
		}
		return argument.Enquote(random)
	}
	argumentValue := ""
	source := argument.ResolveSource(deploymentScript, outputs)

	msg := &transform.DataPayload{Input: json.RawMessage(source)}

	if argument.Interpolation != nil {
		argument.Interpolation.Name = argument.Name
		transformations, err := argument.Interpolation.MungeAndTransform(msg)
		if err != nil {
			util.GetLogger().Info(fmt.Sprintf("failed to resolve argument %s", err), zap.Any("argument", argument))
			return argumentValue
		}
		argumentValue = string(transformations[argument.Name])
	} else {
		switch argument.SourceType {
		case COMMAND_ARG_SOURCE_TYPE_JSON:
			argumentValue = removeWhiteSpace(argument.ResolveSource(deploymentScript, outputs))
		case COMMAND_ARG_SOURCE_TYPE_TEXT:
			argumentValue = argument.ResolveSource(deploymentScript, outputs)
		default:
			argumentValue = argument.ResolveFileName(executable, deploymentScript, outputs)
		}
	}
	return argument.Enquote(argumentValue)
}

func (deploymentScriptList DeploymentScriptList) Execute(exe *ExecutionContext, file string) (DeploymentScriptResults, []error) {
	deploymentScriptResults := make(DeploymentScriptResults)
	var errors []error
	for _, deploymentScript := range deploymentScriptList.DeploymentScripts {
		fmt.Printf("script path = %s\n", deploymentScript.Path)
		deploymentScript.ExecutionContext = exe
		outputs, err := deploymentScript.Execute()
		if err != nil {
			errors = append(errors, err)
			util.GetLogger().Error("deployment scripts failed", zap.Error(err))
		}
		deploymentScriptResults[deploymentScript.Path] = outputs
	}
	return deploymentScriptResults, errors
}

func (deploymentScript DeploymentScript) SourceAsFileName(source string, outputs DeploymentScriptResult) string {
	var err error
	_, ok := deploymentScript.Sources[source]
	filename := fmt.Sprintf("/tmp/%s.%s.json", deploymentScript.Name, source)
	if !ok {
		// need to render an output as a file
		for _, result := range outputs {
			if result.Script == source {
				filename, err = result.WriteOutput()
				if err != nil {
					util.GetLogger().Error(fmt.Sprintf("failed to write %s", filename), zap.Error(err))
				}
			}
		}
	}

	return filename
}

func (deploymentScript DeploymentScript) InitializeSources() error {
	for key, source := range deploymentScript.Sources {
		resolved := source.Resolve(key)
		err := os.WriteFile(deploymentScript.SourceAsFileName(key, nil), []byte(resolved), 0644)
		if err != nil {
			util.GetLogger().Error("failed to write source into a file", zap.Error(err))
			return err
		}
	}
	return nil
}

func (deploymentScript DeploymentScript) Execute() (DeploymentScriptResult, error) {
	outputs := make(DeploymentScriptResult, 0)
	err := deploymentScript.InitializeSources()
	if err != nil {
		return outputs, err
	}
	outputs, err = deploymentScript.SetupScript.Execute("setup", deploymentScript, outputs)
	if err != nil {
		return outputs, err
	}

	// We don't exit here - we need to try to cleanup first
	if !deploymentScript.DryRun {
		outputs, err = deploymentScript.MainScript.Execute("main", deploymentScript, outputs)
		if err != nil {
			util.GetLogger().Error(deploymentScript.Name, zap.Error(err))
		}
	}

	if deploymentScript.Destroy {
		outputs, err = deploymentScript.CleanupScript.Execute("cleanup", deploymentScript, outputs)
		if err != nil {
			return outputs, err
		}
	}
	return outputs, err
}

type Safe struct {
	Message string
}

func safeJson(unsafe []byte) []byte {
	return []byte(strings.ReplaceAll(string(unsafe), "\n", ""))
}

func (command Command) ArgsAsBytes(cmd *exec.Cmd) []byte {
	return []byte(strings.Join(cmd.Args, ""))
}

func (command Command) ScriptReferenceAsFilename() string {
	filename := fmt.Sprintf("/tmp/%s.json", strings.ReplaceAll(command.ScriptReference, ":", "_"))
	return filename
}

func (command Command) ExecuteWithTimeout(cmd *exec.Cmd) (json.RawMessage, error) {
	var output []byte
	var err error
	if command.Sleep.Timeout == 0 {
		switch command.Executable {
		case EXECUTABLE_PERSIST_ARGUMENTS:
			if err == nil {
				err = WriteCommandOutput(command.ScriptReferenceAsFilename(), command.ArgsAsBytes(cmd))
			}
		case EXECUTABLE_AWS:
			output, err = cmd.CombinedOutput()
		case EXECUTABLE_AZURE:
			// azure's cli is whack. this may not work at times
			output, err = cmd.CombinedOutput()
		case EXECUTABLE_ECHO:
			output, err = cmd.Output()
			if err == nil {
				err = WriteCommandOutput(command.Name, output)
			}
		default:
			output, err = cmd.CombinedOutput()
		}
	} else {
		ch := make(chan OutputChannel)

		go func() {
			var channelOutput []byte
			var channelError error
			switch command.Executable {
			case EXECUTABLE_AWS:
				channelOutput, channelError = cmd.CombinedOutput()
			case EXECUTABLE_AZURE:
				channelOutput, channelError = cmd.Output()
			case EXECUTABLE_ECHO:
				channelOutput, channelError = cmd.Output()
				if channelError == nil {
					channelError = WriteCommandOutput(command.Name, channelOutput)
				}

			default:
				channelOutput, channelError = cmd.CombinedOutput()
			}
			ch <- OutputChannel{Output: channelOutput, Error: channelError}
		}()

		select {
		case <-time.After(time.Duration(command.Sleep.Timeout) * time.Second):
			err = fmt.Errorf("'%s' timed out after %d seconds", cmd.String(), command.Sleep.Timeout)
		case x := <-ch:
			output = x.Output
			err = x.Error
		}
	}
	return output, err
}

func CatCommand(filename string) (*exec.Cmd, error) {
	catPath, err := exec.LookPath(EXECUTABLE_CAT)
	if err != nil {
		return nil, err
	}
	catCmd := exec.Command(catPath, filename)
	return catCmd, nil
}

func (command Command) ExecuteCatPipe(deploymentScript DeploymentScript, outputs DeploymentScriptResult) (string, json.RawMessage, error) {
	filename := deploymentScript.SourceAsFileName(command.Source, outputs)
	if command.Replacements != nil {
		interpolated := command.InterpolateSource(deploymentScript, command.Source, outputs)
		filename, _ = command.WriteInterpolatedSource(command.Source, interpolated)
	}

	catCmd, err := CatCommand(filename)
	if err != nil {
		return "", EmptyResult(), err
	}
	theCmd, err := command.BuildCmd(deploymentScript, outputs)
	if err != nil {
		return "", EmptyResult(), err
	}

	pipe, err := catCmd.StdoutPipe()
	if err != nil {
		return "", EmptyResult(), err
	}
	theCmd.Stdin = pipe

	err = catCmd.Start()
	if err != nil {
		return "", EmptyResult(), err
	}

	output, err := theCmd.Output()
	if err != nil {
		return fmt.Sprintf("%s | %s", catCmd.String(), theCmd.String()), EmptyResult(), err
	}

	return fmt.Sprintf("%s | %s", catCmd.String(), theCmd.String()), output, err
}

func (command Command) BuildCmd(deploymentScript DeploymentScript, outputs DeploymentScriptResult) (*exec.Cmd, error) {
	args := make([]string, 1)
	var executablePath string
	var err error

	if len(command.Executable) == 0 {
		command.Executable = EXECUTABLE_AWS
	}
	if command.Executable != EXECUTABLE_PERSIST_ARGUMENTS {
		executablePath, err = exec.LookPath(command.Executable)
		if err != nil {
			return nil, err
		}
		if command.Executable != EXECUTABLE_ECHO {
			if len(command.Name) > 0 {
				args = append(args, command.Name)
			}
			if len(command.SubCommand) > 0 {
				args = append(args, command.SubCommand)
			}
		}
	} else {
		executablePath = "no-command"
	}

	for _, argument := range command.Arguments {
		switch argument.Style {
		case COMMAND_ARG_STYLE_PLAIN:
			// no name
		case COMMAND_ARG_STYLE_LONG:
			args = append(args, fmt.Sprintf("--%s", argument.Name))
		case COMMAND_ARG_STYLE_SHORT:
			args = append(args, fmt.Sprintf("-%s", argument.Name))
		default:
			args = append(args, fmt.Sprintf("--%s", argument.Name))
		}
		args = append(args, argument.Resolve(command.Executable, deploymentScript, outputs))
	}

	var environment []string
	for _, shellScript := range command.Environment {
		vars, err := LoadEnvironmentVariables(command.Directory, shellScript)
		if err != nil {
			return nil, err
		}
		environment = append(environment, vars...)
	}
	cmd := &exec.Cmd{
		Dir:  command.Directory,
		Path: executablePath,
		Env:  environment,
		Args: args,
	}
	return cmd, nil
}

func (command Command) ExecuteNoPipe(deploymentScript DeploymentScript, outputs DeploymentScriptResult) (string, json.RawMessage, error) {
	cmd, err := command.BuildCmd(deploymentScript, outputs)
	if err != nil {
		return "", EmptyResult(), err
	}

	output, err := command.ExecuteWithTimeout(cmd)

	if err != nil {
		return cmd.String(), safeJson(output), err
	}

	return cmd.String(), output, nil
}

func (command Command) SleepBefore(cmdString, label string) {
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("[%s] prior to executing a command in the script: %s\n", Timestamp(), label)
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("\tdelay-before: %d\n", command.Sleep.Before)
	if command.Sensitive {
		fmt.Print("\tcommand contains sensitive data: REDACTED\n")
	} else {
		fmt.Printf("\tcommand to execute: %s\n", cmdString)
	}

	if len(command.Sleep.BeforeMessage) > 0 {
		fmt.Printf("\tmessage: %s\n", command.Sleep.BeforeMessage)
	}

	time.Sleep(time.Duration(command.Sleep.Before) * time.Second)
}

func (command Command) SleepAfter(cmdString, label string, message json.RawMessage, hasError bool) {
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("[%s] after executing a command in the script: %s\n", Timestamp(), label)
	fmt.Println("----------------------------------------------------------------------")
	if hasError {
		fmt.Printf("\tthere was an error at: %s\n", command.ScriptReference)
		fmt.Println("----------------------------------------------------------------------")
		return
	}
	if command.Sensitive {
		fmt.Print("\tcommand was executed with sensitive data: REDACTED\n")
	} else {
		fmt.Printf("\tcommand executed: %s\n", cmdString)
	}
	fmt.Printf("\tdelay-after: %d\n", command.Sleep.After)
	if command.Sensitive {
		fmt.Print("\tresults:\nREDACTED\n")
	} else {
		fmt.Printf("\tresults:\n%s\n", string(message))
	}

	if len(command.Sleep.AfterMessage) > 0 {
		fmt.Printf("\tmessage: %s\n", command.Sleep.AfterMessage)
	}
	time.Sleep(time.Duration(command.Sleep.After) * time.Second)
}

func (command Command) Execute(deploymentScript DeploymentScript, outputs DeploymentScriptResult) (string, json.RawMessage, error) {
	var commandLine string
	var output json.RawMessage
	var err error
	theCmd, err := command.BuildCmd(deploymentScript, outputs)
	if err != nil {
		return "", EmptyResult(), err
	}

	command.SleepBefore(theCmd.String(), deploymentScript.Path)

	if len(command.Source) > 0 && command.Executable != EXECUTABLE_PERSIST_ARGUMENTS {
		commandLine, output, err = command.ExecuteCatPipe(deploymentScript, outputs)
		command.SleepAfter(commandLine, deploymentScript.Path, output, err != nil)
		return commandLine, output, err
	}
	commandLine, output, err = command.ExecuteNoPipe(deploymentScript, outputs)
	command.SleepAfter(commandLine, deploymentScript.Path, output, err != nil)
	return commandLine, output, err
}

func (script Script) Execute(phase string, deploymentScript DeploymentScript, outputs DeploymentScriptResult) (DeploymentScriptResult, error) {
	for index := range script {
		command := script[index]
		command.ScriptReference = fmt.Sprintf("%s:%s:step-%d", deploymentScript.Name, phase, index)
		fmt.Printf("executing: %s\n%s\n\n", command.ScriptReference, command.Description)
		cmd, jsonOutput, err := command.Execute(deploymentScript, outputs)
		if err != nil {
			scriptError := &ScriptError{ErrorMessage: string(jsonOutput)}
			betterMessage, marshalErr := json.Marshal(scriptError)
			if marshalErr != nil {
				util.GetLogger().Error("marshalling aws error failed", zap.Error(marshalErr))
				outputs = append(outputs,
					Result{Sensitive: command.Sensitive,
						Script:      command.ScriptReference,
						Output:      jsonOutput,
						Command:     cmd,
						Description: command.Description,
						Error:       fmt.Sprintf("%s", err)})
			} else {
				outputs = append(outputs,
					Result{Sensitive: command.Sensitive,
						Script:      command.ScriptReference,
						Output:      betterMessage,
						Command:     cmd,
						Description: command.Description,
						Error:       fmt.Sprintf("%s", err)})
			}
			return outputs, err
		}
		if len(jsonOutput) == 0 {
			jsonOutput = EmptyResult()
		}
		outputs = append(outputs,
			Result{Sensitive: command.Sensitive,
				Script:      command.ScriptReference,
				Description: command.Description,
				Output:      jsonOutput,
				Command:     cmd})
	}
	return outputs, nil
}

func (deploymentScripts DeploymentScripts) Execute(exe *ExecutionContext) *DeploymentScriptSuiteResults {
	g, _ := errgroup.WithContext(context.Background())

	// Mutex to protect concurrent access to the results map
	var mutex sync.Mutex
	var results DeploymentScriptSuiteResults
	results.Results = make(map[string]DeploymentScriptResults)
	for key := range deploymentScripts {
		// It's important to copy the loop variable when using it in a goroutine
		deploymentScript := deploymentScripts[key]
		name := key
		g.Go(func() error {
			deployed, errors := deploymentScript.Execute(exe, name)
			if errors != nil {
				util.GetLogger().Info(fmt.Sprintf("errors: %+v", zap.Any("errors", errors)))
			}

			// Use mutex to protect writing to the results map
			mutex.Lock()
			results.Results[name] = deployed
			mutex.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil
	}
	return &results
}

func writeCode(file *os.File, label, value string) error {
	_, err := file.WriteString(fmt.Sprintf("**%s**: `%s`\n\n", label, value))
	if err != nil {
		return err
	}
	return nil
}

func writeCodeBlock(file *os.File, encoding, label, value string) error {
	_, err := file.WriteString(fmt.Sprintf("**%s**:\n\n", label))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", encoding, value))
	if err != nil {
		return err
	}
	return nil
}
func (result Result) AsMarkdown(file *os.File) error {
	var err error
	if result.Sensitive {
		err = writeCode(file, "Command", "REDACTED: Command may contain sensitive data")
		if err != nil {
			return err
		}
	} else {
		err = writeCode(file, "Command", result.Command)
		if err != nil {
			return err
		}
	}
	err = writeCode(file, "Description", result.Description)
	if err != nil {
		return err
	}
	err = writeCode(file, "Script Reference", result.Script)
	if err != nil {
		return err
	}
	if len(result.Error) > 0 {
		err := writeCode(file, "Error", result.Error)
		if err != nil {
			return err
		}
	}
	if len(result.Output) > 0 {
		var err error

		if result.Sensitive {
			err = writeCodeBlock(file, "json", "Output", "REDACTED")
		} else {
			err = writeCodeBlock(file, "json", "Output", string(result.Output))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (deploymentScript DeploymentScript) AsMarkdown(file *os.File) error {
	markdown, err := yaml.Marshal(deploymentScript)
	if err != nil {
		return err
	}
	manifest := &Manifest{}
	manifest.ClientProfile.Contacts = append(manifest.ClientProfile.Contacts, "fred@foo.com")
	environment := &EnvironmentDesc{}
	environment.Operations = append(environment.Operations, "create")
	service := &ServiceDesc{}
	service.Operations = append(environment.Operations, "start")
	manifest.Environments = append(manifest.Environments, *environment)
	manifest.Services = append(manifest.Services, *service)
	manifestBlock, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	writeCodeBlock(file, "yaml", "MANIFEST", string(manifestBlock))
	return writeCodeBlock(file, "yaml", "Executed", string(markdown))
}

func (resultResults DeploymentScriptResults) AsTocEntry(toc *os.File, label, absolutePath, relativePath string) error {
	pieces := strings.Split(absolutePath, relativePath)
	if pieces == nil {
		return nil
	}
	pieces = strings.Split(pieces[1], string(filepath.Separator))
	if pieces == nil {
		return nil
	}
	nameEntry := pieces[len(pieces)-1]
	if len(label) == 0 {
		label = nameEntry
	}
	_, err := toc.WriteString(fmt.Sprintf("[%s](./%s)\n\n", label, nameEntry))
	return err
}

func (resultResults DeploymentScriptResults) AsMarkdown(relativePath, absolutePath string, toc *os.File) error {
	for key, deploymentScriptResult := range resultResults {
		outputFilename := strings.ReplaceAll(key, ".yaml", "")
		outputFilename = GetReportLocation(absolutePath, outputFilename, "results.md")
		output, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer output.Close()
		err = writeCode(output, "Executing", key)
		if err != nil {
			return err
		}
		err = deploymentScriptResult.AsMarkdown(output)
		if err != nil {
			return err
		}
		output.Close()
		resultResults.AsTocEntry(toc, "", outputFilename, relativePath)
	}
	return nil
}

func (deploymentScriptResult DeploymentScriptResult) AsMarkdown(file *os.File) error {
	for _, result := range deploymentScriptResult {
		err := result.AsMarkdown(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetReportLocation(path, originalName, newExtension string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return filepath.Join(path, originalName)
	}
	name := strings.Replace(originalName, pwd, "", 1)
	name = strings.ReplaceAll(name, string(filepath.Separator), "_")
	return fmt.Sprintf("%s.%s", filepath.Join(path, name), newExtension)
}

func (deploymentScriptSuiteResults DeploymentScriptSuiteResults) AsMarkdown(deploymentScripts DeploymentScripts, relativePath, absolutePath string, toc *os.File) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	for key, deploymentScriptResults := range deploymentScriptSuiteResults.Results {
		deploymentScriptList, ok := deploymentScripts[key]
		if ok {
			for _, deploymentScript := range deploymentScriptList.DeploymentScripts {
				deploymentScriptResults.AsTocEntry(toc,
					strings.ReplaceAll(deploymentScript.Path, pwd, ""),
					GetReportLocation(absolutePath, deploymentScript.Path, "dump.md"),
					relativePath)

				output, err := os.OpenFile(GetReportLocation(absolutePath,
					deploymentScript.Path, "dump.md"),
					os.O_APPEND|os.O_CREATE|os.O_WRONLY,
					0644)
				if err != nil {
					return err
				}

				defer output.Close()

				err = deploymentScript.AsMarkdown(output)
				if err != nil {
					return err
				}
				output.Close()
			}
			err := deploymentScriptResults.AsMarkdown(relativePath, absolutePath, toc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Timestamp() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d.%02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func (results DeploymentScriptSuiteResults) Publish(reportsPath string, deploymentScripts DeploymentScripts) (string, error) {
	publishedFile := fmt.Sprintf("%s/README.md", results.Directory)
	tableOfContents, err := os.OpenFile(publishedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return publishedFile, err
	}

	defer tableOfContents.Close()
	_, err = tableOfContents.WriteString("# Executed\n\n")
	if err != nil {
		return publishedFile, err
	}
	err = results.AsMarkdown(deploymentScripts, reportsPath, results.Directory, tableOfContents)
	if err != nil {
		return publishedFile, err
	}

	return publishedFile, nil
}

func (resource Resource) DefaultSetupScript() Script {
	script := make([]Command, 0)
	script = append(script, Command{
		Executable: "aws",
		SubCommand: "get-caller-identity",
		Name:       "sts",
	})
	return script
}

func (resource Resource) GenerateScript() Script {
	script := make([]Command, 0)
	for _, action := range resource.Actions {
		script = append(script, Command{Executable: "aws", SubCommand: Sluggify(action), Name: resource.Package})
	}
	return script
}

func (deploymentScriptList DeploymentScriptList) Generate(path string) error {
	bytes, err := yaml.Marshal(deploymentScriptList.DeploymentScripts)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (resources ResourceList) Generate(filepath string) error {
	deploymentScriptList := &DeploymentScriptList{DeploymentScripts: make([]DeploymentScript, 0)}
	for _, resource := range resources {
		deploymentScript := &DeploymentScript{
			Name:        resource.Resource,
			SetupScript: resource.DefaultSetupScript(),
			MainScript:  resource.GenerateScript(),
		}
		deploymentScriptList.DeploymentScripts = append(deploymentScriptList.DeploymentScripts, *deploymentScript)
	}
	return deploymentScriptList.Generate(filepath)
}

func (manifest *Manifest) Validate() error {
	return nil
}

func LoadManifest(location string) (*Manifest, error) {
	var manifest Manifest
	yfile, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &manifest)

	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (manifest *Manifest) Process() (*ManifestResults, error) {
	os.Setenv("COURIER_CLIENT", manifest.ClientProfile.Name)
	for _, environment := range manifest.Environments {
		os.Setenv("COURIER_ENVIRONMENT", environment.Name)
		for _, database := range manifest.Databases {
			for _, operation := range database.Operations {
				database.Perform(operation, &manifest.ClientProfile, &environment)
			}
		}
	}
	return nil, nil
}

func (database *DatabaseDesc) Perform(operation string, client *ClientProfile, environment *EnvironmentDesc) {

}

func (results *ManifestResults) Publish(manifest *Manifest, path string) error {
	return nil
}
