package main

func simMove(inBoard []int, scoreNorth, scoreSouth int, side position, move int) (newBoard []int, newScoreNorth, newScoreSouth int, newSide position) {
	// if game over
	if side == positionGOver {
		return inBoard, scoreNorth, scoreSouth, side
	}

	board := make([]int, 14)
	copy(board, inBoard)

	// compute N/S offset on board and default next side
	nextSide := positionSouth
	offset := 0
	if side == positionSouth {
		nextSide = positionNorth
		offset = 7
	}

	// change move to hole index
	startIndex := move - 1 + offset

	// valid hole ?
	if board[startIndex] == 0 {
		panic(nil)
	}

	// remove beans from hole
	beans := board[startIndex]
	board[startIndex] = 0

	// move beans to next holes and score
	for index := 1; index <= beans; index++ {
		// compute real index in range 0-13
		realIndex := (startIndex + index) % 14

		// check for scores
		if side == positionNorth {
			// check for north score
			if realIndex == 7 && beans-index >= 0 {
				scoreNorth++

				// last move, in house
				if beans-index == 0 {
					nextSide = positionNorth
					break
				}
				beans--
			}
		} else {
			// check for south score
			if realIndex == 0 && beans-index >= 0 {
				scoreSouth++

				// last move, in house
				if beans-index == 0 {
					if scoreSouth != 1 {
						nextSide = positionSouth
					}
					break
				}
				beans--
			}
		}

		// increment next hole
		board[realIndex]++

		// last move, on board, would possibly steal from other side
		oppositeHoleIndex := 13 - realIndex
		if beans-index == 0 && board[realIndex] == 1 && board[oppositeHoleIndex] != 0 {
			if realIndex < 7 && side == positionNorth {
				scoreNorth += board[realIndex]
				scoreNorth += board[oppositeHoleIndex]

				board[realIndex] = 0
				board[oppositeHoleIndex] = 0
			} else if realIndex > 6 && side == positionSouth {
				scoreSouth += board[realIndex]
				scoreSouth += board[oppositeHoleIndex]

				board[realIndex] = 0
				board[oppositeHoleIndex] = 0
			}
		}
	}

	// check if any of the players finished
	northFinished := true
	southFinished := true
	for hole := 0; hole < len(board); hole++ {
		if board[hole] != 0 {
			if hole < 7 {
				northFinished = false
			} else {
				southFinished = false
			}
		}
	}

	// give all beans to south
	if northFinished {
		for hole := 7; hole < len(board); hole++ {
			scoreSouth += board[hole]
			board[hole] = 0
		}
		nextSide = positionGOver
	}

	// give all beans to north
	if southFinished {
		for hole := 0; hole < 7; hole++ {
			scoreNorth += board[hole]
			board[hole] = 0
		}
		nextSide = positionGOver
	}

	return board, scoreNorth, scoreSouth, nextSide
}
