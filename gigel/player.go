package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type player struct {
	sideNorth []int
	sideSouth []int

	scoreNorth int
	scoreSouth int

	positionOur position
	positionOpp position

	connection net.Conn
	connReader *bufio.Reader
	moveReader *bufio.Reader

	isRunning bool

	ln net.Listener
}

func newPlayer(host, port string) (g player, e error) {
	g.isRunning = false

	// set weights

	// set score
	g.scoreNorth = 0
	g.scoreSouth = 0

	// set game board
	g.sideNorth = make([]int, 7)
	g.sideSouth = make([]int, 7)
	for index := 0; index < numberHoles; index++ {
		g.sideNorth[index] = numberBeans
		g.sideSouth[index] = numberBeans
	}

	// open connection to pipeline
	var err error
	g.ln, err = net.Listen("tcp4", host+":"+port)
	if err != nil {
		return g, err
	}
	g.connection, err = g.ln.Accept()
	if err != nil {
		return g, err
	}

	// open reader on pipeline
	g.connReader = bufio.NewReader(g.connection)
	g.moveReader = bufio.NewReader(os.Stdin)

	return g, nil
}

func (g *player) makeMove(hole int) error {
	// validate hole
	if !isValidHole(hole) {
		return errors.New("invalid hole")
	}

	// write move to pipeline
	_, err := g.connection.Write([]byte(fmt.Sprintf("MOVE;%d\n", hole)))

	return err
}

func (g *player) makeSwap() error {

	// write swap instruction to pipeline
	_, err := g.connection.Write([]byte("SWAP\n"))
	if err != nil {
		return err
	}

	// swap positions
	g.positionOpp, g.positionOur = g.positionOur, g.positionOpp

	return nil
}

func (g *player) getScore(p position) int {
	if p == positionNorth {
		return g.scoreNorth
	}
	return g.scoreSouth
}

func (g *player) getSide(p position) []int {
	if p == positionNorth {
		return g.sideNorth
	}
	return g.sideSouth
}

func (g *player) close() {
	if err := g.connection.Close(); err != nil {
		panic(err)
	}
	if err := g.ln.Close(); err != nil {
		panic(err)
	}
}

func (g *player) update(board string) {

	// iterate new boar
	for index, sVal := range strings.Split(board, ",") {
		// get integers
		iVal, _ := strconv.Atoi(sVal)

		// update sides and scores
		if index < 7 {
			g.sideNorth[index] = iVal
		} else if index == 7 {
			g.scoreNorth = iVal
		} else if index < 15 {
			g.sideSouth[index%8] = iVal
		} else if index == 15 {
			g.scoreSouth = iVal
		}
	}

	if g.positionOpp == positionNorth {
		fmt.Print("\033[0;31m")
	}

	fmt.Printf(" %02d ", g.scoreNorth)
	for idx := 6; idx >= 0; idx-- {
		fmt.Printf(" %2d ", g.sideNorth[idx])
	}
	fmt.Printf("____\n")

	fmt.Print("\033[0m")
	if g.positionOpp == positionSouth {
		fmt.Print("\033[0;31m")
	}
	fmt.Printf("____")

	for idx := 0; idx < 7; idx++ {
		fmt.Printf(" %2d ", g.sideSouth[idx])
	}
	fmt.Printf(" %02d ", g.scoreSouth)
	fmt.Printf("\033[0m\n")
	fmt.Println("------------------------------------")
}

func (g *player) start() (int, error) {
	defer g.close()

	g.isRunning = true

	for g.isRunning {
		// while game is running, scan pipeline
		data, _ := g.connReader.ReadString('\n')
		data = strings.TrimSpace(data)
		args := strings.Split(data, ";")

		// check cases
		switch len(args) {
		// invalid line
		case 0:
			return 0, errors.New("invalid pipeline message, 0")
		// starting position and first move
		case 2:
			if args[0] != "START" {
				return 0, errors.New("invalid pipeline message")
			}

			// we start south
			if args[1] == "South" {
				g.positionOur = positionSouth
				g.positionOpp = positionNorth
			} else {
				g.positionOur = positionNorth
				g.positionOpp = positionSouth
			}

			if g.positionOur == positionFirst {
				g.makeMove(g.makeDecision())
			} else {

				data, _ := g.connReader.ReadString('\n')
				data = strings.TrimSpace(data)
				args := strings.Split(data, ";")
				g.update(args[2])

				// always swap
				g.makeSwap()
			}
			break
		// change in board state, positions and turn
		case 4:
			if args[0] != "CHANGE" {
				return 0, errors.New("invalid pipeline message")
			}
			if args[1] == "SWAP" {
				// swap positions
				g.positionOpp, g.positionOur = g.positionOur, g.positionOpp
			}

			// update game board
			g.update(args[2])

			// on our turn make a move
			if args[3] == "YOU" {
				g.makeMove(g.makeDecision())
			}

			break
		// end of game
		case 1:
			if args[0] != "END" {
				return 0, errors.New("invalid pipeline message")
			}
			g.isRunning = false
		}
	}

	return g.getScore(g.positionOur) - g.getScore(g.positionOpp), nil
}

func (g player) makeDecision() int {

	fmt.Print("MOVE: ")
	text, err := g.moveReader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	text = strings.TrimSpace(text)

	move, err := strconv.Atoi(text)
	if err != nil {
		log.Fatalln(err)
	}

	return move
}
