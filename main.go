package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Speed int

type flower struct {
	file string
	Pos  int
}

func main() {
	Speed = 40
	size := GetTerminalSize()
	tic := 0
	var pos []flower
	var h flower
	h.Pos = size - (size / 3)
	h.file = "1.txt"
	pos = append(pos, h)
	delInPos := false
	jump := 0
	widthjump := 0
	// fmt.Println(pos, h)
	// return

	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	for {
		if tic%100 == 0 && Speed > 20 {
			Speed = Speed - 2
		}
		stdin := ""
		select {
		case stdin, _ = <-ch:
			fmt.Println("Keys pressed:", stdin)
		default:
			fmt.Println("Working..")
		}
		if delInPos {
			pos = pos[1:]
			delInPos = false
		}
		if jump != 0 {
			switch widthjump {
			case 1:
				jump = 4
				widthjump++
			case 2:
				jump = 5
				widthjump++
			case 3:
				jump = 7
				widthjump++
			case 4:
				jump = 8
				widthjump++
			case 5:
				jump = 8
				widthjump++
			case 6:
				jump = 8
				widthjump++
			case 7:
				jump = 9
				widthjump++
			case 8:
				jump = 9
				widthjump++
			case 9:
				jump = 9
				widthjump++
			case 10:
				jump = 8
				widthjump++
			case 11:
				jump = 8
				widthjump++
			case 12:
				jump = 8
				widthjump++
			case 13:
				jump = 7
				widthjump++
			case 14:
				jump = 5
				widthjump++
			case 15:
				jump = 4
				widthjump++
			case 16:
				jump = 0
				widthjump = 0
			}
		}
		if stdin == " " && jump == 0 {
			jump = 2
			widthjump++
		}

		stdin = ""
		speed := time.Millisecond * time.Duration(Speed)
		time.Sleep(speed)
		fraim := rendering(size, tic)

		fraim = addDino(fraim, tic, jump, widthjump)
		for j, i := range pos {
			fraim = barrier(fraim, i)

			if size-i.Pos > 30+rand.Intn(20) && j+1 == len(pos) {
				name := strconv.Itoa(rand.Intn(6))
				var h flower
				h.Pos = size
				h.file = name + ".txt"
				pos = append(pos, h)
			}

			pos[j].Pos = pos[j].Pos - 1
			if i.Pos == 0 {
				delInPos = true
			}
		}

		fraim = addScore(fraim, tic)

		fmt.Println("\033[2J")
		fmt.Println(rand.Intn(10))
		fmt.Println(pos)
		for _, i := range fraim {
			// fmt.Println(len(i))
			fmt.Println(i)
		}
		tic++
	}
}

func rendering(size, itr int) []string {

	var frame []string
	for i := 0; i < 20; i++ {
		line := ""
		if i > 16 {
			for j := 0; j < size; j++ {
				if j == 0 && itr%2 != 1 {
					j++
				}
				if j%2 != 0 {
					line += "~"
				} else {
					line += " "
				}
			}
		} else {
			for j := 0; j < size; j++ {
				line += " "
			}
		}
		frame = append(frame, line)
	}
	return frame
}

func addDino(frame []string, tic, jump, widthjump int) []string {
	tic = (tic % 4) + 1
	name := fmt.Sprintf("sprite/%v.txt", tic)
	if jump != 0 {
		name = "sprite/5.txt"
		if widthjump%2 == 1 {
			name = "sprite/6.txt"
		}
	}
	dino, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	di := strings.Split(string(dino), "\n")
	for i, k := range di {
		a := frame[20-3-len(di)+i-jump]
		l := len(k)
		a = a[:10] + k + a[l+10:]
		frame[20-3-len(di)+i-jump] = a
	}
	return frame
}

func barrier(frame []string, pos flower) []string {

	cactus, err := os.ReadFile("sprite/barrier/" + pos.file)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	a := strings.Split(string(cactus), "\n")
	ln := len(a)
	fmt.Print(len(a))
	for i, k := range a {
		b := frame[20-3-ln+i]
		l := len(k)
		if l+pos.Pos > len(b) {
			continue
		}
		if b[pos.Pos:l+pos.Pos] != createAbys(l) {
			os.Exit(0)
		}
		b = b[:pos.Pos] + k + b[l+pos.Pos:]
		frame[20-3-ln+i] = b
	}
	return frame
}

func createAbys(l int) string {
	abys := ""
	for i := 0; i < l; i++ {
		abys += " "
	}
	return abys
}

func addScore(frame []string, score int) []string {
	sco := strconv.Itoa(score)
	a := frame[3]
	a = a[:len(a)-5-len(sco)] + sco + a[len(a)-5:]
	frame[3] = a
	return frame
}

func GetTerminalSize() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range out {
		if v == ' ' {
			out = out[i+1 : len(out)-1] // out это массив байтов и в него записали размер терминала , затем пробегаемся по нему
			break
		}
	}
	x, err := strconv.Atoi(string(out)) // Так как мы его превращаем в строку , то используем strconv.Atoi для получения стандартных чисел
	if err != nil {
		fmt.Println("ERROR1")
	}
	return x
}
