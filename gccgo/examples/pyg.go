package pyg

import (
	"os"
	"fmt"
	"gopy"
)

func init() {
	fmt.Printf("pyg.init()\n")
}

func Main() {
	py.InitializeEx(false)
	py.Main(os.Args)
}
