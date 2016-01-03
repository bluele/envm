package main

import (
	"errors"
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/nightlyone/lockfile"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	appName = "envm"
	version = "0.1.0"

	filename = ".envm.yml"
	tempname = filename + ".bak"

	lockfileName = "envm.lck"

	envmHomePath = "ENVM_HOME"
)

var (
	NotFoundConfigFile = errors.New("Not found config file.")
	flock              *lockfile.Lockfile
)

type Store map[string]map[string]string

type Config struct {
	basePath string
}

func NewConfig(basePath string) *Config {
	cfg := &Config{basePath: basePath}
	return cfg
}

func (cfg *Config) getConfPath() (string, bool) {
	path := filepath.Join(cfg.basePath, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, false
	}
	return path, true
}

func (cfg *Config) getTempPath() string {
	return filepath.Join(cfg.basePath, tempname)
}

func (cfg *Config) load(path string) (map[string]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	data := make(map[string]map[string]string)
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cfg *Config) Save(name string, vs map[string]string) error {
	path, ok := cfg.getConfPath()
	if ok {
		data, err := cfg.load(path)
		if err != nil {
			return err
		}
		if _, ok = data[name]; ok {
			return fmt.Errorf("Already exists key `%v`", name)
		}
		data[name] = vs
		return cfg.save(path, data)
	}
	return cfg.save(path, map[string]map[string]string{name: vs})
}

func (cfg *Config) Update(name string, vs map[string]string) error {
	path, ok := cfg.getConfPath()
	if ok {
		data, err := cfg.load(path)
		if err != nil {
			return err
		}
		data[name] = mergeMap(data[name], vs)
		return cfg.save(path, data)
	}
	return cfg.save(path, map[string]map[string]string{name: vs})
}

func (cfg *Config) save(path string, vs map[string]map[string]string) error {
	temppath := cfg.getTempPath()
	err := cfg.write(temppath, vs)
	if err != nil {
		return err
	}
	return os.Rename(temppath, path)
}

func (cfg *Config) write(path string, vs map[string]map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := yaml.Marshal(vs)
	if err != nil {
		return err
	}
	f.Write(content)
	return nil
}

func (cfg *Config) readData() (map[string]map[string]string, error) {
	path, ok := cfg.getConfPath()
	if !ok {
		return nil, NotFoundConfigFile
	}
	return cfg.load(path)
}

func (cfg *Config) NameSpaces() []string {
	data, err := cfg.readData()
	if err != nil {
		return nil
	}
	return sortKeys(data)
}

func mapToEnvCommand(m map[string]string) string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	commands := []string{}
	for _, k := range keys {
		commands = append(commands, fmt.Sprintf("export %v=%#v", k, m[k]))
	}
	return strings.Join(commands, "\n")
}

func printError(err error) {
	fmt.Println(err.Error())
}

/* default: path to home directory
   env: ENVM_HOME: user defined variable
*/
func getConfPath() (string, error) {
	v, err := getEnv(envmHomePath)
	if err == nil {
		return v, nil
	}
	return homedir.Dir()
}

func getLock() (*lockfile.Lockfile, error) {
	if flock == nil {
		lock, err := lockfile.New(filepath.Join(os.TempDir(), lockfileName))
		if err != nil {
			return nil, err
		}
		flock = &lock
	}
	return flock, nil
}

func mustGetLock() *lockfile.Lockfile {
	lock, err := getLock()
	if err != nil {
		panic(err)
	}
	return lock
}

func main() {
	c := cli.NewCLI(appName, version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &InitCommand{}, nil
		},
		"ls": func() (cli.Command, error) {
			return &LSCommand{}, nil
		},
		"new": func() (cli.Command, error) {
			return &NewCommand{}, nil
		},
		"use": func() (cli.Command, error) {
			return &UseCommand{}, nil
		},
		"rm": func() (cli.Command, error) {
			return &RMCommand{}, nil
		},
		"show": func() (cli.Command, error) {
			return &ShowCommand{}, nil
		},
		"update": func() (cli.Command, error) {
			return &UpdateCommand{}, nil
		},
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(status)
}
