package main

import (
	_ "fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
	"fmt"
	"math"
	"strings"
	"strconv"
)

type Status uint

const (
	WIN  Status = iota
	LOSE
	ADD
	MAX  = 64
)

var helpInfo = []string{
	"Enter: Restart the game.",
	"Esc: \t Quit the game.",
}
var colorMap = map[int]termbox.Attribute{
	2:    termbox.ColorWhite,
	4:    termbox.ColorWhite,
	8:    termbox.ColorCyan,
	16:   termbox.ColorCyan,
	32:   termbox.ColorCyan,
	64:   termbox.ColorGreen,
	128:  termbox.ColorMagenta,
	256:  termbox.ColorYellow,
	512:  termbox.ColorBlue,
	1024: termbox.ColorRed,
}
var mrgLen int = 4 // 边距
var arrLen int     // 矩阵长度

var Step int
var Score int

// 输出字符串在中间
func PrintStr(str string) error {
	str += time.Now().Format(" 15:04:05")
	fg := termbox.ColorYellow
	bg := termbox.ColorBlack
	unitX := mrgLen * 5 / 2 * arrLen / 2
	unitY := mrgLen * 3 / 2 * arrLen / 2
	for offsetX, char := range str {
		termbox.SetCell(unitX+offsetX, unitY, char, fg, bg)
	}
	termbox.Flush()
	return nil
}

// 在底部输出字符串
func PrintStrEnd(str string) error {
	str += time.Now().Format(" 15:04:05")
	fg := termbox.ColorYellow
	bg := termbox.ColorBlack
	y := len(helpInfo) + arrLen*mrgLen + mrgLen*2
	for offsetX, char := range str {
		termbox.SetCell(mrgLen+offsetX, y+1, char, fg, bg)
	}
	termbox.Flush()
	return nil
}

// 2048游戏中的16个格子使用4x4二维数组表示
type G2048 [][]int

func (p G2048) Strings() string {
	gameStr := "G2048: \n"
	for _, line := range p {
		for _, unit := range line {
			gameStr += fmt.Sprintf("[%d]", unit)
		}
		gameStr += "\n"
	}
	return gameStr
}

func (p *G2048) Init(cellNumber int) *G2048 {
	arrLen = cellNumber
	var v [][]int
	for i := 0; i < cellNumber; i++ {
		tmp := make([]int, cellNumber)
		v = append(v, tmp)
	}
	*p = v
	return p
}

// 检查游戏是否已经胜利，没有胜利的情况下随机将值为0的元素
// 随机设置为2或者4
func (p *G2048) checkWinOrAdd() Status {
	count := 0
	var emptyUnit []string
	for x, line := range *p {
		for y, num := range line {
			count ++
			if num >= MAX {
				return WIN
			}
			if num == 0 {
				emptyUnit = append(emptyUnit, fmt.Sprintf("%d:%d", x, y))
			}
		}
	}
	if (len(emptyUnit) > 0) {
		randCount := int(math.Floor(float64(arrLen * arrLen / 16)))
		for i := 0; i < randCount; i++ {
			randNum := 2 << (rand.Uint32() % 2)
			ind := rand.Intn(len(emptyUnit))
			xy := strings.Split(emptyUnit[ind], ":")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			(*p)[x][y] = randNum
			Score += randNum
			//PrintStrEnd(fmt.Sprintf("x: %d,y: %d", x, y))
		}
		return ADD
	} else {
		PrintStrEnd("Your lose the game.")
	}
	return LOSE
}

// 初始化2048 / 刷新界面
func (p G2048) initialize() {
	x, y := 0, 0
	fg := termbox.ColorYellow
	bg := termbox.ColorBlack
	termbox.Clear(fg, bg)

	gameInfo := fmt.Sprintf("Score : %d \t Steps: %d", Score, Step)
	for offsetX, char := range gameInfo {
		termbox.SetCell(x+mrgLen+offsetX, y+1, char, fg, bg)
	}

	// 输出提示信息
	for i, info := range helpInfo {
		offsetY := i + arrLen*mrgLen + mrgLen*2
		str := fmt.Sprint(info)
		for j, char := range str {
			offsetX := mrgLen + j
			termbox.SetCell(x+offsetX, y+offsetY, char, fg, bg)
		}
	}

	// 输出 2048 矩阵
	for offsetY, line := range p {
		unitY := y + mrgLen*3/2
		for offsetX, unit := range line {
			unitX := x + offsetX*(mrgLen*2) + mrgLen/2
			str := fmt.Sprint(unit)
			if unit != 0 {
				color := colorMap[unit]
				for offsetChar, char := range str {
					termbox.SetCell(unitX+offsetChar, unitY+offsetY*mrgLen, char, color, bg)
				}
			}
		}
	}
	// 输出 矩阵方格
	for offsetY := 1; offsetY <= arrLen+1; offsetY++ {
		for offsetX := 1; offsetX < x+arrLen*mrgLen*2; offsetX++ {
			termbox.SetCell(x+offsetX, y+offsetY*mrgLen, '─', fg, bg)
		}
	}
	for offsetY := 1 * mrgLen; offsetY <= (arrLen+1)*mrgLen; offsetY++ {
		for offsetX := x; offsetX <= arrLen*mrgLen; offsetX += mrgLen {
			termbox.SetCell(offsetX*mrgLen/2, y+offsetY, '|', fg, bg)
		}
	}
	termbox.Flush()
}

