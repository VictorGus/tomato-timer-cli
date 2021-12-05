package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var ARGUMENTS = map[string]string{
	"--minutes": "-m",
	"--seconds": "-s",
	"--hours":   "-h",
	"-m":        "-m",
	"-s":        "-s",
	"-h":        "-h"}

func PrepareKey(key string) string {
	key = ARGUMENTS[key]
	return strings.ReplaceAll(key, "-", "")
}

func ConvertSliceToMap(array []string) map[string]string {
	var resultMap = make(map[string]string)

	for i := 0; i < len(array); i += 2 {
		var preparedKey = PrepareKey(array[i])
		resultMap[preparedKey] = array[i+1]
	}

	return resultMap
}

func ShowHelp() {
	fmt.Println()
	fmt.Println("----------------")
	fmt.Println("| Tomato Timer |")
	fmt.Println("----------------")
	fmt.Println()
	fmt.Println("POSSIBLE OPTIONS:")
	fmt.Println("   --hours int_value, -h int_value    Amount of hours")
	fmt.Println("   --minutes int_value, -m int_value  Amount of minutes")
	fmt.Println("   --seconds int_value, -s int_value  Amount of seconds")
}

func IsInArray(argKey string, args []string) bool {
	for _, v := range args {
		if argKey == v {
			return true
		}
	}
	return false
}

func ValidateInputArguments(arguments []string) error {
	if len(arguments)%2 != 0 {
		return errors.New("wrong number of input args")
	}

	if IsInArray("-m", arguments) || IsInArray("--minutes", arguments) ||
		IsInArray("-h", arguments) || IsInArray("--hours", arguments) ||
		IsInArray("-s", arguments) || IsInArray("--seconds", arguments) {
		return nil
	} else {
		return errors.New("no required args supplied")
	}
}

func ConvertTimePart(timePart int) string {
	if timePart < 10 {
		return fmt.Sprintf("0%d", timePart)
	} else {
		return fmt.Sprintf("%d", timePart)
	}
}

func ExtractTimePartAsInt(timePart string, argumentsMap map[string]string) (int, error) {
	v, isInPlace := argumentsMap[timePart]
	if isInPlace {
		convertedValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		return convertedValue, nil
	} else {
		return 0, nil
	}
}

func Beep() {
	f, err := os.Open("./resources/mixkit-alarm-tone-996.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N((time.Second / 10)))
	if err != nil {
		log.Fatal(err)
	}

	speaker.Play(streamer)
	time.Sleep(3 * time.Second)
}

func main() {
	arguments := os.Args[1:]
	validationErr := ValidateInputArguments(arguments)
	if validationErr != nil {
		ShowHelp()
		return
	}

	argumentsMap := ConvertSliceToMap(arguments)
	ticker := time.NewTicker(1 * time.Second)
	doneChannel := make(chan bool)

	hours, err := ExtractTimePartAsInt("h", argumentsMap)
	if err != nil {
		log.Fatal(err)
	}

	minutes, err := ExtractTimePartAsInt("m", argumentsMap)
	if err != nil {
		log.Fatal(err)
	}

	seconds, err := ExtractTimePartAsInt("s", argumentsMap)
	if err != nil {
		log.Fatal(err)
	}

	timeSumSeconds := 3600*hours + 60*minutes + seconds

	fmt.Println("#--")

	go func() {
		for i := timeSumSeconds; ; i-- {
			h := i / 3600
			m := i / 60 % 60
			s := i % 60
			select {
			case <-doneChannel:
				Beep()
				return
			case <-ticker.C:
				fmt.Printf("\r#-- Time left: %s:%s:%s", ConvertTimePart(h), ConvertTimePart(m), ConvertTimePart(s))
			}
		}
	}()

	time.Sleep((time.Duration(timeSumSeconds) + 1) * time.Second)
	doneChannel <- true
	ticker.Stop()
	fmt.Println("\n#-- Done!")
	fmt.Println("#--")
	fmt.Println()
	time.Sleep(1 * time.Second)
}
