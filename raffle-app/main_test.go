package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWinners(t *testing.T) {
	tests := []struct {
		name                  string
		currentRaffle         Raffle
		expectedGroup2Winners []User
		expectedGroup3Winners []User
		expectedGroup4Winners []User
		expectedGroup5Winners []User
	}{
		{
			name: "No Winners",
			currentRaffle: Raffle{
				PotSize: 110,
				Users: []User{
					{Name: "User1", Tickets: []Ticket{{Numbers: []int{1, 2, 3, 4, 5}}}},
					{Name: "User2", Tickets: []Ticket{{Numbers: []int{6, 7, 8, 9, 10}}}},
				},
				Winner: Ticket{Numbers: []int{11, 12, 13, 14, 15}},
			},
			expectedGroup2Winners: []User{},
			expectedGroup3Winners: []User{},
			expectedGroup4Winners: []User{},
			expectedGroup5Winners: []User{},
		},
		{
			name: "Winners",
			currentRaffle: Raffle{
				PotSize: 125,
				Users: []User{
					{Name: "User1", Tickets: []Ticket{{Numbers: []int{1, 2, 3, 4, 5}}}},
					{Name: "User2", Tickets: []Ticket{{Numbers: []int{2, 7, 8, 9, 10}}}},
					{Name: "User3", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 9, 10}}}},
					{Name: "User4", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 10, 11}}}},
					{Name: "User5", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 11, 12}}}},
				},
				Winner: Ticket{Numbers: []int{3, 7, 8, 11, 12}},
			},
			expectedGroup2Winners: []User{{Name: "User2", Tickets: []Ticket{{Numbers: []int{2, 7, 8, 9, 10}}}}},
			expectedGroup3Winners: []User{{Name: "User3", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 9, 10}}}}},
			expectedGroup4Winners: []User{{Name: "User4", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 10, 11}}}}},
			expectedGroup5Winners: []User{{Name: "User5", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 11, 12}}}}},
		},
		{
			name: "Co-winners",
			currentRaffle: Raffle{
				PotSize: 110,
				Users: []User{
					{Name: "User4", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 10, 11}}}},
					{Name: "User5", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 12, 13}}}},
				},
				Winner: Ticket{Numbers: []int{3, 7, 8, 11, 12}},
			},
			expectedGroup4Winners: []User{
				{Name: "User4", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 10, 11}}}},
				{Name: "User5", Tickets: []Ticket{{Numbers: []int{3, 7, 8, 12, 13}}}}},
			expectedGroup5Winners: []User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group2Winners, group3Winners, group4Winners, group5Winners := getWinners(&tt.currentRaffle)
			assert.ElementsMatch(t, group2Winners, tt.expectedGroup2Winners)
			assert.ElementsMatch(t, group3Winners, tt.expectedGroup3Winners)
			assert.ElementsMatch(t, group4Winners, tt.expectedGroup4Winners)
			assert.ElementsMatch(t, group5Winners, tt.expectedGroup5Winners)
		})
	}
}

func TestRunRaffle(t *testing.T) {
	tests := []struct {
		name              string
		raffle            Raffle
		expectedRewards   map[string]float64
		expectedRemaining float64
	}{
		{
			name: "Test not open",
			raffle: Raffle{
				rand:    rand.New(rand.NewSource(1)),
				PotSize: 110,
				Open:    false,
				Users: []User{
					{Name: "User4", Tickets: []Ticket{{Numbers: []int{12, 13, 3, 15, 2}}}},
					{Name: "User5", Tickets: []Ticket{{Numbers: []int{12, 13, 3, 15, 2}}}},
				},
			},
			expectedRewards:   map[string]float64{},
			expectedRemaining: 110,
		},
		{
			name: "Test no winners",
			raffle: Raffle{
				rand:    rand.New(rand.NewSource(1)),
				PotSize: 110,
				Open:    true,
				Users: []User{
					{Name: "User4", Tickets: []Ticket{{Numbers: []int{1, 2, 4, 5, 6}}}},
					{Name: "User5", Tickets: []Ticket{{Numbers: []int{1, 2, 4, 5, 6}}}},
				},
			},
			expectedRewards:   map[string]float64{},
			expectedRemaining: 110,
		},
		{
			name: "Test remaining pot",
			raffle: Raffle{
				rand:    rand.New(rand.NewSource(1)),
				PotSize: 110,
				Open:    true,
				Users: []User{
					{Name: "User4", Tickets: []Ticket{{Numbers: []int{12, 13, 3, 15, 2}}}},
					{Name: "User5", Tickets: []Ticket{{Numbers: []int{12, 13, 3, 15, 2}}}},
				},
			},
			expectedRewards:   map[string]float64{},
			expectedRemaining: 55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runRaffle(&tt.raffle)

			assert.Equal(t, tt.expectedRemaining, tt.raffle.PotSize)
		})
	}
}

func TestBuyTicketsForUser(t *testing.T) {
	// Mock Raffle for testing
	mockRaffle := Raffle{
		rand:    rand.New(rand.NewSource(1)),
		PotSize: 100,
		Users:   []User{},
		Winner:  Ticket{},
	}

	tests := []struct {
		name               string
		input              string
		expectedRaffle     Raffle
		expectedUser       User
		expectedNumTickets int
	}{
		{
			name:           "Valid Input",
			input:          "John,3",
			expectedRaffle: Raffle{PotSize: 115}, // Initial pot size + 3 tickets * $5
			expectedUser: User{Name: "John", Tickets: []Ticket{{
				Numbers: []int{12, 13, 3, 15, 2},
			}, {
				Numbers: []int{4, 11, 6, 2, 1},
			}, {
				Numbers: []int{15, 2, 13, 14, 7},
			}}}, // 3 tickets generated
			expectedNumTickets: 3,
		},
		{
			name:               "Invalid Input",
			input:              "InvalidInput",
			expectedRaffle:     Raffle{PotSize: 100}, // Pot size remains unchanged
			expectedUser:       User{},               // No tickets generated
			expectedNumTickets: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultRaffle, resultUser, resultNumTickets := buyTicketsForUser(mockRaffle, tt.input)
			assert.Equal(t, tt.expectedRaffle.PotSize, resultRaffle.PotSize)
			assert.Equal(t, resultUser, tt.expectedUser)
			assert.Equal(t, resultNumTickets, tt.expectedNumTickets)
		})
	}
}

func TestStartNewDraw(t *testing.T) {
	tests := []struct {
		name           string
		input          Raffle
		expectedRaffle Raffle
	}{
		{
			name: "Test new draw",
			input: Raffle{
				PotSize: 100,
				Open:    false,
			},
			expectedRaffle: Raffle{PotSize: 200, Open: true},
		},
		{
			name: "Test already open draw",
			input: Raffle{
				PotSize: 100,
				Open:    true,
			},
			expectedRaffle: Raffle{PotSize: 100, Open: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultRaffle := startNewDraw(tt.input)
			assert.Equal(t, tt.expectedRaffle.PotSize, resultRaffle.PotSize)
			assert.Equal(t, resultRaffle.Open, tt.expectedRaffle.Open)
		})
	}
}
