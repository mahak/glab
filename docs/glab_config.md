## glab config

Set and get glab settings

### Synopsis

Get and set key/value strings.

		Current respected settings:

		- token: Your gitlab access token, defaults to environment variables
		- gitlab_uri: if unset, defaults to https://gitlab.com
		- browser: if unset, defaults to environment variables
		- editor: if unset, defaults to environment variables.
		- visual: alternative for editor. if unset, defaults to environment variables.
		- glamour_style: Your desired markdown renderer style. Options are dark, light, notty. Custom styles are allowed set a custom style
https://github.com/charmbracelet/glamour#styles
	

### Options

```
  -g, --global   use global config file
```

### Options inherited from parent commands

```
      --help   Show help for command
```

### SEE ALSO

* [glab](glab.md)	 - A GitLab CLI Tool
* [glab config get](glab_config_get.md)	 - Prints the value of a given configuration key
* [glab config init](glab_config_init.md)	 - Shows a prompt to set basic glab configuration
* [glab config set](glab_config_set.md)	 - Updates configuration with the value of a given key

###### Auto generated by spf13/cobra on 4-Oct-2020
