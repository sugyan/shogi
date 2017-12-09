package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"time"

	"github.com/sugyan/shogi"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type (
	// BoardStyle type
	BoardStyle int
	// GridStyle type
	GridStyle int
	// PieceStyle type
	PieceStyle int
)

// board image styles
const (
	BoardDefault BoardStyle = iota
	BoardRandom
	BoardKayaA
	BoardKayaB
	BoardKayaC
	BoardKayaD
	BoardPlywood
	BoardFoldaway
	BoardWellused
	BoardPaper
	BoardStripe
	BoardPlain
)

// grid image styles
const (
	GridDefault GridStyle = iota
	GridRandom
	GridDot
	GridDotXY
	GridNoDot
	GridNoDotXY
	GridHandWritten
)

// piece image styles
const (
	PieceDefault PieceStyle = iota
	PieceRandom
	PieceKinki
	PieceKinkiTorafu
	PieceRyoko
	PieceRyokoTorafu
	PieceDirty
)

// StyleOptions type
type StyleOptions struct {
	Board     BoardStyle
	Grid      GridStyle
	Piece     PieceStyle
	HighLight *shogi.Position
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate function
func Generate(state *shogi.State, options *StyleOptions) (image.Image, error) {
	options = checkOptions(options)
	boardImg, err := boardImage(options.Board)
	if err != nil {
		return nil, err
	}
	gridImg, err := gridImage(options.Grid)
	if err != nil {
		return nil, err
	}

	xStep := 387.0 / 9.0
	yStep := 432.0 / 9.0
	xOffset := 138
	dst := image.NewRGBA(image.Rectangle{Min: image.ZP, Max: boardImg.Bounds().Size().Add(image.Pt(276, 0))})
	for i := 0; i < dst.Bounds().Dx(); i++ {
		for j := 0; j < dst.Bounds().Dy(); j++ {
			dst.Set(i, j, color.White)
		}
	}
	// board
	draw.Draw(dst, dst.Bounds().Add(image.Pt(xOffset, 0)), boardImg, boardImg.Bounds().Min, draw.Over)
	draw.Draw(dst, dst.Bounds().Add(image.Pt(xOffset, 0)), gridImg, gridImg.Bounds().Min, draw.Over)
	if options.HighLight != nil {
		focusImage, err := loadImage("data/focus/focus_bold_o.png")
		i, j := options.HighLight.Rank-1, 9-options.HighLight.File
		if err != nil {
			return nil, err
		}
		r := dst.Bounds().
			Add(image.Pt(xOffset+11, 11)).
			Add(image.Pt(int(xStep*float64(j)), int(yStep*float64(i))))
		draw.Draw(dst, r, focusImage, focusImage.Bounds().Min, draw.Over)

	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := state.Board[i][j]
			if b != nil {
				pieceImg, err := pieceImage(b.Turn, b.Piece, options.Piece)
				if err != nil {
					return nil, err
				}
				r := dst.Bounds().
					Add(image.Pt(xOffset+11, 11)).
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
		case shogi.TurnBlack:
			offset = image.Pt(
				xOffset+boardImg.Bounds().Dx()+2,
				int(yStep*float64(-4))+boardImg.Bounds().Dy(),
			)
		case shogi.TurnWhite:
			offset = image.Pt(2, 0)
		}
		for i := 0; i < 4; i++ {
			for j := 0; j < 2; j++ {
				k := i*2 + j
				if k < len(pieces) {
					var data *capturedPiecesData
					switch turn {
					case shogi.TurnBlack:
						data = pieces[k]
					case shogi.TurnWhite:
						data = pieces[len(pieces)-k-1]
					}
					pieceImg, err := pieceImage(turn, data.piece, options.Piece)
					if err != nil {
						return nil, err
					}
					r := dst.Bounds().
						Add(offset).
						Add(image.Pt(
							int(66*float64(j)),
							int(yStep*float64(i))))
					draw.Draw(dst, r, pieceImg, pieceImg.Bounds().Min, draw.Over)
					if data.num > 1 {
						o := r.Bounds().Min.Add(pieceImg.Bounds().Max).Add(image.Pt(-1, -5))
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

func checkOptions(so *StyleOptions) *StyleOptions {
	if so == nil {
		so = &StyleOptions{}
	} else {
		if so.Board == BoardRandom {
			so.Board = []BoardStyle{
				BoardKayaA,
				BoardKayaB,
				BoardKayaC,
				BoardKayaD,
				BoardPlywood,
				BoardFoldaway,
				BoardWellused,
				BoardPaper,
				BoardStripe,
				BoardPlain,
			}[rand.Intn(10)]
		}
		if so.Grid == GridRandom {
			so.Grid = []GridStyle{
				GridDot,
				GridDotXY,
				GridNoDot,
				GridNoDotXY,
				GridHandWritten,
			}[rand.Intn(5)]
		}
		if so.Piece == PieceRandom {
			so.Piece = []PieceStyle{
				PieceKinki,
				PieceKinkiTorafu,
				PieceRyoko,
				PieceRyokoTorafu,
				PieceDirty,
			}[rand.Intn(5)]
		}
	}
	return so
}

func boardImage(style BoardStyle) (image.Image, error) {
	var assetName string
	switch style {
	case BoardKayaA:
		assetName = "data/ban/ban_kaya_a.png"
	case BoardKayaB:
		assetName = "data/ban/ban_kaya_b.png"
	case BoardKayaC:
		assetName = "data/ban/ban_kaya_c.png"
	case BoardKayaD:
		assetName = "data/ban/ban_kaya_d.png"
	case BoardPlywood:
		assetName = "data/ban/ban_gohan.png"
	case BoardFoldaway:
		assetName = "data/ban/ban_oritatami.png"
	case BoardWellused:
		assetName = "data/ban/ban_dirty.png"
	case BoardPaper:
		assetName = "data/ban/ban_paper.png"
	case BoardStripe:
		assetName = "data/ban/ban_stripe.png"
	case BoardPlain:
		assetName = "data/ban/ban_muji.png"
	default:
		assetName = "data/ban/ban_kaya_a.png"
	}
	return loadImage(assetName)
}

func gridImage(style GridStyle) (image.Image, error) {
	var assetName string
	switch style {
	case GridDot:
		assetName = "data/masu/masu_dot.png"
	case GridDotXY:
		assetName = "data/masu/masu_dot_xy.png"
	case GridNoDot:
		assetName = "data/masu/masu_nodot.png"
	case GridNoDotXY:
		assetName = "data/masu/masu_nodot_xy.png"
	case GridHandWritten:
		assetName = "data/masu/masu_handwriting.png"
	default:
		assetName = "data/masu/masu_dot_xy.png"
	}
	return loadImage(assetName)
}

func pieceImage(turn shogi.Turn, piece shogi.Piece, style PieceStyle) (image.Image, error) {
	var piecesDirName string
	switch style {
	case PieceKinki:
		piecesDirName = "data/koma_kinki"
	case PieceKinkiTorafu:
		piecesDirName = "data/koma_kinki_torafu"
	case PieceRyoko:
		piecesDirName = "data/koma_ryoko"
	case PieceRyokoTorafu:
		piecesDirName = "data/koma_ryoko_torafu"
	case PieceDirty:
		piecesDirName = "data/koma_dirty"
	default:
		piecesDirName = "data/koma_ryoko"
	}

	prefix := "+"
	if turn == shogi.TurnWhite {
		prefix = "-"
	}
	code := piece.String()
	img, err := loadImage(fmt.Sprintf("%s/%s%s.png", piecesDirName, prefix, code))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadImage(assetName string) (image.Image, error) {
	data, err := Asset(assetName)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewBuffer(data))
}

type capturedPiecesData struct {
	piece shogi.Piece
	num   int
}

func arrangeCapturedPieces(cp *shogi.CapturedPieces) []*capturedPiecesData {
	results := []*capturedPiecesData{}
	if cp.HI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.HI,
			num:   cp.HI,
		})
	}
	if cp.KA > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.KA,
			num:   cp.KA,
		})
	}
	if cp.KI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.KI,
			num:   cp.KI,
		})
	}
	if cp.GI > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.GI,
			num:   cp.GI,
		})
	}
	if cp.KE > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.KE,
			num:   cp.KE,
		})
	}
	if cp.KY > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.KY,
			num:   cp.KY,
		})
	}
	if cp.FU > 0 {
		results = append(results, &capturedPiecesData{
			piece: shogi.FU,
			num:   cp.FU,
		})
	}
	return results
}