// 左旋转 90度
func (p *G2048) Left90() {
	tmp := new(G2048).Init(arrLen)
	for x, line := range *p {
		for y, unit := range line {
			(*tmp)[arrLen-1-y][x] = unit
		}
	}
	*p = *tmp
}

// 右旋转 90度
func (p *G2048) Right90() {
	tmp := new(G2048).Init(arrLen)
	for x, line := range *p {
		for y, unit := range line {
			(*tmp)[y][arrLen-1-x] = unit
		}
	}
	*p = *tmp
}

// 旋转 180度
func (p *G2048) Right180() {
	tmp := new(G2048).Init(arrLen)
	for x, line := range *p {
		for y, unit := range line {
			(*tmp)[arrLen-1-x][arrLen-1-y] = unit
		}
	}
	*p = *tmp
}

// 向上移动合并转换为向上移动合并
func (p *G2048) mergeUp() (bool, bool) {
	p.Left90()
	change, isFull := p.mergeLeft()
	p.Right90()
	return change, isFull
}

// 向下移动合并转换为向下移动合并
func (p *G2048) mergeDown() (bool, bool) {
	p.Right90()
	change, isFull := p.mergeLeft()
	p.Left90()
	return change, isFull
}

/// 向左移动合并转换为向左移动合并
func (p *G2048) mergeLeft() (bool, bool) {
	isChange := false
	isFull := true
	// 将元素向左靠拢
	for x, line := range *p {
		empty := -1
		for y, unit := range line {
			if unit == 0 && empty == -1 {
				empty = y
			}
			if unit != 0 && empty > -1 {
				isChange = true
				(*p)[x][empty] = unit
				(*p)[x][y] = 0
				empty = y
			}
		}
	}
	//  将相同元素进行合并
	for x, line := range *p {
		for y := 0; y < arrLen; y++ {
			if line[y] == 0 {
				if (y < arrLen-1) {
					isChange = true
					(*p)[x][y], (*p)[x][y+1] = (*p)[x][y+1], (*p)[x][y]
				}
				isFull = false
			} else if y < arrLen-1 && line[y] == line[y+1] {
				isChange = true
				(*p)[x][y] *= 2
				(*p)[x][y+1] = 0
				isFull = false
			}
		}
	}
	return isChange || isFull, isFull
}

/// 向右移动合并转换为向上移动合并
func (p *G2048) mergeRight() (bool, bool) {
	p.Right180()
	change, isFull := p.mergeLeft()
	p.Right180()
	return change, isFull
}

// 检查按键，做出不同的移动动作或者退出程序
func (p *G2048) mrgeAndReturnKey() termbox.Key {
	var changed, isFull bool
Loop:
	changed, isFull = false, false
	ch := make(chan termbox.Event)
	go func() {
		ev := termbox.PollEvent()
		ch <- ev
	}()
	for {
		switch ev := <-ch; ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				changed, isFull = p.mergeUp()
				PrintStrEnd("Your press upKey:" + fmt.Sprintf("%t", changed))
			case termbox.KeyArrowDown:
				changed, isFull = p.mergeDown()
				PrintStrEnd("Your press downKey:" + fmt.Sprintf("%t", changed))
			case termbox.KeyArrowLeft:
				changed, isFull = p.mergeLeft()
				PrintStrEnd("Your press leftKey:" + fmt.Sprintf("%t", changed))
			case termbox.KeyArrowRight:
				changed, isFull = p.mergeRight()
				PrintStrEnd("Your press rightKey:" + fmt.Sprintf("%t", changed))
			case termbox.KeyEsc, termbox.KeyEnter:
				changed = true
			default:
				changed = false
			}
			if changed {
				if (!isFull) {
					Step++
				}
				return ev.Key
			} else {
				close(ch)
				goto Loop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// 重置
func (p *G2048) clear() {
	Step = 0
	Score = 0
	newGame := new(G2048)
	newGame.Init(arrLen)
	*p = *newGame
}

// 开始游戏
func (p *G2048) Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	rand.Seed(time.Now().UnixNano())
Loop:
	p.clear()
	for {
		p.initialize()
		status := p.checkWinOrAdd()
		switch status {
		case LOSE:
			PrintStrEnd("Game Lose.")
			break
		case WIN:
			PrintStrEnd("Game Win")
			break
		}
		enKey := p.mrgeAndReturnKey()
		if enKey == termbox.KeyEsc {
			return
		} else if enKey == termbox.KeyEnter {
			goto Loop
		}
	}
}

func main() {
	game := new(G2048)
	game.Init(4)
	game.Run()
}
