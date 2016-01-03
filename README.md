# envm

Environment variables management made easy.

You can easily management various API token for each environment.

# Installation

## Setup

1. go get!

	```sh
	$ go get -u github.com/bluele/envm
	```

2. Add envm init to your shell to enable commands and autocompletion.

	```sh
	$ echo 'eval "$(envm init -)"' >> ~/.bash_profile
	```

	Same as in previous step, use ~/.bashrc on Ubuntu, or ~/.zshrc for Zsh.

3. Restart your shell so that PATH changes take effect. (Opening a new
   terminal tab will usually do it.) Now check if envm was set up:

    ```sh
    $ type envm
    #=> "envm is a function"
    ```

# Usage

```sh
# Create a new environment variable set.
$ OPTION=1 envm new app.dev API_TOKEN API_SECRET OPTION

$ envm new app.prod API_TOKEN API_SECRET

# List of environ variable sets
$ envm ls
app.dev
app.prod

# Show environment variables of specified dataset
$ envm show app.dev
export API_TOKEN="xxxxxx"
export API_SECRET="xxxxxx"
export OPTION="1"


# Other terminal session
# 
$ envm use app.dev

$ echo $API_TOKEN
xxxxxx

# Update the dataset with specified variables
$ OPTION=0 envm update app.dev OPTION

$ envm show app.dev
export API_TOKEN="xxxxxx"
export API_SECRET="xxxxxx"
export OPTION="0"

# Remove specified dataset
$ envm rm app.dev
remove app.dev? [Y/N] Y

$ envm ls
app.prod
```

# ENVIRONMENT VARIABLES

*ENVM_HOME*

	`envm` create a config file at your homedir. (~/.envm.yml)
	If you custome this path, you can change this with environment variable `ENVM_HOME`.


# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>