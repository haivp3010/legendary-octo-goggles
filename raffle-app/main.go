package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Ticket represents a raffle ticket with a set of numbers.
type Ticket struct {
	Numbers []int
}

func (t Ticket) PrintNumbers() string {
	var nums []string
	for _, num := range t.Numbers {
		nums = append(nums, strconv.Itoa(num))
	}
	return strings.Join(nums, " ")
}

// User represents a raffle participant with purchased tickets.
type User struct {
	Name    string
	Tickets []Ticket
}

// Raffle represents the state of a raffle draw.
type Raffle struct {
	rand    *rand.Rand
	scanner *bufio.Scanner
	Open    bool
	PotSize float64
	Users   []User
	Winner  Ticket
}

// Rewards represents the reward structure for different prize groups.
type Rewards struct {
	Group2 float64
	Group3 float64
	Group4 float64
	Group5 float64
}

func main() {
	currentRaffle := Raffle{
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		scanner: bufio.NewScanner(os.Stdin),
	}

	for {
		clear()
		fmt.Println("Welcome to My Raffle App")
		fmt.Printf("Status: %s\n\n", getRaffleStatus(currentRaffle))
		fmt.Println("[1] Start a New Draw")
		fmt.Println("[2] Buy Tickets")
		fmt.Println("[3] Run Raffle")
		fmt.Println()

		var choice string
		fmt.Print("Enter your choice: ")
		currentRaffle.scanner.Scan()
		choice = currentRaffle.scanner.Text()

		switch choice {
		case "1":
			currentRaffle = StartNewDraw(currentRaffle)
		case "2":
			currentRaffle = BuyTickets(currentRaffle)
		case "3":
			RunRaffle(&currentRaffle)
		default:
			fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
		}
	}
}

func StartNewDraw(raffle Raffle) Raffle {
	clear()
	raffle = startNewDraw(raffle)
	fmt.Printf("New Raffle draw has been started. Initial pot size: $%s\n", trimTrailingZeroes(raffle.PotSize))
	fmt.Print("Press any key to return to the main menu")
	keyboard.GetSingleKey()
	return raffle
}

func BuyTickets(currentRaffle Raffle) Raffle {
	clear()
	if !currentRaffle.Open {
		fmt.Println("Draw has not started")
	} else {

		fmt.Print("Enter your name, number of tickets to purchase (e.g., James,1): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		var (
			user       User
			numTickets int
		)
		currentRaffle, user, numTickets = buyTicketsForUser(currentRaffle, input)
		if numTickets > 0 {
			fmt.Printf("Hi %s, you have purchased %d ticket%s\n", user.Name, numTickets, func() string {
				if numTickets > 1 {
					return "s"
				}
				return ""
			}())
			for i, ticket := range user.Tickets {
				fmt.Printf("Ticket %d: %s\n", i+1, ticket.PrintNumbers())
			}
			fmt.Println()
			currentRaffle.Users = append(currentRaffle.Users, user)
		}
	}

	fmt.Print("Press any key to return to the main menu")
	keyboard.GetSingleKey()

	return currentRaffle
}

func RunRaffle(currentRaffle *Raffle) {
	clear()
	runRaffle(currentRaffle)

	fmt.Print("Press any key to return to the main menu")
	keyboard.GetSingleKey()
}

