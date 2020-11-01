package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// ExecName returns the basename of the executable
func ExecName() string {
	// const progname string = "central-server"
	progname, _ := os.Executable()
	progname = path.Base(progname)
	return progname
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: fmt.Sprintf(`To load completions:

Bash:

$ source <(%[1]s completion bash)

# To load completions for each session, execute once:
Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
MacOS:
  $ %[1]s completion bash > /usr/local/etc/bash_completion.d/%[1]s

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

# You will need to start a new shell for this setup to take effect.

Fish:

$ %[1]s completion fish | source

# To load completions for each session, execute once:
$ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish
`, ExecName()),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
