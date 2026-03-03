package cli

import "fmt"

// overriden during build time via ldflags
var version = "dev"

const helpCmdStr = `usage: silo [command]

commands:
  (none)      open the editor
  wizard      re-run the setup wizard
  pwd         print the notes directory
  changedir   change the notes directory
  reset       delete config and start fresh
  version     print version
  help        show this message`

func cmdPwd() {}

func cmdVersion() {}

func cmdChangeDir() {}

func cmdReset() {}

func cmdHelp() {
	fmt.Println(helpCmdStr)
}
