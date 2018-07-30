package chess

import (
	"fmt"
)

func isCoordStartPosition(playerColor PlayerColor, piece Piece, rank Rank) bool {

	if playerColor == PlayerWhite {
		// white
		if piece == Pawn {
			return rank == Rank2
		}
	} else {
		// black
		if piece == Pawn {
			return rank == Rank7
		}
	}

	return false
}

// Direction is a custom type used to describe a direction from a square on the board
type Direction int

const (
	North Direction = iota
	NorthWest
	NorthEast
	South
	SouthWest
	SouthEast
	East
	West
)

type DirectionSlope []int

var (
	SlopeNorth     DirectionSlope = []int{1, 0}
	SlopeNorthEast DirectionSlope = []int{1, 1}
	SlopeEast      DirectionSlope = []int{0, 1}
	SlopeSouthEast DirectionSlope = []int{-1, 1}
	SlopeSouth     DirectionSlope = []int{-1, 0}
	SlopeSouthWest DirectionSlope = []int{-1, -1}
	SlopeWest      DirectionSlope = []int{0, -1}
	SlopeNorthWest DirectionSlope = []int{1, -1}
)

// Rank is a custom type that represents a horizontal row (rank) on the chess board
type Rank int

// File is a custom type that represents a vertical column (file) on the chess board
type File int

const (
	RankNone Rank = 0
	Rank1    Rank = 1
	Rank2    Rank = 2
	Rank3    Rank = 3
	Rank4    Rank = 4
	Rank5    Rank = 5
	Rank6    Rank = 6
	Rank7    Rank = 7
	Rank8    Rank = 8

	FileNone File = 0
	FileA    File = 1
	FileB    File = 2
	FileC    File = 3
	FileD    File = 4
	FileE    File = 5
	FileF    File = 6
	FileG    File = 7
	FileH    File = 8
)

var ranksOrder = []Rank{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8}
var FilesOrder = []File{FileA, FileB, FileC, FileD, FileE, FileF, FileG, FileH}

var RankToRankView = map[Rank]string{
	Rank1: "1",
	Rank2: "2",
	Rank3: "3",
	Rank4: "4",
	Rank5: "5",
	Rank6: "6",
	Rank7: "7",
	Rank8: "8",
}

var FileToFileView = map[File]string{
	FileA: "a",
	FileB: "b",
	FileC: "c",
	FileD: "d",
	FileE: "e",
	FileF: "f",
	FileG: "g",
	FileH: "h",
}

// Coord represents a board square by rank and file
type Coord struct {
	Rank Rank
	File File
}

// GetRankFile returns the rank and file
func (c Coord) GetRankFile() (Rank, File) {
	return c.Rank, c.File
}

// NewCoord returns a new rank and file coordinate
func NewCoord(file File, rank Rank) Coord {
	return Coord{
		File: file,
		Rank: rank,
	}
}

// GetValidMoves returns a list of valid coordinates the piece can be moved to
func GetValidMoves(playerColor PlayerColor, piece Piece, boardState BoardState, coord Coord) []Coord {
	switch piece {
	case Pawn:
		return getValidMovesPawn(playerColor, boardState, coord)
	case King:
		return getValidMovesKing(playerColor, boardState, coord)
	case Knight:
		return getValidMovesKnight(playerColor, boardState, coord)
	case Rook:
		return getValidMovesRook(playerColor, boardState, coord)
	case Bishop:
		return getValidMovesForBishop(playerColor, boardState, coord)
	case Queen:
		return getValidMovesForQueen(playerColor, boardState, coord)
	}
	return []Coord{}
}

