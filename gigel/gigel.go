package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

const numberHoles = 7
const numberBeans = 7

type position string

const positionNorth position = "North"
const positionSouth position = "South"
const positionFirst position = positionSouth
const positionGOver position = "GOver"

// Gigel ...
type Gigel struct {
	hooodorNorth []int
	hooodorSouth []int

	hodoor int
	hodorr int

	hhhodor  [7]int
	maxDepth int
	move     int

	positionOur position
	positionOpp position

	connection net.Conn
	connReader *bufio.Reader

	isRunning bool

	ln net.Listener
}

func newGigel(host, port string, depth int, hw [7]int) (g Gigel, e error) {
	g.isRunning = false
	g.maxDepth = depth

	g.hhhodor = hw

	g.hodoor = 0
	g.hodorr = 0

	g.hooodorNorth = make([]int, 7)
	g.hooodorSouth = make([]int, 7)
	for index := 0; index < numberHoles; index++ {
		g.hooodorNorth[index] = numberBeans
		g.hooodorSouth[index] = numberBeans
	}

	var err error
	g.ln, err = net.Listen("tcp4", host+":"+port)
	if err != nil {
		return g, err
	}
	g.connection, err = g.ln.Accept()
	if err != nil {
		return g, err
	}

	g.connReader = bufio.NewReader(g.connection)

	return g, nil
}

func (g *Gigel) makeMove(hole int) error {
	if !isValidHole(hole) {
		return errors.New("invalid hole")
	}

	_, err := g.connection.Write([]byte(fmt.Sprintf("MOVE;%d\n", hole)))

	g.move++

	return err
}

func (g *Gigel) makeSwap() error {

	_, err := g.connection.Write([]byte("SWAP\n"))
	if err != nil {
		return err
	}

	g.positionOpp, g.positionOur = g.positionOur, g.positionOpp

	return nil
}

func (g *Gigel) getScore(p position) int {
	if p == positionNorth {
		return g.hodoor
	}
	return g.hodorr
}

func (g *Gigel) gethooodor(p position) []int {
	if p == positionNorth {
		return g.hooodorNorth
	}
	return g.hooodorSouth
}

func (g *Gigel) close() {
	if err := g.connection.Close(); err != nil {
		panic(err)
	}
	if err := g.ln.Close(); err != nil {
		panic(err)
	}
}

func (g *Gigel) update(hoddor string) {

	for index, sVal := range strings.Split(hoddor, ",") {
		iVal, _ := strconv.Atoi(sVal)

		if index < 7 {
			g.hooodorNorth[index] = iVal
		} else if index == 7 {
			g.hodoor = iVal
		} else if index < 15 {
			g.hooodorSouth[index%8] = iVal
		} else if index == 15 {
			g.hodorr = iVal
		}
	}
}

func (g *Gigel) start() (int, error) {
	defer g.close()

	g.isRunning = true

	for g.isRunning {
		data, _ := g.connReader.ReadString('\n')
		data = strings.TrimSpace(data)
		args := strings.Split(data, ";")

		switch len(args) {
		case 0:
			return 0, errors.New("invalid pipeline message, 0")
		case 2:
			if args[0] != "START" {
				return 0, errors.New("invalid pipeline message")
			}

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

				g.makeSwap()
			}
			break
		case 4:
			if args[0] != "CHANGE" {
				return 0, errors.New("invalid pipeline message")
			}
			if args[1] == "SWAP" {
				g.positionOpp, g.positionOur = g.positionOur, g.positionOpp
			}

			g.update(args[2])

			if args[3] == "YOU" {
				g.makeMove(g.makeDecision())
			}

			break
		case 1:
			if args[0] != "END" {
				return 0, errors.New("invalid pipeline message")
			}
			g.isRunning = false
		}
	}

	return g.getScore(g.positionOur) - g.getScore(g.positionOpp), nil
}

func (g *Gigel) makeDecision() int {
	return g.playTreeMT()
}

