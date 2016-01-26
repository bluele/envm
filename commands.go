package main

import (
	"fmt"
	"os"
)

type baseCommand struct{}

func (cm *baseCommand) lock() error {
	return mustGetLock().TryLock()
}

func (cm *baseCommand) unLock() error {
	return mustGetLock().Unlock()
}

type InitCommand struct {
	baseCommand
}

func (cm *InitCommand) Help() string {
	return "init command"
}

func (cm *InitCommand) Run(args []string) int {
	script := `envm () {
        local command
        command="$1"
        if [ "$#" -gt 0 ]
        then
                shift
        fi
        case "$command" in
           (use) eval "%vcommand envm "$command" "$@"%v" ;;
           (*) command envm "$command" "$@" ;;
        esac
}`
	fmt.Print(fmt.Sprintf(script, "`", "`"))
	return 0
}

func (cm *InitCommand) Synopsis() string {
	return "$ envm init"
}

type NewCommand struct {
	baseCommand
}

func (cm *NewCommand) Help() string {
	return "new command"
}

func (cm *NewCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	cm.lock()
	defer cm.unLock()
	name := args[0]
	configPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cfg := NewConfig(configPath)
	vs := make(map[string]string)
	for _, key := range args[1:] {
		vs[key] = os.Getenv(key)
	}
	if err := cfg.Save(name, vs); err != nil {
		printError(err)
		return 1
	}
	return 0
}

func (cm *NewCommand) Synopsis() string {
	return "$ envm new <name> [<env-name>…​]"
}

type LSCommand struct {
	baseCommand
}

func (cm *LSCommand) Help() string {
	return "ls command"
}

func (cm *LSCommand) Run(args []string) int {
	confPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cm.lock()
	defer cm.unLock()
	cfg := NewConfig(confPath)
	for _, key := range cfg.NameSpaces() {
		fmt.Println(key)
	}
	return 0
}

func (cm *LSCommand) Synopsis() string {
	return "$ envm ls"
}

type UseCommand struct {
	baseCommand
}

func (cm *UseCommand) Help() string {
	return "use command"
}

func (cm *UseCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	cm.lock()
	defer cm.unLock()

	name := args[0]
	confPath, err := getConfPath()
	if err != nil {
		return 1
	}
	cfg := NewConfig(confPath)
	data, err := cfg.readData()
	if err != nil {
		printError(err)
		return 1
	}
	m, ok := data[name]
	if !ok {
		return 0
	}
	fmt.Print(mapToEnvCommand(m))
	return 0
}

func (cm *UseCommand) Synopsis() string {
	return "$ envm use <name>"
}

type RMCommand struct {
	baseCommand
}

func (cm *RMCommand) Help() string {
	return "rm command"
}

func (cm *RMCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	cm.lock()
	defer cm.unLock()
	name := args[0]
	configPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cfg := NewConfig(configPath)
	path, ok := cfg.getConfPath()
	if !ok {
		return 0
	}
	data, err := cfg.load(path)
	if err != nil {
		printError(err)
		return 1
	}
	if _, ok = data[name]; !ok {
		return 0
	}

	fmt.Printf(`remove %v? [Y/N] `, name)
	if !askYesOrNo(os.Stdin) {
		return 0
	}
	delete(data, name)

	if err = cfg.save(path, data); err != nil {
		printError(err)
		return 1
	}
	return 0
}

func (cm *RMCommand) Synopsis() string {
	return "$ envm rm <name>"
}

type ShowCommand struct {
	baseCommand
}

func (cm *ShowCommand) Help() string {
	return "show command"
}

func (cm *ShowCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	cm.lock()
	defer cm.unLock()
	name := args[0]
	configPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cfg := NewConfig(configPath)
	path, ok := cfg.getConfPath()
	if !ok {
		return 0
	}
	data, err := cfg.load(path)
	if err != nil {
		printError(err)
		return 1
	}
	m, ok := data[name]
	if !ok {
		return 0
	}
	fmt.Println(mapToEnvCommand(m))
	return 0
}

func (cm *ShowCommand) Synopsis() string {
	return "$ envm show <name>"
}

type UpdateCommand struct {
	baseCommand
}

func (cm *UpdateCommand) Help() string {
	return "update command"
}

func (cm *UpdateCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	cm.lock()
	defer cm.unLock()
	name := args[0]
	configPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cfg := NewConfig(configPath)
	path, ok := cfg.getConfPath()
	if !ok {
		return 0
	}
	data, err := cfg.load(path)
	if err != nil {
		printError(err)
		return 1
	}
	vs, ok := data[name]
	if !ok {
		vs = make(map[string]string)
	}
	for _, key := range args[1:] {
		vs[key] = os.Getenv(key)
	}
	if err := cfg.Update(name, vs); err != nil {
		printError(err)
		return 1
	}
	return 0
}

func (cm *UpdateCommand) Synopsis() string {
	return "$ envm update <name>"
}

type CheckCommand struct {
	baseCommand
}

func (cm *CheckCommand) Help() string {
	return "check command"
}

func (cm *CheckCommand) Run(args []string) int {
	if len(args) < 1 {
		fmt.Println(cm.Synopsis())
		return 1
	}
	configPath, err := getConfPath()
	if err != nil {
		printError(err)
		return 1
	}
	cfg := NewConfig(configPath)
	path, ok := cfg.getConfPath()
	if !ok {
		return 0
	}
	data, err := cfg.load(path)
	if err != nil {
		printError(err)
		return 1
	}

	for _, name := range args {
		vs, ok := data[name]
		if ok {
			fmt.Printf("=== Check %v\n", name)
		} else {
			fmt.Printf("=== Not found %v\nFAIL %v", name, name)
			continue
		}
		passed := true
		for k, v := range vs {
			if os.Getenv(k) == v {
				fmt.Printf("✔ %v\n", k)
			} else {
				fmt.Printf("✗ %v\n", k)
				passed = false
			}
		}
		if passed {
			fmt.Printf("ok   %v\n\n", name)
		} else {
			fmt.Printf("FAIL %v\n\n", name)
		}
	}
	return 0
}

func (cm *CheckCommand) Synopsis() string {
	return "$ envm check [<name>…​]"
}
