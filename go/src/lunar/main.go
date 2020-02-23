package main

import (
	"bufio"
	"fmt"
	"lander"
	"math"
	"os"
	"strconv"
	"strings"
)

func landing(lem *lander.Lander, burnRate float64, burnTime float64) {
	lem.CalcImpact(burnRate, burnTime)
	endGame(lem)
}

func outOfFuel(lem *lander.Lander) {
	fmt.Printf("\nFUEL OUT AT %0.2f SECONDS\n", lem.ElapsedTime)
	secondsToImpact := ((-1.0 * lem.Velocity) + math.Sqrt(
		math.Pow(lem.Velocity, 2.0)+(2.0*lem.Altitude*
			lander.Gravity))) / lander.Gravity
	lem.Velocity = lem.Velocity + (lander.Gravity * secondsToImpact)
	endGame(lem)
}

func endGame(lem *lander.Lander) {
	velocityMph := lem.Velocity * 3600.0
	fmt.Printf("\nON THE MOON AT %0.2f SECONDS\nIMPACT VELOCITY %0.2f MPH\n",
		lem.ElapsedTime, velocityMph)
	if velocityMph < 1.2 {
		fmt.Println("PERFECT LANDING!!")
	} else if velocityMph <= 10.0 {
		fmt.Println("GOOD LANDING (COULD BE BETTER)")
	} else if velocityMph <= 60.0 {
		fmt.Println("CRAFT DAMAGE... YOU'RE STRANDED HERE")
		fmt.Println("UNTIL A RESCUE PARTY ARRIVES.")
		fmt.Println("I HOPE YOU HAVE ENOUGH OXYGEN!")
	} else {
		fmt.Println("THAT'S ONE SMALL IMPACT FOR THE MOON,")
		fmt.Println("ONE GIANT BOOM FOR YOUR LANDER!")
		fmt.Printf("YOU BLASTED A NEW CRATER %0.0f FEET DEEP!",
			velocityMph*0.227)
	}
}

func intro() {
	fmt.Println("")
	fmt.Println("              LUNAR")
	fmt.Println("CREATIVE COMPUTING MORRISTOWN, NJ")
	fmt.Println("")
	fmt.Println("THIS IS A COMPUTER SIMULATION OF AN")
	fmt.Println("APOLLO LUNAR LANDING CAPSULE.")
	fmt.Println("")
	fmt.Println("THE ON-BOARD COMPUTER HAS FAILED (IT WAS")
	fmt.Println("MADE BY XEROX) SO YOU HAVE TO LAND THE")
	fmt.Println("CAPSULE MANUALLY.")
	fmt.Println("")
}

func printHeader() {
	fmt.Println(" SEC  MILES  FEET    MPH    FUEL ")
}

func printStatus(lem *lander.Lander) {
	miles, fracFeet := math.Modf(lem.Altitude)
	feet := int(5280.0 * fracFeet)
	mph := int(3600.0 * lem.Velocity)
	fmt.Printf("%4.0f    %3.0f  %4d  %5d   %5.0f  ",
		lem.ElapsedTime, miles, feet, mph, lem.Fuel)
}

func getInput(in *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to read input: ", err)
		os.Exit(1)
	}
	return input
}

func getBurnRate(in *bufio.Reader) float64 {
	var burnRate float64
	for {
		var err error
		input := getInput(in, "RATE? ")
		burnRate, err = strconv.ParseFloat(strings.TrimSpace(input), 64)
		if err != nil {
			fmt.Print("PLEASE ENTER A BURN RATE: ")
			continue
		}
		if burnRate < 0.0 || burnRate > 200.0 {
			fmt.Println("PLEASE ENTER A BURN RATE")
			fmt.Print("BETWEEN 0 AND 200: ")
			continue
		}
		break
	}
	return burnRate
}

func printStartMessage() {
	fmt.Println("")
	fmt.Println("SET THE BURN RATE OF THE RETRO ROCKETS")
	fmt.Println("TO ANY VALUE BETWEEN 0 (FREE FALL) AND")
	fmt.Println("200 (MAXIMUM BURN) IN POUNDS PER SECOND.")
	fmt.Println("SET A NEW BURN RATE EVERY 10 SECONDS.")
	fmt.Println("")
	fmt.Println("CAPSULE DRY WEIGHT IS 16,500 LBS;")
	fmt.Println("INITIAL FUEL IS 16,500 LBS.")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("GOOD LUCK!")
	fmt.Println("")
}

func runGame(inputSource *bufio.Reader) {

	// These should be options, passed to this function as a struct...
	var lem = lander.Lander{
		CapsuleMass: 16500.0,
		Fuel:        16500.0,
		Kinematics: lander.Kinematics{
			Altitude:    120.0,
			Velocity:    1.0,
			ElapsedTime: 0.0,
		},
	}
	lem.TotalMass = lem.CapsuleMass + lem.Fuel

	var thisPeriod, burnRate, burnTime float64
	var newK lander.Kinematics

	printStartMessage()
	printHeader()

game:
	for {
		printStatus(&lem)
		burnRate = getBurnRate(inputSource)
		thisPeriod = 10.0
		for {
			if lem.OutOfFuel() {
				outOfFuel(&lem)
				break game
			}
			if thisPeriod < 0.001 {
				break
			}
			burnTime = lem.ActualBurnTime(burnRate, thisPeriod)
			newK = lem.CalcDynamics(burnRate, burnTime)
			if newK.Altitude <= 0.0 {
				landing(&lem, burnRate, burnTime)
				// landing() does not itself end the game
				break game
			}
			if lem.Velocity > 0.0 {
				if newK.Velocity < 0.0 {
					burnTime = lem.UpwardBurnTime(burnRate)
					newK = lem.CalcDynamics(burnRate, burnTime)
					if newK.Altitude <= 0.0 {
						landing(&lem, burnRate, burnTime)
						break game
					}
				}
			}
			lem.UpdateLander(burnRate, burnTime, newK)
			thisPeriod -= burnTime
		}
	}
}

func main() {

	var response string
	fromStdin := bufio.NewReader(os.Stdin)

	intro()
	for {
		runGame(fromStdin)
		response = getInput(fromStdin, "\nTRY AGAIN?? ")
		response = strings.ToLower(strings.TrimSpace(response))
		if !strings.HasPrefix(response, "y") {
			break
		}
	}
}
