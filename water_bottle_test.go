package main

import (
	"fmt"
	"testing"
)

func TestWaterBottleGameCreation(t *testing.T) {
	// Test valid game creation
	game, err := NewWaterBottleGame(5, 4, 2, 3)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if game.N != 5 || game.M != 4 || game.J != 2 || game.K != 3 {
		t.Errorf("Game parameters not set correctly")
	}

	// Test invalid parameters
	_, err = NewWaterBottleGame(2, 4, 2, 3) // N <= J
	if err == nil {
		t.Error("Expected error for N <= J")
	}

	_, err = NewWaterBottleGame(5, 0, 2, 3) // M <= 0
	if err == nil {
		t.Error("Expected error for M <= 0")
	}

	_, err = NewWaterBottleGame(5, 4, 2, 0) // K <= 0
	if err == nil {
		t.Error("Expected error for K <= 0")
	}
}

func TestPourLogic(t *testing.T) {
	game, _ := NewWaterBottleGame(4, 3, 1, 2)

	// Manually set up a simple test state
	game.bottles[0] = Bottle{Color(0), Color(0), Color(1)} // Red, Red, Blue
	game.bottles[1] = Bottle{Color(1), Color(1)}           // Blue, Blue
	game.bottles[2] = Bottle{}                             // Empty
	game.bottles[3] = Bottle{}                             // Empty
	game.emptyCount = 2

	// Test pouring from bottle 0 to empty bottle 2
	success, moved := game.Pour(0, 2)
	if !success || moved != 1 {
		t.Errorf("Expected success=true, moved=1, got success=%v, moved=%d", success, moved)
	}

	// Check result
	if len(game.bottles[0]) != 2 || game.bottles[0][1] != Color(0) {
		t.Error("Source bottle state incorrect after pour")
	}
	if len(game.bottles[2]) != 1 || game.bottles[2][0] != Color(1) {
		t.Error("Target bottle state incorrect after pour")
	}

	// Test pouring same color
	success, moved = game.Pour(1, 2) // Blue to Blue
	if !success || moved != 2 {
		t.Errorf("Expected success=true, moved=2, got success=%v, moved=%d", success, moved)
	}

	// Test invalid pour (different colors)
	success, moved = game.Pour(0, 2) // Red to Blue
	if success {
		t.Error("Expected failure when pouring different colors")
	}
}

func TestGameStateGeneration(t *testing.T) {
	game, err := NewWaterBottleGame(4, 3, 1, 2)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	err = game.generateInitialState()
	if err != nil {
		t.Fatalf("Failed to generate initial state: %v", err)
	}

	// Verify total water amount
	totalWater := 0
	for _, bottle := range game.bottles {
		totalWater += len(bottle)
	}

	expectedWater := (game.N - game.J) * game.M
	if totalWater != expectedWater {
		t.Errorf("Expected total water %d, got %d", expectedWater, totalWater)
	}

	// Verify color distribution
	colorCounts := make(map[Color]int)
	for _, bottle := range game.bottles {
		for _, color := range bottle {
			colorCounts[color]++
		}
	}

	// Each color count should be divisible by M
	for color, count := range colorCounts {
		if count%game.M != 0 {
			t.Errorf("Color %d count (%d) not divisible by M (%d)", color, count, game.M)
		}
	}
}

func TestWinCondition(t *testing.T) {
	game, _ := NewWaterBottleGame(3, 2, 1, 2)

	// Set up a winning state
	game.bottles[0] = Bottle{Color(0), Color(0)} // Full bottle of color 0
	game.bottles[1] = Bottle{Color(1), Color(1)} // Full bottle of color 1
	game.bottles[2] = Bottle{}                   // Empty bottle
	game.emptyCount = 1

	if !game.IsWon() {
		t.Error("Expected winning state")
	}

	// Set up a non-winning state (mixed colors)
	game.bottles[0] = Bottle{Color(0), Color(1)} // Mixed colors
	if game.IsWon() {
		t.Error("Expected non-winning state for mixed colors")
	}

	// Set up a non-winning state (not full)
	game.bottles[0] = Bottle{Color(0)} // Not full
	game.bottles[1] = Bottle{Color(1), Color(1)}
	if game.IsWon() {
		t.Error("Expected non-winning state for non-full bottle")
	}
}

// Benchmark the initial state generation
func BenchmarkGenerateInitialState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game, _ := NewWaterBottleGame(6, 4, 2, 3)
		game.generateInitialState()
	}
}

// Example test that demonstrates usage
func ExampleWaterBottleGame() {
	// Create a small game
	game, _ := NewWaterBottleGame(4, 3, 1, 2)

	fmt.Printf("Created game with %d bottles, capacity %d each\n", game.N, game.M)
	fmt.Printf("Empty bottles: %d, Colors: %d\n", game.J, game.K)

	// Output:
	// Created game with 4 bottles, capacity 3 each
	// Empty bottles: 1, Colors: 2
}