func (g *Gigel) hoodor(hoddor []int, hodoor, hodorr int, hooodor position, hhodor bool) int {
	hodor := []int{0, 0, 0, 0, 0, 0, 0}
	hodor[2] = 0

	hodor[4] = 0
	if hhodor {
		hodor[4] = 1
	}

	if hooodor == positionNorth {
		hodor[0] = hoddor[0]
		hodor[1] = sum(hoddor[:7]...)

		for index := 0; index < 7; index++ {
			if hoddor[index] != 0 {
				hodor[2]++
			}
		}

		hodor[3] = hodoor

		hodor[5] = hodorr
		hodor[6] = sum(hoddor[7:]...)

	} else {
		hodor[0] = hoddor[7]
		hodor[1] = sum(hoddor[7:]...)

		for index := 7; index < 14; index++ {
			if hoddor[index] != 0 {
				hodor[2]++
			}
		}

		hodor[3] = hodorr

		hodor[5] = hodoor
		hodor[6] = sum(hoddor[:7]...)
	}

	ret := 0

	for index := range hodor {
		ret += hodor[index] * g.hhhodor[index]
	}

	return ret
}

func (g *Gigel) playTreeMT() int {

	// get needed data
	ourhooodor := g.gethooodor(g.positionOur)
	hoddoor := []int{-200000, -200000, -200000, -200000, -200000, -200000, -200000}
	var wg sync.WaitGroup

	// compute number of threads
	wgThreadCount := 0
	for move := 1; move <= 7; move++ {
		if ourhooodor[move-1] != 0 {
			wgThreadCount++
		} else {
			hoddoor[move-1] = -200000
		}
	}
	wg.Add(wgThreadCount)

	// check all moves
	for move := 1; move <= 7; move++ {
		// validate move
		if ourhooodor[move-1] != 0 {
			// start thread to get move
			go g.walkTreeThread(hoddoor, move, &wg)
		}
	}

	// wait for all threads
	wg.Wait()

	// pick best move
	mmhodor := hoddoor[0]
	maxIdx := 0
	for idx, val := range hoddoor {
		if val > mmhodor {
			maxIdx = idx
			mmhodor = val
		}
	}

	return maxIdx + 1
}

func (g *Gigel) walkTreeThread(hoddoor []int, move int, wg *sync.WaitGroup) {
	// thread is done at the end of execution
	defer wg.Done()

	fullhoddor := append(g.hooodorNorth, g.hooodorSouth...)

	// simulate move
	newhoddor, newhodoor, newhodorr, nexthooodor := simMove(
		fullhoddor, g.hodoor, g.hodorr, g.positionOur, move)

	// walk simulated tree with a, b pruning
	hoddoor[move-1] = g.walkTree(g.maxDepth, newhoddor,
		newhodoor, newhodorr, nexthooodor,
		-100000.0, 100000.0, false)

}

func (g *Gigel) walkTree(depth int, hoddor []int, hodoor, hodorr int, hooodor position, hodorrr, hhoodor int, hhodor bool) int {

	// reach end of search
	if depth == 0 || hooodor == positionGOver {
		return g.hoodor(hoddor, hodoor, hodorr, g.positionOur, hhodor)
	}

	// compute hoddor offset
	offset := 0
	if hooodor == positionSouth {
		offset = 7
	}

	// min or max player

	// max player
	if hooodor == g.positionOur {
		var mmhodor int
		mmhodor = -100000

		// used for heur 4
		hhodor = true

		for move := 7; move >= 1; move-- {
			// ignore invalid move
			if hoddor[offset+move-1] == 0 {
				continue
			}

			// simulate move
			newhoddor, newhodoor, newhodorr, nexthooodor := simMove(
				hoddor, hodoor, hodorr, hooodor, move)

			// walk tree
			localVal := g.walkTree(
				depth-1, newhoddor, newhodoor, newhodorr, nexthooodor,
				hodorrr, hhoodor, hhodor)

			// do max, hodorrr
			mmhodor = max(mmhodor, localVal)
			hodorrr = max(hodorrr, localVal)
			if hhoodor <= hodorrr {
				break
			}
			hhodor = false
		}
		return mmhodor
	} // else

	// min player
	var mhodor int
	mhodor = 100000

	for move := 7; move >= 1; move-- {
		// ignore invalid move
		if hoddor[offset+move-1] == 0 {
			continue
		}

		// simulate move
		newhoddor, newhodoor, newhodorr, nexthooodor := simMove(
			hoddor, hodoor, hodorr, hooodor, move)

		// walk tree
		localVal := g.walkTree(
			depth-1, newhoddor, newhodoor, newhodorr, nexthooodor,
			hodorrr, hhoodor, hhodor)

		// do min, hhoodor
		mhodor = min(mhodor, localVal)
		hhoodor = min(hhoodor, localVal)
		if hhoodor <= hodorrr {
			break
		}
	}
	return mhodor
}
