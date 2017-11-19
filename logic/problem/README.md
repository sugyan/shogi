# 詰将棋 (Japanese chess problem)

## Example

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sugyan/shogi/format/csa"
	"github.com/sugyan/shogi/logic/problem/generator"
	"github.com/sugyan/shogi/logic/problem/solver"
)

func main() {
	// generate
	state := generator.Generate(generator.ProblemType3)
	println(csa.InitialState1(state))
	// solve
	answer, err := solver.Solve(state)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer: %s\n", strings.Join(answer, " -> "))
}
```

Output:

```
P1 *  *  *  *  *  * +GI *  * 
P2 *  *  *  *  *  *  *  * -OU
P3 *  *  *  *  *  *  * +FU * 
P4 *  *  *  *  *  *  *  * -KE
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00HI
P-00AL

Answer: ▲2二歩成 -> △1三玉 -> ▲2三飛
```