func generateTicket(raffle Raffle) Ticket {
	uniqueNumbers := make([]int, 0, 5)
	usedNumbers := make(map[int]bool)

	for len(uniqueNumbers) < 5 {
		num := raffle.rand.Intn(15) + 1

		if !usedNumbers[num] {
			usedNumbers[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}
	return Ticket{Numbers: uniqueNumbers}
}

func getWinners(currentRaffle *Raffle) ([]User, []User, []User, []User) {
	var group2Winners, group3Winners, group4Winners, group5Winners []User

	for _, user := range currentRaffle.Users {
		matchedCount := countMatchingNumbers(user.Tickets[0].Numbers, currentRaffle.Winner.Numbers)
		switch matchedCount {
		case 2:
			group2Winners = append(group2Winners, user)
		case 3:
			group3Winners = append(group3Winners, user)
		case 4:
			group4Winners = append(group4Winners, user)
		case 5:
			group5Winners = append(group5Winners, user)
		}
	}

	return group2Winners, group3Winners, group4Winners, group5Winners
}

func countMatchingNumbers(ticketNumbers, winningNumbers []int) int {
	matchingNumbers := 0
	winningNumCount := make(map[int]int)

	for _, num := range winningNumbers {
		winningNumCount[num]++
	}

	for _, num := range ticketNumbers {
		if count, ok := winningNumCount[num]; ok && count > 0 {
			matchingNumbers++
			winningNumCount[num]--
		}
	}

	return matchingNumbers
}

func calculateRewards(currentRaffle Raffle) Rewards {
	totalPot := currentRaffle.PotSize
	return Rewards{
		Group2: 0.1 * totalPot,
		Group3: 0.15 * totalPot,
		Group4: 0.25 * totalPot,
		Group5: 0.5 * totalPot,
	}
}

func calculateTotalRewards(rewards Rewards, group2Winners, group3Winners, group4Winners, group5Winners []User) float64 {
	var total float64
	if len(group2Winners) > 0 {
		total += rewards.Group2
	}
	if len(group3Winners) > 0 {
		total += rewards.Group3
	}
	if len(group4Winners) > 0 {
		total += rewards.Group4
	}
	if len(group5Winners) > 0 {
		total += rewards.Group5
	}
	return total
}

func displayWinners(groupName string, winners []User, rewardPercentage float64) {
	fmt.Printf("%s:\n", groupName)
	if len(winners) == 0 {
		fmt.Println("Nil")
	}
	for _, winner := range winners {
		fmt.Printf("%s with %d winning ticket(s)- $%s\n", winner.Name, len(winner.Tickets), trimTrailingZeroes(rewardPercentage/float64(len(winners))))
	}
	fmt.Println()
}

func getRaffleStatus(currentRaffle Raffle) string {
	if !currentRaffle.Open {
		return "Draw has not started"
	}
	return fmt.Sprintf("Draw is ongoing. Raffle pot size is $%s", trimTrailingZeroes(currentRaffle.PotSize))
}

func startNewDraw(raffle Raffle) Raffle {
	if !raffle.Open {
		raffle.PotSize += 100
		raffle.Open = true
	}
	return raffle
}

func buyTicketsForUser(currentRaffle Raffle, input string) (Raffle, User, int) {
	var user User
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		fmt.Println("Invalid input format. Please enter name and number of tickets.")
		return currentRaffle, user, 0
	}

	user.Name = parts[0]
	numTickets := parseNumber(parts[1])

	for i := 0; i < numTickets; i++ {
		ticket := generateTicket(currentRaffle)
		user.Tickets = append(user.Tickets, ticket)
		currentRaffle.PotSize += 5 // Each ticket costs $5
	}
	return currentRaffle, user, numTickets
}

func runRaffle(currentRaffle *Raffle) {
	if !currentRaffle.Open {
		fmt.Println("Draw has not started")
	} else {
		fmt.Println("Running Raffle..")
		rewards := calculateRewards(*currentRaffle)

		currentRaffle.Winner = generateTicket(*currentRaffle)
		fmt.Printf("Winning Ticket is %s\n", currentRaffle.Winner.PrintNumbers())

		group2Winners, group3Winners, group4Winners, group5Winners := getWinners(currentRaffle)

		displayWinners("Group 2 Winners", group2Winners, rewards.Group2)
		displayWinners("Group 3 Winners", group3Winners, rewards.Group3)
		displayWinners("Group 4 Winners", group4Winners, rewards.Group4)
		displayWinners("Group 5 Winners (Jackpot)", group5Winners, rewards.Group5)

		// Reset the raffle state for the next draw
		currentRaffle.PotSize = currentRaffle.PotSize - calculateTotalRewards(rewards, group2Winners, group3Winners, group4Winners, group5Winners)
		currentRaffle.Open = false
	}
}

func parseNumber(s string) int {
	s = strings.TrimSpace(s)
	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error parsing number:", err)
	}
	return num
}

func trimTrailingZeroes(num float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", num), "0"), ".")
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clear() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}
