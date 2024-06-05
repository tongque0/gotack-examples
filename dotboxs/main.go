package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Dab struct {
	Human         int
	Robot         int
	Record        [][]int
	TimeoutOffset float64
	Turn          int
	Now           int
}

type Robot struct {
	dab *Dab
}

func NewRobot(dab *Dab) *Robot {
	return &Robot{dab: dab}
}

// GetBestMove now returns int32 values for H and V, H表示横线，第几行第几列的横线，V表示竖线，第几行第几列的竖线
func (r *Robot) GetBestMove() (int32, int32, error) {
	var s0, s1 int
	s0, s1 = r.dab.Human, r.dab.Robot

	h, v := int32(0), int32(0)
	for _, move := range r.dab.Record {
		x, y := move[1], move[2]
		if move[0] == 0 {
			v |= (1 << (y*6 + x))
		} else {
			h |= (1 << (x*6 + y))
		}
	}
	algorithm := "quctann"
	timeout := uint(10+60*r.dab.TimeoutOffset) * 1000

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		return 0, 0, fmt.Errorf("error connecting: %v", err)
	}
	defer conn.Close()

	arg := map[string]interface{}{
		"method": "Server.MakeMove",
		"params": []interface{}{map[string]interface{}{
			"Algorithm": algorithm,
			"Board": map[string]interface{}{
				"H":    h,
				"V":    v,
				"S":    []int{s0, s1},
				"Now":  r.dab.Now,
				"Turn": r.dab.Turn,
			},
			"Timeout": timeout,
		}},
		"id": int(time.Now().Unix()),
	}

	data, err := json.Marshal(arg)
	if err != nil {
		return 0, 0, fmt.Errorf("error marshaling: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return 0, 0, fmt.Errorf("error writing: %v", err)
	}

	buff := make([]byte, 4096)
	n, err := conn.Read(buff)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(buff[:n], &res); err != nil {
		return 0, 0, fmt.Errorf("error unmarshaling: %v", err)
	}

	if result, ok := res["result"].(map[string]interface{}); ok {
		hResult, hOk := result["H"].(float64)
		vResult, vOk := result["V"].(float64)
		if hOk && vOk {
			return int32(hResult), int32(vResult), nil
		} else {
			return 0, 0, fmt.Errorf("unexpected result format: %v", result)
		}
	} else {
		return 0, 0, fmt.Errorf("unexpected response format: %v", res)
	}
}
func DecodeMove(hMask int32, vMask int32) (int, int, int, error) {
	// Determine which mask to process
	var mask int32
	var lineType int // 0 for horizontal (H), 1 for vertical (V)

	if hMask != 0 && vMask != 0 {
		return 0, 0, 0, fmt.Errorf("both masks cannot be non-zero")
	} else if hMask != 0 {
		mask = hMask
		lineType = 0
	} else if vMask != 0 {
		mask = vMask
		lineType = 1
	} else {
		return 0, 0, 0, fmt.Errorf("both masks cannot be zero")
	}

	// Decode the mask
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			pos := i*6 + j
			if (mask & (1 << pos)) != 0 {
				return lineType, i + 1, j + 1, nil
			}
		}
	}

	return 0, 0, 0, fmt.Errorf("no move found in the provided mask")
}

func EncodeMove(lineType int, row int, col int) (int32, int32, error) {
	if row < 1 || row > 6 || col < 1 || col > 6 {
		return 0, 0, fmt.Errorf("row and col must be within 1 to 6")
	}

	// Calculate the bit position
	pos := (row-1)*6 + (col - 1)

	if lineType == 0 {
		// Horizontal line type
		return (1 << pos), 0, nil
	} else if lineType == 1 {
		// Vertical line type
		return 0, (1 << pos), nil
	} else {
		return 0, 0, fmt.Errorf("invalid line type: %d", lineType)
	}
}
func parseMove(move string) (k, i, j int, err error) {
	if len(move) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid move format")
	}
	k = int(move[0] - 'A') // 'H' or 'V' converted to 0 or 1
	i = int(move[1] - 'A')
	j = int(move[2] - 'A')
	return k, i, j, nil
}

func main() {
	gameState := &Dab{
		Human:         0,
		Robot:         1,
		Record:        [][]int{},
		TimeoutOffset: 0,
		Turn:          1,
		Now:           1, // Robot starts first
	}
	robotAI := NewRobot(gameState)
	robotSide := 1 // Assume robot is white

	input := os.Stdin
	output := os.Stdout

	var command string
	for {
		fmt.Fscanf(input, "%s", &command)
		fmt.Print(robotAI.dab.Record)
		switch command {
		case "new":
			var playerColor string
			fmt.Fscanf(input, "%s", &playerColor)
			if playerColor == "black" {
				robotSide = 0
			} else {
				robotSide = 1
			}
			gameState.Human = robotSide ^ 1
			gameState.Robot = robotSide
			gameState.Record = [][]int{}
			gameState.Turn = 1
			gameState.Now = robotSide
			if robotSide == 0 {
				horiz, vert, err := robotAI.GetBestMove()
				if err != nil {
					fmt.Fprintln(output, "Error getting next move:", err)
					continue
				}
				lineType, row, col, _ := DecodeMove(horiz, vert)
				maskH, maskV, _ := EncodeMove(lineType, row, col)
				gameState.Record = append(gameState.Record, []int{lineType, int(maskH), int(maskV)})
				moveCommand := fmt.Sprintf("move 1 %c%c%c\n", lineType+'A', row+'A', col+'A')
				fmt.Fprint(output, moveCommand)
			}
		case "move":
			var numMoves int
			fmt.Fscanf(input, "%d", &numMoves)
			for idx := 0; idx < numMoves; idx++ {
				var moveStr string
				fmt.Fscanf(input, "%s", &moveStr)
				lineType, row, col, _ := parseMove(moveStr)
				fmt.Println(moveStr, lineType, row, col)
				maskH, maskV, _ := EncodeMove(lineType, row, col)
				gameState.Record = append(gameState.Record, []int{lineType, int(maskH), int(maskV)})
				gameState.Turn++
				responseHoriz, responseVert, _ := robotAI.GetBestMove()
				responseType, responseRow, responseCol, _ := DecodeMove(responseHoriz, responseVert)
				responseMaskH, responseMaskV, _ := EncodeMove(responseType, responseRow, responseCol)
				gameState.Record = append(gameState.Record, []int{responseType, int(responseMaskH), int(responseMaskV)})
				responseCommand := fmt.Sprintf("move 1 %c%c%c\n", responseType+'A', responseRow+'A', responseCol+'A')
				fmt.Fprint(output, responseCommand)
			}
		case "name?":
			fmt.Fprintln(output, "name DotV2.0.88")
		case "end", "quit":
			fmt.Fprintln(output, "Quit!")
			return
		}
	}
}
