# Shogi image

## Example

```go
package main

import (
	"image/png"
	"log"
	"os"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/util/image"
)

func main() {
	state := shogi.NewState()
	state.SetBoard(5, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.OU})
	state.SetBoard(5, 3, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
	state.Captured[shogi.TurnBlack].Add(shogi.KI)
	img, err := image.Generate(state, &image.StyleOptions{
		Piece: image.PieceKinki,
	})
	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
```

Output:

![](https://user-images.githubusercontent.com/80381/33804002-910b6c82-dddf-11e7-90c6-87b09871c7c7.png)


## License

[Shogi images by muchonovski](http://mucho.girly.jp/bona/) below `data` directory are under a [Creative Commons 表示-非営利 2.1 日本 License](http://creativecommons.org/licenses/by-nc/2.1/jp/).
