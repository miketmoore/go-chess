package main

import (
	"fmt"
	_ "image/png"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/chess"
	"github.com/miketmoore/pgn"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/image/colornames"
	"golang.org/x/text/language"
)

const screenW = 600
const screenH = 600
const squareSize float64 = 50
const displayFontPath = "assets/kenney_fontpackage/Fonts/Kenney Future Narrow.ttf"
const bodyFontPath = "assets/kenney_fontpackage/Fonts/Kenney Pixel Square.ttf"
const translationFile = "i18n/en.toml"
const lang = "en-US"

func run() {
	// i18n
	bundle := &i18n.Bundle{DefaultLanguage: language.English}

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(translationFile)

	localizer := i18n.NewLocalizer(bundle, "en")

	i18nTitle := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Title",
		},
	})
	i18nPressAnyKey := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "PressAnyKey",
		},
	})

	// Game data
	model := chess.Model{
		PGN:          pgn.PGN{},
		BoardState:   chess.InitialOnBoardState(),
		Draw:         true,
		WhitesMove:   true,
		CurrentState: chess.StateTitle,
	}

	// Setup GUI window
	cfg := pixelgl.WindowConfig{
		Title:  i18nTitle,
		Bounds: pixel.R(0, 0, screenW, screenH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Prepare display text
	displayFace, err := chess.LoadTTF(displayFontPath, 80)
	if err != nil {
		panic(err)
	}

	displayAtlas := text.NewAtlas(displayFace, text.ASCII)
	displayOrig := pixel.V(screenW/2, screenH/2)
	displayTxt := text.New(displayOrig, displayAtlas)

	// Prepare body text
	bodyFace, err := chess.LoadTTF(bodyFontPath, 12)
	if err != nil {
		panic(err)
	}

	// Build body text
	bodyAtlas := text.NewAtlas(bodyFace, text.ASCII)
	bodyOrig := pixel.V(screenW/2, screenH/2)
	bodyTxt := text.New(bodyOrig, bodyAtlas)

	// Title
	titleStr := "Chess"
	fmt.Fprintln(displayTxt, titleStr)

	// Sub-title
	pressAnyKeyStr := i18nPressAnyKey
	fmt.Fprintln(bodyTxt, pressAnyKeyStr)

	// Make board
	boardThemeName := "sandcastle"
	boardW := squareSize * 8
	boardOriginX := (screenW - int(boardW)) / 2
	squares, squareOriginByCoords := chess.NewBoardView(
		float64(boardOriginX),
		150,
		squareSize,
		chess.BoardThemes[boardThemeName]["black"],
		chess.BoardThemes[boardThemeName]["white"],
	)

	// Make pieces
	drawer := chess.NewSpriteByColor()

	validDestinations := []chess.Coord{}

	for !win.Closed() {

		if win.JustPressed(pixelgl.KeyQ) {
			fmt.Println(model.PGN.String())
			os.Exit(0)
		}

		switch model.CurrentState {
		case chess.StateTitle:
			if model.Draw {
				win.Clear(colornames.Black)

				// Draw title text
				c := displayTxt.Bounds().Center()
				heightThird := screenH / 5
				c.Y = c.Y - float64(heightThird)
				displayTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(c)))

				// Draw secondary text
				bodyTxt.Color = colornames.White
				bodyTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(bodyTxt.Bounds().Center())))

				model.Draw = false
			}

			if win.JustPressed(pixelgl.KeyEnter) || win.JustPressed(pixelgl.MouseButtonLeft) {
				model.CurrentState = chess.StateDraw
				win.Clear(colornames.Black)
				model.Draw = true
			}
		case chess.StateDraw:
			if model.Draw {
				draw(win, model.BoardState, drawer, squares)
				model.Draw = false
				model.CurrentState = chess.StateSelectPiece
			}
		case chess.StateSelectPiece:
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				square := chess.FindSquareByVec(squares, win.MousePosition())
				if square != nil {
					coord, ok := chess.GetCoordByXY(
						squareOriginByCoords,
						square.OriginX,
						square.OriginY,
					)
					if ok {
						occupant, isOccupied := model.BoardState[coord]
						if occupant.Color == model.CurrentPlayerColor() && isOccupied {
							validDestinations = chess.GetValidMoves(
								model.CurrentPlayerColor(),
								occupant.Piece,
								model.BoardState,
								coord,
							)
							if len(validDestinations) > 0 {
								model.PieceToMove = occupant
								model.MoveStartCoord = coord
								model.CurrentState = chess.DrawValidMoves
								model.Draw = true
							}
						}

					}

				}
			}
		case chess.DrawValidMoves:
			if model.Draw {
				draw(win, model.BoardState, drawer, squares)
				chess.HighlightSquares(win, squares, validDestinations, colornames.Greenyellow)
				model.Draw = false
				model.CurrentState = chess.StateSelectDestination
			}
		case chess.StateSelectDestination:
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mpos := win.MousePosition()
				square := chess.FindSquareByVec(squares, mpos)
				if square != nil {
					coord, ok := chess.GetCoordByXY(squareOriginByCoords, square.OriginX, square.OriginY)
					if ok {
						occupant, isOccupied := model.BoardState[coord]
						isValid := chess.FindInSliceCoord(validDestinations, coord)
						if isValid && chess.IsDestinationValid(model.WhitesMove, isOccupied, occupant) {
							move(&model, coord)

							inCheckData := chess.GetInCheckData(model.BoardState)
							if inCheckData.InCheck {
								fmt.Println("Someone is in check!")
								// TODO add check notation to PGN data
								if len(inCheckData.WhiteThreateningBlack) > 0 {
									fmt.Println("Black is in check by white")
								}
								if len(inCheckData.BlackThreateningWhite) > 0 {
									fmt.Println("White is in check by black")
								}
							}
							fmt.Println(model.PGN.String())
						} else {
							model.CurrentState = chess.StateSelectPiece
						}
					}
				}
			}
		}

		win.Update()
	}
}

