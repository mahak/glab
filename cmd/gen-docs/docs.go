package main

import (
	"github.com/profclems/glab/commands"
	"log"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(commands.NewCmdRoot(&cmdutils.Factory{}, "", ""), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
