// package repl provides a Read Eval Print Loop functionality for the RQL DB
package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/spencercdixon/rql/rql"
)

const PromptTemplate = "rql(%s)=# "

func Start(db *rql.Database, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PromptTemplate, db.Path)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		if strings.HasPrefix(line, "\\") {
			handleMetaCommand(line)
		} else {
			fmt.Println("Executing")
		}
	}
}

func handleMetaCommand(line string) {
	if line == `\q` {
		fmt.Println("Exiting...")
	} else if line == `\?` {
		printHelp()
	}
}

func printHelp() {
	fmt.Println("\n\rrql help:")
	fmt.Println(` \q - quit `)
	fmt.Println(` \? - help `)
	println()
}
