package image

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/sugyan/shogi"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

// Convert function
func Convert(state *shogi.State) (image.Image, error) {
	boardImg, err := loadImage("data/board.png")
	if err != nil {
		return nil, err
	}

	xStep := 684.0 / 9.0
	yStep := 751.0 / 9.0
	xOffset := xStep * 3.0
	dst := image.NewRGBA(image.Rectangle{Min: image.ZP, Max: boardImg.Bounds().Size().Add(image.Pt(int(xStep*6), 0))})
	// board
	draw.Draw(dst, dst.Bounds().Add(image.Pt(int(xOffset), 0)), boardImg, boardImg.Bounds().Min, draw.Over)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			p := state.Board[i][j]
			if p != nil {
				pieceImg, err := loadImage(fmt.Sprintf("data/%s.png", p.Code()))
				if err != nil {
					return nil, err
				}
				if p.Turn() == shogi.TurnSecond {
					pieceImg = rotate180(pieceImg)
				}
				r := dst.Bounds().
					Add(image.Pt(int(xOffset)+33, 28)).
					Add(image.Pt(int(xStep*float64(j)), int(yStep*float64(i))))
				draw.Draw(dst, r, pieceImg, pieceImg.Bounds().Min, draw.Over)
			}
		}
	}
	// captured
	for turn, captured := range state.Captured {
		pieces := arrangeCapturedPieces(turn, captured)
		var offset image.Point
		switch turn {
		case shogi.TurnFirst:
			offset = image.Pt(
				int(xOffset)+boardImg.Bounds().Dx()+5,
				int(yStep*float64(-4))+boardImg.Bounds().Dy(),
			)
		case shogi.TurnSecond:
			offset = image.Pt(5, 0)
		}
		for i := 0; i < 4; i++ {
			for j := 0; j < 2; j++ {
				k := i*2 + j
				if k < len(pieces) {
					var data *capturedPiecesData
					switch turn {
					case shogi.TurnFirst:
						data = pieces[k]
					case shogi.TurnSecond:
						data = pieces[len(pieces)-k-1]
					}
					pieceImg, err := loadImage(fmt.Sprintf("data/%s.png", data.piece.Code()))
					if err != nil {
						return nil, err
					}
					if turn == shogi.TurnSecond {
						pieceImg = rotate180(pieceImg)
					}
					r := dst.Bounds().
						Add(offset).
						Add(image.Pt(
							int(xStep*1.5*float64(j)),
							int(yStep*float64(i))))
					draw.Draw(dst, r, pieceImg, pieceImg.Bounds().Min, draw.Over)
					if data.num > 1 {
						// draw "x2" etc.
					}
				}
			}
		}
	}

	return dst, nil
}

func loadImage(assetName string) (image.Image, error) {
	data, err := Asset(assetName)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewBuffer(data))
}

func rotate180(img image.Image) image.Image {
	dst := image.NewRGBA(img.Bounds())
	draw.NearestNeighbor.Transform(dst,
		f64.Aff3{
			-1.0, 0.0, float64(img.Bounds().Dx()),
			0.0, -1.0, float64(img.Bounds().Dy()),
		}, img, img.Bounds(), draw.Over, nil)
	return dst
}

type capturedPiecesData struct {
	piece shogi.Piece
	num   int
}

func arrangeCapturedPieces(turn shogi.Turn, cp *shogi.CapturedPieces) []*capturedPiecesData {
	results := []*capturedPiecesData{}
	if cp.Hi > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.HI),
			num:   cp.Hi,
		})
	}
	if cp.Ka > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.KA),
			num:   cp.Ka,
		})
	}
	if cp.Ki > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.KI),
			num:   cp.Ki,
		})
	}
	if cp.Gi > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.GI),
			num:   cp.Gi,
		})
	}
	if cp.Ke > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.KE),
			num:   cp.Ke,
		})
	}
	if cp.Ky > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.KY),
			num:   cp.Ky,
		})
	}
	if cp.Fu > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(turn, shogi.FU),
			num:   cp.Fu,
		})
	}
	return results
}