func move(model *chess.Model, destCoord chess.Coord) {
	model.CurrentState = chess.StateDraw
	model.Draw = true
	model.MoveDestinationCoord = destCoord

	// _, isCapture := model.BoardState[destCoord]

	if model.WhitesMove {
		model.PGN.Movetext = append(model.PGN.Movetext, pgn.MovetextEntry{
			White: fmt.Sprintf("%s%s", chess.FileToFileView[destCoord.File], chess.RankToRankView[destCoord.Rank]),
		})
	} else {
		move := model.PGN.Movetext[len(model.PGN.Movetext)-1]
		move.Black = fmt.Sprintf("%s%s", chess.FileToFileView[destCoord.File], chess.RankToRankView[destCoord.Rank])
		model.PGN.Movetext[len(model.PGN.Movetext)-1] = move
	}

	model.BoardState[destCoord] = model.PieceToMove
	delete(model.BoardState, model.MoveStartCoord)
	model.WhitesMove = !model.WhitesMove
}

func main() {
	pixelgl.Run(run)
}

func draw(win *pixelgl.Window, boardState chess.BoardState, drawer chess.Drawer, squares chess.BoardMap) {
	// Draw board
	for _, square := range squares {
		square.Shape.Draw(win)
	}

	// Draw pieces in the correct position
	for coord, livePieceData := range boardState {
		var set chess.PieceSpriteSet
		if livePieceData.Color == chess.PlayerBlack {
			set = drawer.Black
		} else {
			set = drawer.White
		}

		var piece *pixel.Sprite
		switch livePieceData.Piece {
		case chess.Bishop:
			piece = set.Bishop
		case chess.King:
			piece = set.King
		case chess.Knight:
			piece = set.Knight
		case chess.Pawn:
			piece = set.Pawn
		case chess.Queen:
			piece = set.Queen
		case chess.Rook:
			piece = set.Rook
		}

		chess.DrawPiece(win, squares, piece, coord)
	}
}
