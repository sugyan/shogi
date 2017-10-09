package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // for decoding jpeg
	_ "image/png"  // for decoding png

	"github.com/sugyan/shogi"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

// Generate function
func Generate(state *shogi.State) (image.Image, error) {
	boardImg, err := loadImage("data/board.jpg")
	if err != nil {
		return nil, err
	}

	xStep := 540.0 / 9.0
	yStep := 576.0 / 9.0
	xOffset := xStep * 3.0
	dst := image.NewRGBA(image.Rectangle{Min: image.ZP, Max: boardImg.Bounds().Size().Add(image.Pt(int(xStep*6), 0))})
	for i := 0; i < dst.Bounds().Dx(); i++ {
		for j := 0; j < dst.Bounds().Dy(); j++ {
			dst.Set(i, j, color.White)
		}
	}
	// board
	draw.Draw(dst, dst.Bounds().Add(image.Pt(int(xOffset), 0)), boardImg, boardImg.Bounds().Min, draw.Over)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := state.Board[i][j]
			if bp != nil {
				pieceImg, err := loadPieceImage(bp.Turn, bp.Piece)
				if err != nil {
					return nil, err
				}
				r := dst.Bounds().
					Add(image.Pt(int(xOffset)+30, 30)).
					Add(image.Pt(int(xStep*float64(j)), int(yStep*float64(i))))
				draw.Draw(dst, r, pieceImg, pieceImg.Bounds().Min, draw.Over)
			}
		}
	}
	// captured
	for turn, captured := range state.Captured {
		pieces := arrangeCapturedPieces(captured)
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
					pieceImg, err := loadPieceImage(turn, data.piece)
					if err != nil {
						return nil, err
					}
					r := dst.Bounds().
						Add(offset).
						Add(image.Pt(
							int(xStep*1.5*float64(j)),
							int(yStep*float64(i))))
					draw.Draw(dst, r, pieceImg, pieceImg.Bounds().Min, draw.Over)
					if data.num > 1 {
						o := r.Bounds().Min.Add(pieceImg.Bounds().Max).Add(image.Pt(0, -5))
						drawer := &font.Drawer{
							Dst:  dst,
							Src:  image.Black,
							Face: inconsolata.Bold8x16,
							Dot: fixed.Point26_6{
								X: fixed.Int26_6(o.X << 6),
								Y: fixed.Int26_6(o.Y << 6),
							},
						}
						drawer.DrawString(fmt.Sprintf("x%d", data.num))
					}
				}
			}
		}
	}

	return dst, nil
}

func loadPieceImage(turn shogi.Turn, piece *shogi.Piece) (image.Image, error) {
	code := piece.Code().String()
	if piece.Code() == shogi.OU && turn == shogi.TurnSecond {
		code += "_"
	}
	img, err := loadImage(fmt.Sprintf("data/%s.png", code))
	if err != nil {
		return nil, err
	}
	if turn == shogi.TurnSecond {
		img = rotate180(img)
	}
	return img, nil
}

func loadImage(assetName string) (image.Image, error) {
	data, err := Asset(assetName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return img, nil
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
	piece *shogi.Piece
	num   int
}

func arrangeCapturedPieces(cp *shogi.CapturedPieces) []*capturedPiecesData {
	results := []*capturedPiecesData{}
	if cp.HI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.HI),
			num:   cp.HI,
		})
	}
	if cp.KA > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.KA),
			num:   cp.KA,
		})
	}
	if cp.KI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.KI),
			num:   cp.KI,
		})
	}
	if cp.GI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.GI),
			num:   cp.GI,
		})
	}
	if cp.KE > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.KE),
			num:   cp.KE,
		})
	}
	if cp.KY > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.KY),
			num:   cp.KY,
		})
	}
	if cp.FU > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.NewPiece(shogi.FU),
			num:   cp.FU,
		})
	}
	return results
}