func getValidMovesPawn(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {
	enemyColor := GetOppositeColor(playerColor)
	valid := []Coord{}

	yChange := 1
	if playerColor == PlayerBlack {
		yChange = -1
	}

	// get two spaces north or south
	coords := GetCoordsBySlopeAndDistance(currCoord, yChange, 0, 2)
	if !isOccupied(boardState, coords[0]) {
		valid = append(valid, coords[0])
		if isCoordStartPosition(playerColor, Pawn, currCoord.Rank) && !isOccupied(boardState, coords[1]) {
			valid = append(valid, coords[1])
		}
	}

	// pawn attack moves
	if playerColor == PlayerWhite {
		// NW
		coord, ok := GetCoordBySlopeAndDistance(currCoord, 1, 1)
		if ok && isOccupiedByColor(boardState, coord, enemyColor) {
			valid = append(valid, coord)
		}

		// NE
		coord, ok = GetCoordBySlopeAndDistance(currCoord, 1, -1)
		if ok && isOccupiedByColor(boardState, coord, enemyColor) {
			valid = append(valid, coord)
		}
	} else {
		// SW
		coord, ok := GetCoordBySlopeAndDistance(currCoord, -1, -1)
		if ok && isOccupiedByColor(boardState, coord, enemyColor) {
			valid = append(valid, coord)
		}

		// SE
		coord, ok = GetCoordBySlopeAndDistance(currCoord, -1, 1)
		if ok && isOccupiedByColor(boardState, coord, enemyColor) {
			valid = append(valid, coord)
		}
	}

	return valid
}

func getValidMovesKing(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {
	enemyColor := GetOppositeColor(playerColor)
	valid := []Coord{}

	coords := GetCoordsBySlopeAndDistanceAll(currCoord, 1)
	for _, coord := range coords {
		if !isOccupied(boardState, coord) || isOccupiedByColor(boardState, coord, enemyColor) {
			valid = append(valid, coord)
		}
	}

	return valid
}

func getValidMovesRook(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {
	enemyColor := GetOppositeColor(playerColor)
	valid := []Coord{}

	slopes := []DirectionSlope{
		SlopeNorth,
		SlopeEast,
		SlopeSouth,
		SlopeWest,
	}

	for _, slope := range slopes {
		yChange := slope[0]
		xChange := slope[1]
		for i := 0; i < 8; i++ {
			coords := GetCoordsBySlopeAndDistance(currCoord, yChange, xChange, i)
			for _, coord := range coords {
				if !isOccupied(boardState, coord) {
					valid = append(valid, coord)
				} else if isOccupiedByColor(boardState, coord, enemyColor) {
					valid = append(valid, coord)
					break
				} else {
					break
				}
			}
		}
	}

	return valid
}

func getValidMovesKnight(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {

	valid := []Coord{}

	var slopes = [][]int{
		[]int{2, -1},
		[]int{2, 1},
		[]int{1, 2},
		[]int{-1, 2},
		[]int{-2, 1},
		[]int{-2, -1},
		[]int{1, -2},
		[]int{-1, -2},
	}

	for _, slope := range slopes {
		coords := GetCoordsBySlopeAndDistance(currCoord, slope[0], slope[1], 1)
		for _, coord := range coords {
			if !isOccupied(boardState, coord) || isOccupiedByColor(boardState, coord, GetOppositeColor(playerColor)) {
				valid = append(valid, coord)
			}
		}

	}

	return valid
}

func getValidMovesForBishop(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {
	enemyColor := GetOppositeColor(playerColor)
	valid := []Coord{}

	slopes := []DirectionSlope{
		SlopeNorthEast,
		SlopeSouthEast,
		SlopeSouthWest,
		SlopeNorthWest,
	}

	for _, slope := range slopes {
		yChange := slope[0]
		xChange := slope[1]
		for i := 0; i < 8; i++ {
			coords := GetCoordsBySlopeAndDistance(currCoord, yChange, xChange, i)
			for _, coord := range coords {
				if !isOccupied(boardState, coord) {
					valid = append(valid, coord)
				} else if isOccupiedByColor(boardState, coord, enemyColor) {
					valid = append(valid, coord)
					break
				} else {
					break
				}
			}
		}
	}

	return valid
}

func getValidMovesForQueen(playerColor PlayerColor, boardState BoardState, currCoord Coord) []Coord {
	diagonals := getValidMovesForBishop(playerColor, boardState, currCoord)
	horizontals := getValidMovesRook(playerColor, boardState, currCoord)
	valid := []Coord{}
	valid = append(valid, diagonals...)
	valid = append(valid, horizontals...)
	return valid
}

// IsDestinationValid checks if the specified color can move to the
func IsDestinationValid(whitesMove bool, isOccupied bool, occupant OnBoardData) bool {
	if isOccupied {
		if whitesMove && occupant.Color == PlayerBlack {
			return true
		} else if !whitesMove && occupant.Color == PlayerWhite {
			return true
		}
	} else {
		return true
	}
	return false
}

func isKingInCheck(boardState BoardState, playerColor PlayerColor) (bool, []ThreateningPiece) {

	// get the specified  player's king coord
	kingCoord, ok := findKingCoordByColor(boardState, playerColor)
	if !ok {
		fmt.Println("could not find the king")
	}

	enemyColor := GetOppositeColor(playerColor)

	inCheck := false
	threateningPieces := []ThreateningPiece{}

	// loop through all pieces
	for coord, pieceData := range boardState {
		// check if enemy has king in check
		if pieceData.Piece != King && pieceData.Color == enemyColor {
			// fmt.Printf("checking if %s is putting the king in check\n", pieceData.Piece)

			// get valid moves for enemy piece
			moves := GetValidMoves(enemyColor, pieceData.Piece, boardState, coord)

			// check if any of the moves currently put the king in check
			for _, move := range moves {
				// if any of these moves is where the king is, then it is in check
				if move.Rank == kingCoord.Rank && move.File == kingCoord.File {
					fmt.Printf("%s has the king in check\n", pieceData.Piece)
					inCheck = true
					threateningPieces = append(threateningPieces, ThreateningPiece{
						Color: pieceData.Color,
						Piece: pieceData.Piece,
						Coord: coord,
					})
					break
				}
			}
		}
	}

	return inCheck, threateningPieces
}

func findKingCoordByColor(boardState BoardState, playerColor PlayerColor) (Coord, bool) {
	for coord, pieceData := range boardState {
		if pieceData.Piece == King && pieceData.Color == playerColor {
			return coord, true
		}
	}
	return Coord{}, false
}

type ThreateningPiece struct {
	Color PlayerColor
	Piece Piece
	Coord Coord
}

type InCheckData struct {
	InCheck                                      bool
	BlackThreateningWhite, WhiteThreateningBlack []ThreateningPiece
}

func GetInCheckData(boardState BoardState) InCheckData {
	data := InCheckData{
		BlackThreateningWhite: buildThreateningPieceSlice(boardState, PlayerWhite),
		WhiteThreateningBlack: buildThreateningPieceSlice(boardState, PlayerBlack),
	}
	data.InCheck = len(data.BlackThreateningWhite) > 0 || len(data.WhiteThreateningBlack) > 0
	return data
}

func buildThreateningPieceSlice(boardState BoardState, color PlayerColor) []ThreateningPiece {
	data := []ThreateningPiece{}
	if ok, threateningPieces := isKingInCheck(boardState, color); ok {
		fmt.Printf("%s is in check by %d %s pieces\n", color, len(threateningPieces), GetOppositeColor(color))
		for _, threateningPiece := range threateningPieces {
			fmt.Printf("%s %s is threatening from file %s rank %s\n",
				threateningPiece.Color,
				threateningPiece.Piece,
				FileToFileView[threateningPiece.Coord.File],
				RankToRankView[threateningPiece.Coord.Rank],
			)
			data = append(data, threateningPiece)
		}
	}
	return data
}
