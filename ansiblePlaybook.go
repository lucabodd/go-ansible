package ansibler

import (
	"errors"

	"os"

	common "github.com/apenella/go-common-utils/data"
)

const (
	// AnsiblePlaybookBin is the ansible-playbook binary file value
	AnsiblePlaybookBin = "ansible-playbook"

	// ConnectionFlag is the connection flag for ansible-playbook
	ConnectionFlag = "--connection"

	// ExtraVarsFlag is the extra variables flag for ansible-playbook
	ExtraVarsFlag = "--extra-vars"

	// FlushCacheFlag is the flush cache flag for ansible-playbook
	FlushCacheFlag = "--flush-cache"

	// InventoryFlag is the inventory flag for ansible-playbook
	InventoryFlag = "--inventory"

	// LimitFlag is the limit flag for ansible-playbook
	LimitFlag = "--limit"

	// ListHostsFlag is the list hosts flag for ansible-playbook
	ListHostsFlag = "--list-hosts"

	// ListTagsFlag is the list tags flag for ansible-playbook
	ListTagsFlag = "--list-tags"

	// ListTasksFlag is the list tasks flag for ansible-playbook
	ListTasksFlag = "--list-tasks"

	// TagsFlag is the tags flag for ansible-playbook
	TagsFlag = "--tags"

	// SyntaxCheckFlag is the syntax check flag for ansible-playbook
	SyntaxCheckFlag = "--syntax-check"

	// VaultPasswordFileFlag is the vault password file flag for ansible-playbook
	VaultPasswordFileFlag = "--vault-password-file"

	// AnsibleForceColorEnv is the environment variable which forces color mode
	AnsibleForceColorEnv = "ANSIBLE_FORCE_COLOR"
)

// PlaybookCmd object is the main object which defines the `ansible-playbook` command and how to execute it.
type PlaybookCmd struct {
	// Exec is the executor item
	Exec Executor
	// Playbook is the ansible's playbook name to be used
	Playbook string
	// Options are the ansible's playbook options
	Options *PlaybookOptions
	// ConnectionOptions are the ansible's playbook specific options for connection
	ConnectionOptions *PlaybookConnectionOptions
	// Writer manages the output
	Res PlaybookResults
}

// PlaybookOptions object has those parameters described on `Options` section within ansible-playbook's man page, and which defines which should be the ansible-playbook execution behavior.
type PlaybookOptions struct {
	// ExtraVars is a map of extra variables used on ansible-playbook execution
	ExtraVars map[string]interface{}
	// FlushCache clear the fact cache for every host in inventory
	FlushCache bool
	// Inventory specify inventory host path
	Inventory string
	// Limit is selected hosts additional pattern
	Limit string
	// ListHosts outputs a list of matching hosts
	ListHosts bool
	// ListTags list all available tags
	ListTags bool
	// ListTasks
	ListTasks bool
	// Tags list all tasks that would be executed
	Tags string
}

// PlaybookConnectionOptions object has those parameters described on `Connections Options` section within ansible-playbook's man page, and which defines how to connect to hosts.
type PlaybookConnectionOptions struct {
	// Connection is the type of connection used by ansible-playbook
	Connection string
}

type PlaybookResults struct {
	Stdout string
	TimeElapsed string
}

// AnsibleForceColor change to a forced color mode
func AnsibleForceColor() {
	os.Setenv(AnsibleForceColorEnv, "true")
}

// Run method runs the ansible-playbook
func (p *PlaybookCmd) Run() error {
	if p == nil {
		return errors.New("(ansible:Run) PlaybookCmd is nil")
	}

	// Generate the command to be run
	cmd, err := p.Command()
	if err != nil {
		return errors.New("(ansible:Run) -> " + err.Error())
	}

	err = p.Exec.Execute(cmd[0], cmd[1:])
	p.Res.Stdout = p.Exec.Stdout
	p.Res.TimeElapsed = p.Exec.TimeElapsed

	// Execute the command an return
	return err
}

// Command generate the ansible-playbook command which will be executed
func (p *PlaybookCmd) Command() ([]string, error) {
	cmd := []string{}
	// Set the ansible-playbook binary file
	cmd = append(cmd, AnsiblePlaybookBin)

	// Determine the options to be set
	if p.Options != nil {
		options, err := p.Options.GenerateCommandOptions()
		if err != nil {
			return nil, errors.New("(ansible::Command) -> " + err.Error())
		}
		for _, option := range options {
			cmd = append(cmd, option)
		}
	}

	// Determine the connection options to be set
	if p.ConnectionOptions != nil {
		options, err := p.ConnectionOptions.GenerateCommandConnectionOptions()
		if err != nil {
			return nil, errors.New("(ansible::Command) -> " + err.Error())
		}
		for _, option := range options {
			cmd = append(cmd, option)
		}
	}

	// Include the ansible playbook
	cmd = append(cmd, p.Playbook)

	return cmd, nil
}

// GenerateCommandOptions return a list of options flags to be used on ansible-playbook execution
func (o *PlaybookOptions) GenerateCommandOptions() ([]string, error) {
	cmd := []string{}

	if o == nil {
		return nil, errors.New("(ansible::GenerateCommandOptions) PlaybookOptions is nil")
	}

	if o.FlushCache {
		cmd = append(cmd, FlushCacheFlag)
	}

	if o.Inventory != "" {
		cmd = append(cmd, InventoryFlag)
		cmd = append(cmd, o.Inventory)
	}

	if o.Limit != "" {
		cmd = append(cmd, LimitFlag)
		cmd = append(cmd, o.Limit)
	}

	if o.ListHosts {
		cmd = append(cmd, ListHostsFlag)
	}

	if o.ListTags {
		cmd = append(cmd, ListTagsFlag)
	}

	if o.ListTasks {
		cmd = append(cmd, ListTasksFlag)
	}

	if o.Tags != "" {
		cmd = append(cmd, TagsFlag)
		cmd = append(cmd, o.Tags)
	}

	if len(o.ExtraVars) > 0 {
		cmd = append(cmd, ExtraVarsFlag)
		extraVars, err := o.generateExtraVarsCommand()
		if err != nil {
			return nil, errors.New("(ansible::GenerateCommandOptions) -> " + err.Error())
		}
		cmd = append(cmd, extraVars)
	}

	return cmd, nil
}

// generateExtraVarsCommand return an string which is a json structure having all the extra variable
func (o *PlaybookOptions) generateExtraVarsCommand() (string, error) {

	extraVars, err := common.ObjectToJSONString(o.ExtraVars)
	if err != nil {
		return "", errors.New("(ansible::generateExtraVarsCommand) -> " + err.Error())
	}
	return extraVars, nil
}

// AddExtraVar registers a new extra variable on ansible-playbook options item
func (o *PlaybookOptions) AddExtraVar(name string, value interface{}) error {

	if o.ExtraVars == nil {
		o.ExtraVars = map[string]interface{}{}
	}
	_, exists := o.ExtraVars[name]
	if exists {
		return errors.New("(ansible::AddExtraVar) ExtraVar '" + name + "' already exist")
	}

	o.ExtraVars[name] = value

	return nil
}

// GenerateCommandConnectionOptions return a list of connection options flags to be used on ansible-playbook execution
func (o *PlaybookConnectionOptions) GenerateCommandConnectionOptions() ([]string, error) {
	cmd := []string{}

	if o.Connection != "" {
		cmd = append(cmd, ConnectionFlag)
		cmd = append(cmd, o.Connection)
	}

	return cmd, nil
}
