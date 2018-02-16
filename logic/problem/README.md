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
	"github.com/sugyan/shogi/record"
)

func main() {
	// generate
	state, _ := generator.Generate(generator.Type3)
	// solve
	answer := solver.Solve(state)

	record := &record.Record{
		State: state,
		Moves: answer,
	}
	println(record.ConvertToString(csa.NewConverter(nil)))

	answerStr, err := state.MoveStrings(answer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Answer: %s\n", strings.Join(answerStr, " -> "))
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
+
+2322TO
-1213OU
+0023HI

Answer: ▲2二歩成 -> △1三玉 -> ▲2三飛
```
