package commands

import (
	"fmt"

	aliasCmd "github.com/profclems/glab/commands/alias"
	"github.com/profclems/glab/commands/cmdutils"
	completionCmd "github.com/profclems/glab/commands/completion"
	configCmd "github.com/profclems/glab/commands/config"
	"github.com/profclems/glab/commands/help"
	issueCmd "github.com/profclems/glab/commands/issue"
	labelCmd "github.com/profclems/glab/commands/label"
	mrCmd "github.com/profclems/glab/commands/mr"
	pipelineCmd "github.com/profclems/glab/commands/pipeline"
	projectCmd "github.com/profclems/glab/commands/project"
	releaseCmd "github.com/profclems/glab/commands/release"
	updateCmd "github.com/profclems/glab/commands/update"
	versionCmd "github.com/profclems/glab/commands/version"
	"github.com/profclems/glab/internal/glrepo"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RootCmd is the main root/parent command
func NewCmdRoot(f *cmdutils.Factory, version, buildDate string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:           "glab <command> <subcommand> [flags]",
		Short:         "A GitLab CLI Tool",
		Long:          `GLab is an open source Gitlab Cli tool bringing GitLab to your command line`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
	$ glab issue create
	$ glab mr list --all
	$ glab mr checkout 123
	$ glab pipeline ci view
	`),
		Annotations: map[string]string{
			"help:environment": heredoc.Doc(`
			GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
			prompted to authenticate and overrides any previously stored credentials.
			Can be set in the config with 'glab config set token xxxxxx'

			GITLAB_URI or GITLAB_HOST: specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com). Default is https://gitlab.com.

			REMOTE_ALIAS or GIT_REMOTE_URL_VAR: git remote variable or alias that contains the gitlab url.
			Can be set in the config with 'glab config set remote_alias origin'

			VISUAL, EDITOR (in order of precedence): the editor tool to use for authoring text.
			Can be set in the config with 'glab config set editor vim'

			BROWSER: the web browser to use for opening links.
			Can be set in the config with 'glab config set browser mybrowser'

			GLAMOUR_STYLE: environment variable to set your desired markdown renderer style
			Available options are (dark|light|notty) or set a custom style
			https://github.com/charmbracelet/glamour#styles
		`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				fmt.Printf("Unknown command: %s\n", args[0])
				return cmd.Usage()
			} else if ok, _ := cmd.Flags().GetBool("version"); ok {
				return versionCmd.NewCmdVersion(version, buildDate).RunE(cmd, args)
			}

			return cmd.Help()
		},
	}

	rootCmd.PersistentFlags().Bool("help", false, "Show help for command")
	rootCmd.SetHelpFunc(help.RootHelpFunc)
	rootCmd.SetUsageFunc(help.RootUsageFunc)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err == pflag.ErrHelp {
			return err
		}
		return &cmdutils.FlagError{Err: err}
	})

	// Child commands
	rootCmd.AddCommand(aliasCmd.NewCmdAlias(f))
	rootCmd.AddCommand(configCmd.NewCmdConfig(f))
	rootCmd.AddCommand(completionCmd.NewCmdCompletion())
	rootCmd.AddCommand(versionCmd.NewCmdVersion(version, buildDate))

	// below here at the commands that require the "intelligent" BaseRepo resolver
	repoResolvingCmdFactory := *f
	repoResolvingCmdFactory.BaseRepo = resolvedBaseRepo(f)
	rootCmd.AddCommand(issueCmd.NewCmdIssue(&repoResolvingCmdFactory))
	rootCmd.AddCommand(labelCmd.NewCmdLabel(&repoResolvingCmdFactory))
	rootCmd.AddCommand(mrCmd.NewCmdMR(&repoResolvingCmdFactory))
	rootCmd.AddCommand(pipelineCmd.NewCmdPipeline(&repoResolvingCmdFactory))
	rootCmd.AddCommand(projectCmd.NewCmdRepo(&repoResolvingCmdFactory))
	rootCmd.AddCommand(releaseCmd.NewCmdRelease(&repoResolvingCmdFactory))
	rootCmd.Flags().BoolP("version", "v", false, "show glab version information")
	rootCmd.AddCommand(updateCmd.NewCheckUpdateCmd(version, buildDate))
	return rootCmd
}

//
func resolvedBaseRepo(f *cmdutils.Factory) func() (glrepo.Interface, error) {
	return func() (glrepo.Interface, error) {
		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}
		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}
		repoContext, err := glrepo.ResolveRemotesToRepos(remotes, httpClient, "")
		if err != nil {
			return nil, err
		}
		baseRepo, err := repoContext.BaseRepo(true)
		if err != nil {
			return nil, err
		}

		return baseRepo, nil
	}
}
