package main
import (
"fmt"
"os"
)

func main() {
  args := os.Args
for _, arg := range args {
 fmt.Println(arg)
 }
panic("error")
}
