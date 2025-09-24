package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Color represents a water color (0-based)
type Color int

// Bottle represents a water bottle with layers of colored water
type Bottle []Color

// Move represents a single pour operation
type Move struct {
	From   int // Source bottle index
	To     int // Target bottle index
	Amount int // Amount of water moved
	Color  Color // Color of water moved
}

// WaterBottleGame represents the game state
type WaterBottleGame struct {
	bottles       []Bottle // All bottles in the game
	N             int      // Total number of bottles
	M             int      // Capacity of each bottle
	J             int      // Number of empty bottles
	K             int      // Number of different colors
	emptyCount    int      // Current number of empty bottles
	reverseSteps  []Move   // Record of reverse operations for validation
}

// NewWaterBottleGame creates a new game with given parameters
func NewWaterBottleGame(N, M, J, K int) (*WaterBottleGame, error) {
	if N <= J {
		return nil, fmt.Errorf("total bottles (%d) must be greater than empty bottles (%d)", N, J)
	}
	if K <= 0 {
		return nil, fmt.Errorf("number of colors (%d) must be positive", K)
	}
	if M <= 0 {
		return nil, fmt.Errorf("bottle capacity (%d) must be positive", M)
	}

	totalWater := (N - J) * M
	if totalWater%M != 0 {
		return nil, fmt.Errorf("total water volume must be divisible by bottle capacity")
	}

	game := &WaterBottleGame{
		bottles:    make([]Bottle, N),
		N:          N,
		M:          M,
		J:          J,
		K:          K,
		emptyCount: J,
	}

	// Initialize empty bottles
	for i := range game.bottles {
		game.bottles[i] = make(Bottle, 0, M)
	}

	return game, nil
}

// generateInitialState creates a solvable initial game state using reverse generation
func (g *WaterBottleGame) generateInitialState() error {
	// Use default difficulty calculation
	difficulty := g.calculateDifficulty()
	return g.generateInitialStateWithSteps(difficulty)
}

// generateInitialStateWithSteps creates initial state with specified reverse steps
func (g *WaterBottleGame) generateInitialStateWithSteps(reverseSteps int) error {
	rand.Seed(time.Now().UnixNano())

	// Check if parameters are reasonable
	totalWater := (g.N - g.J) * g.M
	maxPossibleColors := totalWater / g.M
	if g.K > maxPossibleColors {
		return fmt.Errorf("参数不合理：总水量%d，每种颜色至少需要%d单位，最多只能有%d种颜色，但要求%d种",
			totalWater, g.M, maxPossibleColors, g.K)
	}

	// Use reverse generation: start from solved state and work backwards
	return g.generateByReverseWithSteps(reverseSteps)
}

// generateByReverse creates initial state by working backwards from solved state
func (g *WaterBottleGame) generateByReverse() error {
	difficulty := g.calculateDifficulty()
	return g.generateByReverseWithSteps(difficulty)
}

// generateByReverseWithSteps creates initial state with specified reverse steps
func (g *WaterBottleGame) generateByReverseWithSteps(reverseSteps int) error {
	// Step 1: Create perfect solved state
	if err := g.createSolvedState(); err != nil {
		return err
	}
	
	// Initialize reverse steps recording
	g.reverseSteps = make([]Move, 0, reverseSteps)
	
	// Step 2: Apply reverse operations to create puzzle
	fmt.Printf("🎲 正在进行 %d 步逆向操作生成谜题...\n", reverseSteps)
	
	actualSteps := 0
	totalAttempts := 0
	
	for step := 0; step < reverseSteps; step++ {
		if actualSteps > 0 && actualSteps%20 == 0 {
			fmt.Printf("   进度: %d 有效步数 (尝试了 %d 次)\n", actualSteps, totalAttempts)
		}
		
		// Try to find a valid reverse operation
		maxAttempts := 50
		success := false
		
		for attempt := 0; attempt < maxAttempts; attempt++ {
			totalAttempts++
			if g.tryReverseOperationWithRecord() {
				success = true
				actualSteps++
				break
			}
		}
		
		if !success {
			// If we can't find more reverse operations, we're done
			fmt.Printf("   在第 %d 有效步数完成（总共尝试 %d 次，无法继续逆向操作）\n", actualSteps, totalAttempts)
			break
		}
	}
	
	successRate := float64(actualSteps) / float64(totalAttempts) * 100
	fmt.Printf("🎯 逆向生成完成！\n")
	fmt.Printf("   - 有效逆向操作: %d 步\n", actualSteps)
	fmt.Printf("   - 总尝试次数: %d 次\n", totalAttempts)
	fmt.Printf("   - 成功率: %.1f%%\n", successRate)
	
	// Step 3: Validate that we can restore the original state
	if err := g.validateReverseSteps(); err != nil {
		return fmt.Errorf("逆向步骤验证失败: %v", err)
	}
	
	fmt.Println("✅ 逆向步骤验证成功！所有操作都可以还原")
	return nil
}

// createSolvedState creates the perfect solved state
func (g *WaterBottleGame) createSolvedState() error {
	// Calculate how many bottles each color needs
	baseBottlesPerColor := (g.N - g.J) / g.K
	extraBottles := (g.N - g.J) % g.K

	bottleIndex := 0
	for colorID := 0; colorID < g.K; colorID++ {
		bottlesForThisColor := baseBottlesPerColor
		if colorID < extraBottles {
			bottlesForThisColor++
		}

		// Fill bottles with this color
		for b := 0; b < bottlesForThisColor; b++ {
			if bottleIndex >= g.N-g.J {
				return fmt.Errorf("bottle index overflow")
			}

			// Fill bottle completely with single color
			for i := 0; i < g.M; i++ {
				g.bottles[bottleIndex] = append(g.bottles[bottleIndex], Color(colorID))
			}
			bottleIndex++
		}
	}

	// Remaining bottles are empty
	for i := bottleIndex; i < g.N; i++ {
		g.bottles[i] = make(Bottle, 0, g.M)
	}

	g.emptyCount = g.J
	return nil
}

// calculateDifficulty determines how many reverse steps to take
func (g *WaterBottleGame) calculateDifficulty() int {
	// Base difficulty on game complexity
	totalBottles := g.N
	totalColors := g.K
	capacity := g.M

	// More bottles/colors/capacity = more possible moves = higher difficulty
	baseDifficulty := totalBottles * totalColors * capacity / 4

	// Add some randomness (±25%)
	variation := baseDifficulty / 4
	difficulty := baseDifficulty + rand.Intn(variation*2+1) - variation

	// Ensure reasonable bounds
	minDifficulty := max(10, totalBottles*2)
	maxDifficulty := totalBottles * totalColors * capacity

	if difficulty < minDifficulty {
		difficulty = minDifficulty
	}
	if difficulty > maxDifficulty {
		difficulty = maxDifficulty
	}

	return difficulty
}

// tryReverseOperation attempts one simple reverse operation (random pour)
func (g *WaterBottleGame) tryReverseOperation() bool {
	return g.tryReverseOperationWithRecord()
}

// tryReverseOperationWithRecord attempts reverse operation and records the move
func (g *WaterBottleGame) tryReverseOperationWithRecord() bool {
	// Find all non-empty bottles as potential sources
	var sources []int
	for i, bottle := range g.bottles {
		if len(bottle) > 0 {
			sources = append(sources, i)
		}
	}
	
	if len(sources) == 0 {
		return false // No water to pour
	}
	
	// Try multiple random combinations
	maxAttempts := 20
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Pick random source
		sourceIdx := sources[rand.Intn(len(sources))]
		sourceBottle := g.bottles[sourceIdx]
		
		if len(sourceBottle) == 0 {
			continue
		}
		
		// Get the top color
		topColor := sourceBottle[len(sourceBottle)-1]
		
		// Count how many of this color are on top
		maxPourAmount := 0
		for j := len(sourceBottle) - 1; j >= 0 && sourceBottle[j] == topColor; j-- {
			maxPourAmount++
		}
		
		// Pick random amount to pour (1 to maxPourAmount)
		pourAmount := rand.Intn(maxPourAmount) + 1
		
		// Find all valid targets (reverse operation can pour anywhere with space!)
		var targets []int
		for i, bottle := range g.bottles {
			if i == sourceIdx {
				continue // Can't pour to self
			}

			// In reverse operation, we can pour to ANY bottle with space
			hasSpace := len(bottle)+pourAmount <= g.M

			if hasSpace {
				targets = append(targets, i)
			}
		}
		
		if len(targets) > 0 {
			// Pick random target and perform the pour
			targetIdx := targets[rand.Intn(len(targets))]
			
			// Save current state before attempting the move
			stateBefore := g.copyGameState()
			
			// Perform the reverse move
			if g.performSimplePour(sourceIdx, targetIdx, pourAmount) {
				// Immediately test if this move can be reversed using forward rules
				canReverse, _ := g.Pour(targetIdx, sourceIdx)
				
				if canReverse {
					// The move is valid! Restore state and record it
					g.restoreGameState(stateBefore)
					
					// Re-perform the reverse move and record it
					if g.performSimplePour(sourceIdx, targetIdx, pourAmount) {
						move := Move{
							From:   sourceIdx,
							To:     targetIdx,
							Amount: pourAmount,
							Color:  topColor,
						}
						g.reverseSteps = append(g.reverseSteps, move)
						return true
					}
				} else {
					// This move cannot be reversed, restore state and try again
					g.restoreGameState(stateBefore)
					continue // Try next target or next attempt
				}
			}
		}
	}
	
	return false // Couldn't find any valid pour
}

// performSimplePour executes a simple pour operation (used in reverse generation)
func (g *WaterBottleGame) performSimplePour(fromIdx, toIdx, amount int) bool {
	if fromIdx < 0 || fromIdx >= len(g.bottles) ||
		toIdx < 0 || toIdx >= len(g.bottles) ||
		amount <= 0 {
		return false
	}

	sourceBottle := &g.bottles[fromIdx]
	targetBottle := &g.bottles[toIdx]

	// Basic validation
	if len(*sourceBottle) < amount {
		return false
	}

	if len(*targetBottle)+amount > g.M {
		return false
	}

	// Get the color we're moving
	color := (*sourceBottle)[len(*sourceBottle)-1]

	// Verify we have enough of this color on top
	colorCount := 0
	for j := len(*sourceBottle) - 1; j >= 0 && (*sourceBottle)[j] == color; j-- {
		colorCount++
		if colorCount >= amount {
			break
		}
	}

	if colorCount < amount {
		return false
	}

	// In reverse generation, we allow pouring any color onto any color
	// This creates mixed states that can be solved using forward rules

	// Track empty bottle count changes
	wasSourceEmpty := len(*sourceBottle) == 0
	wasTargetEmpty := len(*targetBottle) == 0

	// Perform the pour
	for i := 0; i < amount; i++ {
		*targetBottle = append(*targetBottle, color)
	}
	*sourceBottle = (*sourceBottle)[:len(*sourceBottle)-amount]

	// Update empty count
	nowSourceEmpty := len(*sourceBottle) == 0
	nowTargetEmpty := len(*targetBottle) == 0

	if !wasSourceEmpty && nowSourceEmpty {
		g.emptyCount++
	}
	if wasTargetEmpty && !nowTargetEmpty {
		g.emptyCount--
	}

	return true
}

// validateReverseSteps verifies that all reverse steps can be undone to restore solved state
func (g *WaterBottleGame) validateReverseSteps() error {
	// Save current state
	currentState := g.copyGameState()
	
	// Apply reverse steps in reverse order (forward direction)
	fmt.Printf("🔍 验证 %d 步逆向操作的可还原性...\n", len(g.reverseSteps))
	
	for i := len(g.reverseSteps) - 1; i >= 0; i-- {
		move := g.reverseSteps[i]
		
		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败", 
				len(g.reverseSteps)-i, move.To, move.From)
		}
		
		if (i+1)%50 == 0 {
			fmt.Printf("   验证进度: %d/%d\n", len(g.reverseSteps)-i, len(g.reverseSteps))
		}
	}
	
	// Check if we're back to solved state
	if !g.IsWon() {
		return fmt.Errorf("还原后的状态不是完美解状态")
	}
	
	// Restore the generated initial state
	g.restoreGameState(currentState)
	return nil
}

// copyGameState creates a deep copy of the current game state
func (g *WaterBottleGame) copyGameState() [][]Color {
	state := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		state[i] = make([]Color, len(bottle))
		copy(state[i], bottle)
	}
	return state
}

// restoreGameState restores the game to a previous state
func (g *WaterBottleGame) restoreGameState(state [][]Color) {
	for i, bottleState := range state {
		g.bottles[i] = make(Bottle, len(bottleState))
		copy(g.bottles[i], bottleState)
	}
	
	// Recalculate empty count
	g.emptyCount = 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			g.emptyCount++
		}
	}
}

// GetReverseSteps returns the recorded reverse steps for analysis  
func (g *WaterBottleGame) GetReverseSteps() []Move {
	return g.reverseSteps
}

// getColorName returns the Chinese name of a color for display
func getColorName(color Color) string {
	names := []string{"红", "蓝", "绿", "黄", "橙", "紫", "棕", "黑", "白", "粉"}
	if int(color) < len(names) {
		return names[color]
	}
	return fmt.Sprintf("色%d", color)
}

// tryAggressivePour performs random pours to create more mixing
func (g *WaterBottleGame) tryAggressivePour() bool {
	// Find all bottles with water
	var sources []int
	for i, bottle := range g.bottles {
		if len(bottle) > 0 {
			sources = append(sources, i)
		}
	}

	if len(sources) == 0 {
		return false
	}

	// Try multiple random sources
	for attempts := 0; attempts < 5; attempts++ {
		sourceIdx := sources[rand.Intn(len(sources))]
		sourceBottle := g.bottles[sourceIdx]

		if len(sourceBottle) == 0 {
			continue
		}

		// Get top color and its continuous count
		topColor := sourceBottle[len(sourceBottle)-1]
		maxAmount := 0
		for j := len(sourceBottle) - 1; j >= 0 && sourceBottle[j] == topColor; j-- {
			maxAmount++
		}

		// Pour a random amount (1 to all continuous)
		pourAmount := rand.Intn(maxAmount) + 1

		// Find all possible targets (including creating more mixed states)
		var targets []int
		for i, bottle := range g.bottles {
			if i == sourceIdx {
				continue
			}

			// More liberal pouring rules for better mixing
			canPour := len(bottle) == 0 ||
				(len(bottle) > 0 && bottle[len(bottle)-1] == topColor) ||
				(len(bottle) < g.M/2) // Allow pouring to less-full bottles

			hasSpace := len(bottle)+pourAmount <= g.M

			if canPour && hasSpace {
				targets = append(targets, i)
			}
		}

		if len(targets) > 0 {
			targetIdx := targets[rand.Intn(len(targets))]
			targetBottle := g.bottles[targetIdx]

			// Only pour if colors match or target is empty
			if len(targetBottle) == 0 || targetBottle[len(targetBottle)-1] == topColor {
				return g.performPour(sourceIdx, targetIdx, pourAmount)
			}
		}
	}

	return false
}

// tryBreakCompleteBottle breaks up completed single-color bottles
func (g *WaterBottleGame) tryBreakCompleteBottle() bool {
	// Find complete single-color bottles
	var completeBottles []int
	for i, bottle := range g.bottles {
		if len(bottle) == g.M && g.isSingleColor(bottle) {
			completeBottles = append(completeBottles, i)
		}
	}

	if len(completeBottles) == 0 {
		return false
	}

	sourceIdx := completeBottles[rand.Intn(len(completeBottles))]
	sourceBottle := g.bottles[sourceIdx]
	color := sourceBottle[0]

	// Pour a significant amount (1/3 to 2/3 of bottle)
	minPour := g.M / 3
	maxPour := g.M * 2 / 3
	if minPour < 1 {
		minPour = 1
	}
	if maxPour > g.M-1 {
		maxPour = g.M - 1
	}

	pourAmount := rand.Intn(maxPour-minPour+1) + minPour

	// Find targets (prefer non-empty bottles for more mixing)
	var targets []int
	for i, bottle := range g.bottles {
		if i == sourceIdx {
			continue
		}

		hasSpace := len(bottle)+pourAmount <= g.M
		canPour := len(bottle) == 0 ||
			(len(bottle) > 0 && bottle[len(bottle)-1] == color)

		if hasSpace && canPour {
			targets = append(targets, i)
		}
	}

	if len(targets) > 0 {
		targetIdx := targets[rand.Intn(len(targets))]
		return g.performPour(sourceIdx, targetIdx, pourAmount)
	}

	return false
}

// tryOriginalSplit uses the original splitting strategy
func (g *WaterBottleGame) tryOriginalSplit() bool {
	var candidates []struct {
		bottleIdx  int
		moveAmount int
		color      Color
	}

	for i, bottle := range g.bottles {
		if len(bottle) <= 1 {
			continue
		}

		topColor := bottle[len(bottle)-1]
		blockSize := 0

		for j := len(bottle) - 1; j >= 0 && bottle[j] == topColor; j-- {
			blockSize++
		}

		if blockSize > 1 {
			for splitSize := 1; splitSize < blockSize; splitSize++ {
				candidates = append(candidates, struct {
					bottleIdx  int
					moveAmount int
					color      Color
				}{i, splitSize, topColor})
			}
		}
	}

	if len(candidates) == 0 {
		return false
	}

	candidate := candidates[rand.Intn(len(candidates))]

	var targets []int
	for i, bottle := range g.bottles {
		if i == candidate.bottleIdx {
			continue
		}

		canMove := len(bottle) == 0 ||
			(len(bottle) > 0 && bottle[len(bottle)-1] == candidate.color)
		hasSpace := len(bottle)+candidate.moveAmount <= g.M

		if canMove && hasSpace {
			targets = append(targets, i)
		}
	}

	if len(targets) > 0 {
		targetIdx := targets[rand.Intn(len(targets))]
		return g.performPour(candidate.bottleIdx, targetIdx, candidate.moveAmount)
	}

	return false
}

// performPour executes a pour operation safely
func (g *WaterBottleGame) performPour(fromIdx, toIdx, amount int) bool {
	if fromIdx < 0 || fromIdx >= len(g.bottles) ||
		toIdx < 0 || toIdx >= len(g.bottles) ||
		amount <= 0 {
		return false
	}

	sourceBottle := &g.bottles[fromIdx]
	targetBottle := &g.bottles[toIdx]

	if len(*sourceBottle) < amount {
		return false
	}

	// Get the color to move
	color := (*sourceBottle)[len(*sourceBottle)-1]

	// Verify we can move this amount of this color
	colorCount := 0
	for j := len(*sourceBottle) - 1; j >= 0 && (*sourceBottle)[j] == color; j-- {
		colorCount++
		if colorCount >= amount {
			break
		}
	}

	if colorCount < amount {
		return false
	}

	// Check target capacity and color compatibility
	if len(*targetBottle)+amount > g.M {
		return false
	}

	if len(*targetBottle) > 0 && (*targetBottle)[len(*targetBottle)-1] != color {
		return false
	}

	// Track empty state changes
	wasSourceEmpty := len(*sourceBottle) == 0
	wasTargetEmpty := len(*targetBottle) == 0

	// Perform the move
	for i := 0; i < amount; i++ {
		*targetBottle = append(*targetBottle, color)
	}
	*sourceBottle = (*sourceBottle)[:len(*sourceBottle)-amount]

	// Update empty count
	nowSourceEmpty := len(*sourceBottle) == 0
	nowTargetEmpty := len(*targetBottle) == 0

	if !wasSourceEmpty && nowSourceEmpty {
		g.emptyCount++
	}
	if wasTargetEmpty && !nowTargetEmpty {
		g.emptyCount--
	}

	return true
}

// validateReverseSteps verifies that all reverse steps can be undone to restore solved state
func (g *WaterBottleGame) validateReverseSteps() error {
	// Save current state
	currentState := g.copyGameState()
	
	// Apply reverse steps in reverse order (forward direction)
	fmt.Printf("🔍 验证 %d 步逆向操作的可还原性...\n", len(g.reverseSteps))
	
	for i := len(g.reverseSteps) - 1; i >= 0; i-- {
		move := g.reverseSteps[i]
		
		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败", 
				len(g.reverseSteps)-i, move.To, move.From)
		}
		
		if (i+1)%50 == 0 {
			fmt.Printf("   验证进度: %d/%d\n", len(g.reverseSteps)-i, len(g.reverseSteps))
		}
	}
	
	// Check if we're back to solved state
	if !g.IsWon() {
		return fmt.Errorf("还原后的状态不是完美解状态")
	}
	
	// Restore the generated initial state
	g.restoreGameState(currentState)
	return nil
}

// copyGameState creates a deep copy of the current game state
func (g *WaterBottleGame) copyGameState() [][]Color {
	state := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		state[i] = make([]Color, len(bottle))
		copy(state[i], bottle)
	}
	return state
}

// restoreGameState restores the game to a previous state
func (g *WaterBottleGame) restoreGameState(state [][]Color) {
	for i, bottleState := range state {
		g.bottles[i] = make(Bottle, len(bottleState))
		copy(g.bottles[i], bottleState)
	}
	
	// Recalculate empty count
	g.emptyCount = 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			g.emptyCount++
		}
	}
}

// GetReverseSteps returns the recorded reverse steps for analysis  
func (g *WaterBottleGame) GetReverseSteps() []Move {
	return g.reverseSteps
}

// getColorName returns the Chinese name of a color for display
func getColorName(color Color) string {
	names := []string{"红", "蓝", "绿", "黄", "橙", "紫", "棕", "黑", "白", "粉"}
	if int(color) < len(names) {
		return names[color]
	}
	return fmt.Sprintf("色%d", color)
}

// generatePerfectDistribution creates a state with perfect color distribution
func (g *WaterBottleGame) generatePerfectDistribution() error {
	totalWater := (g.N - g.J) * g.M

	// Calculate exact distribution
	baseAmount := totalWater / g.K / g.M * g.M // Each color gets multiple of M
	remainder := totalWater - baseAmount*g.K

	// Create all water units with perfect distribution
	var allWater []Color
	for colorID := 0; colorID < g.K; colorID++ {
		amount := baseAmount
		if colorID < remainder/g.M {
			amount += g.M // Distribute remainder evenly
		}
		for i := 0; i < amount; i++ {
			allWater = append(allWater, Color(colorID))
		}
	}

	// Shuffle the water thoroughly
	for i := 0; i < 10; i++ { // Multiple shuffles for better randomness
		rand.Shuffle(len(allWater), func(i, j int) {
			allWater[i], allWater[j] = allWater[j], allWater[i]
		})
	}

	// Distribute to bottles
	nonEmptyBottles := g.N - g.J
	waterIndex := 0

	for i := 0; i < nonEmptyBottles && waterIndex < len(allWater); i++ {
		// Random bottle size, but ensure we can distribute all water
		remainingBottles := nonEmptyBottles - i
		remainingWater := len(allWater) - waterIndex
		minSize := max(1, (remainingWater-remainingBottles*g.M+remainingBottles)/remainingBottles)
		maxSize := min(g.M, remainingWater-remainingBottles+1)

		if maxSize < minSize {
			maxSize = minSize
		}

		bottleSize := minSize
		if maxSize > minSize {
			bottleSize = rand.Intn(maxSize-minSize+1) + minSize
		}

		for j := 0; j < bottleSize && waterIndex < len(allWater); j++ {
			g.bottles[i] = append(g.bottles[i], allWater[waterIndex])
			waterIndex++
		}
	}

	// Distribute any remaining water
	for waterIndex < len(allWater) {
		bottleIdx := rand.Intn(nonEmptyBottles)
		if len(g.bottles[bottleIdx]) < g.M {
			g.bottles[bottleIdx] = append(g.bottles[bottleIdx], allWater[waterIndex])
			waterIndex++
		}
	}

	g.emptyCount = g.J
	return nil
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// tryGenerateState attempts to generate one random initial state
func (g *WaterBottleGame) tryGenerateState() error {
	totalWater := (g.N - g.J) * g.M

	// Create color distribution ensuring each color total is divisible by M
	colorCounts := make([]int, g.K)
	remaining := totalWater

	// Distribute water among colors, ensuring each can fill complete bottles
	for i := 0; i < g.K-1; i++ {
		maxPossible := remaining - (g.K-i-1)*g.M // Reserve at least M for each remaining color
		if maxPossible < g.M {
			maxPossible = g.M
		}
		// Random multiple of M
		bottles := rand.Intn(maxPossible/g.M) + 1
		colorCounts[i] = bottles * g.M
		remaining -= colorCounts[i]
	}
	colorCounts[g.K-1] = remaining

	if colorCounts[g.K-1] <= 0 || colorCounts[g.K-1]%g.M != 0 {
		return fmt.Errorf("invalid color distribution")
	}

	// Create all water units
	var allWater []Color
	for colorID := 0; colorID < g.K; colorID++ {
		for i := 0; i < colorCounts[colorID]; i++ {
			allWater = append(allWater, Color(colorID))
		}
	}

	// Shuffle the water
	rand.Shuffle(len(allWater), func(i, j int) {
		allWater[i], allWater[j] = allWater[j], allWater[i]
	})

	// Distribute water to non-empty bottles
	nonEmptyBottles := g.N - g.J
	waterIndex := 0

	for i := 0; i < nonEmptyBottles; i++ {
		bottleSize := rand.Intn(g.M) + 1 // Random size from 1 to M
		if waterIndex+bottleSize > len(allWater) {
			bottleSize = len(allWater) - waterIndex
		}

		for j := 0; j < bottleSize && waterIndex < len(allWater); j++ {
			g.bottles[i] = append(g.bottles[i], allWater[waterIndex])
			waterIndex++
		}
	}

	// Distribute remaining water
	for waterIndex < len(allWater) {
		bottleIdx := rand.Intn(nonEmptyBottles)
		if len(g.bottles[bottleIdx]) < g.M {
			g.bottles[bottleIdx] = append(g.bottles[bottleIdx], allWater[waterIndex])
			waterIndex++
		}
	}

	g.emptyCount = g.J
	return nil
}

// resetBottles clears all bottles for a new generation attempt
func (g *WaterBottleGame) resetBottles() {
	for i := range g.bottles {
		g.bottles[i] = g.bottles[i][:0]
	}
	g.emptyCount = g.J
}

// Pour attempts to pour water from one bottle to another
// Returns: success, amount moved
func (g *WaterBottleGame) Pour(fromBottle, toBottle int) (bool, int) {
	if fromBottle < 0 || fromBottle >= g.N || toBottle < 0 || toBottle >= g.N {
		return false, 0 // Invalid bottle indices
	}

	if fromBottle == toBottle {
		return false, 0 // Cannot pour to same bottle
	}

	from := &g.bottles[fromBottle]
	to := &g.bottles[toBottle]

	if len(*from) == 0 {
		return false, 0 // Cannot pour from empty bottle
	}

	if len(*to) >= g.M {
		return false, 0 // Target bottle is full
	}

	// Get the top color from source bottle
	topColor := (*from)[len(*from)-1]

	// Check if we can pour to target bottle
	if len(*to) > 0 && (*to)[len(*to)-1] != topColor {
		return false, 0 // Top colors don't match
	}

	// Count how many consecutive top colors we can pour
	fromIndex := len(*from) - 1
	for fromIndex >= 0 && (*from)[fromIndex] == topColor {
		fromIndex--
	}
	fromIndex++ // Now fromIndex points to the first occurrence of topColor from top

	availableAmount := len(*from) - fromIndex
	targetSpace := g.M - len(*to)
	pourAmount := min(availableAmount, targetSpace)

	if pourAmount <= 0 {
		return false, 0
	}

	// Perform the pour
	for i := 0; i < pourAmount; i++ {
		*to = append(*to, topColor)
	}
	*from = (*from)[:len(*from)-pourAmount]

	// Update empty bottle count
	wasFromEmpty := len(*from) == pourAmount
	wasToEmpty := len(*to) == pourAmount

	if wasFromEmpty && !wasToEmpty {
		g.emptyCount++
	} else if !wasFromEmpty && wasToEmpty {
		g.emptyCount--
	}

	return true, pourAmount
}

// IsWon checks if the game is won
func (g *WaterBottleGame) IsWon() bool {
	nonEmptyBottles := 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			continue
		}

		// Check if bottle is full and single-colored
		if len(bottle) != g.M {
			return false
		}

		color := bottle[0]
		for _, c := range bottle {
			if c != color {
				return false
			}
		}
		nonEmptyBottles++
	}

	return nonEmptyBottles == g.N-g.J
}

// GetState returns the current game state for display
func (g *WaterBottleGame) GetState() [][]Color {
	result := make([][]Color, g.N)
	for i, bottle := range g.bottles {
		result[i] = make([]Color, len(bottle))
		copy(result[i], bottle)
	}
	return result
}

// PrintState prints the current game state
func (g *WaterBottleGame) PrintState() {
	colorEmojis := []string{"🔴", "🔵", "🟢", "🟡", "🟠", "🟣", "🟤", "⚫", "⚪", "🔸"}

	fmt.Printf("\n🎮 当前游戏状态 (总瓶数:%d, 容量:%d, 空瓶:%d, 颜色数:%d):\n", g.N, g.M, g.J, g.K)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, bottle := range g.bottles {
		fmt.Printf("%d号瓶: ", i)
		if len(bottle) == 0 {
			fmt.Print("[空瓶子]")
		} else {
			fmt.Print("[")
			for j, color := range bottle {
				if j > 0 {
					fmt.Print(" ")
				}
				if int(color) < len(colorEmojis) {
					fmt.Printf("%s", colorEmojis[color])
				} else {
					fmt.Printf("%d", color)
				}
			}
			fmt.Print("]")
		}

		// 显示容量条
		filled := len(bottle)
		empty := g.M - filled

		// 防止负数导致panic
		if empty < 0 {
			empty = 0
			fmt.Printf(" ⚠️OVERFLOW⚠️ ")
		}

		fmt.Printf(" %s", strings.Repeat("█", min(filled, g.M)))
		fmt.Printf("%s", strings.Repeat("░", empty))
		fmt.Printf(" (%d/%d)", filled, g.M)

		// 检查是否是完成的瓶子（满瓶且单色）
		if len(bottle) == g.M && g.isSingleColor(bottle) {
			fmt.Print(" ✅完成")
		}
		fmt.Println()
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📊 空瓶子数量: %d\n", g.emptyCount)
	if g.IsWon() {
		fmt.Println("🎉 游戏胜利！所有瓶子都完成了！🎉")
	} else {
		fmt.Println("🎯 继续加油！目标：让每个瓶子都装满单一颜色")
	}
	fmt.Println()
}

// Helper function to check if a bottle contains only one color
func (g *WaterBottleGame) isSingleColor(bottle Bottle) bool {
	if len(bottle) == 0 {
		return true
	}
	firstColor := bottle[0]
	for _, color := range bottle {
		if color != firstColor {
			return false
		}
	}
	return true
}

// validateReverseSteps verifies that all reverse steps can be undone to restore solved state
func (g *WaterBottleGame) validateReverseSteps() error {
	// Save current state
	currentState := g.copyGameState()
	
	// Apply reverse steps in reverse order (forward direction)
	fmt.Printf("🔍 验证 %d 步逆向操作的可还原性...\n", len(g.reverseSteps))
	
	for i := len(g.reverseSteps) - 1; i >= 0; i-- {
		move := g.reverseSteps[i]
		
		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败", 
				len(g.reverseSteps)-i, move.To, move.From)
		}
		
		if (i+1)%50 == 0 {
			fmt.Printf("   验证进度: %d/%d\n", len(g.reverseSteps)-i, len(g.reverseSteps))
		}
	}
	
	// Check if we're back to solved state
	if !g.IsWon() {
		return fmt.Errorf("还原后的状态不是完美解状态")
	}
	
	// Restore the generated initial state
	g.restoreGameState(currentState)
	return nil
}

// copyGameState creates a deep copy of the current game state
func (g *WaterBottleGame) copyGameState() [][]Color {
	state := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		state[i] = make([]Color, len(bottle))
		copy(state[i], bottle)
	}
	return state
}

// restoreGameState restores the game to a previous state
func (g *WaterBottleGame) restoreGameState(state [][]Color) {
	for i, bottleState := range state {
		g.bottles[i] = make(Bottle, len(bottleState))
		copy(g.bottles[i], bottleState)
	}
	
	// Recalculate empty count
	g.emptyCount = 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			g.emptyCount++
		}
	}
}

// GetReverseSteps returns the recorded reverse steps for analysis  
func (g *WaterBottleGame) GetReverseSteps() []Move {
	return g.reverseSteps
}

// getColorName returns the Chinese name of a color for display
func getColorName(color Color) string {
	names := []string{"红", "蓝", "绿", "黄", "橙", "紫", "棕", "黑", "白", "粉"}
	if int(color) < len(names) {
		return names[color]
	}
	return fmt.Sprintf("色%d", color)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isSolvable uses multiple checks to verify if the current state is solvable
func (g *WaterBottleGame) isSolvable() bool {
	// 1. Basic mathematical constraints
	if !g.checkBasicConstraints() {
		return false
	}

	// 2. Advanced solvability check using BFS search
	return g.checkSolvabilityWithSearch()
}

// checkBasicConstraints verifies basic mathematical requirements
func (g *WaterBottleGame) checkBasicConstraints() bool {
	// Count each color
	colorCounts := make(map[Color]int)
	for _, bottle := range g.bottles {
		for _, color := range bottle {
			colorCounts[color]++
		}
	}

	// Check if each color count is divisible by M
	for _, count := range colorCounts {
		if count%g.M != 0 {
			return false
		}
	}

	// Check total water amount
	totalWater := 0
	for _, bottle := range g.bottles {
		totalWater += len(bottle)
	}

	expectedWater := (g.N - g.J) * g.M
	if totalWater != expectedWater {
		return false
	}

	// Must have at least one empty bottle for operations
	if g.emptyCount < 1 {
		return false
	}

	return true
}

// validateReverseSteps verifies that all reverse steps can be undone to restore solved state
func (g *WaterBottleGame) validateReverseSteps() error {
	// Save current state
	currentState := g.copyGameState()
	
	// Apply reverse steps in reverse order (forward direction)
	fmt.Printf("🔍 验证 %d 步逆向操作的可还原性...\n", len(g.reverseSteps))
	
	for i := len(g.reverseSteps) - 1; i >= 0; i-- {
		move := g.reverseSteps[i]
		
		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败", 
				len(g.reverseSteps)-i, move.To, move.From)
		}
		
		if (i+1)%50 == 0 {
			fmt.Printf("   验证进度: %d/%d\n", len(g.reverseSteps)-i, len(g.reverseSteps))
		}
	}
	
	// Check if we're back to solved state
	if !g.IsWon() {
		return fmt.Errorf("还原后的状态不是完美解状态")
	}
	
	// Restore the generated initial state
	g.restoreGameState(currentState)
	return nil
}

// copyGameState creates a deep copy of the current game state
func (g *WaterBottleGame) copyGameState() [][]Color {
	state := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		state[i] = make([]Color, len(bottle))
		copy(state[i], bottle)
	}
	return state
}

// restoreGameState restores the game to a previous state
func (g *WaterBottleGame) restoreGameState(state [][]Color) {
	for i, bottleState := range state {
		g.bottles[i] = make(Bottle, len(bottleState))
		copy(g.bottles[i], bottleState)
	}
	
	// Recalculate empty count
	g.emptyCount = 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			g.emptyCount++
		}
	}
}

// GetReverseSteps returns the recorded reverse steps for analysis  
func (g *WaterBottleGame) GetReverseSteps() []Move {
	return g.reverseSteps
}

// getColorName returns the Chinese name of a color for display
func getColorName(color Color) string {
	names := []string{"红", "蓝", "绿", "黄", "橙", "紫", "棕", "黑", "白", "粉"}
	if int(color) < len(names) {
		return names[color]
	}
	return fmt.Sprintf("色%d", color)
}

// checkSolvabilityWithSearch uses a limited BFS to check if solution exists
func (g *WaterBottleGame) checkSolvabilityWithSearch() bool {
	// For performance, we only do a limited search
	maxDepth := 10
	maxStates := 1000

	type GameState struct {
		bottles [][]Color
		depth   int
	}

	// Create initial state
	initialState := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		initialState[i] = make([]Color, len(bottle))
		copy(initialState[i], bottle)
	}

	queue := []GameState{{bottles: initialState, depth: 0}}
	visited := make(map[string]bool)
	statesChecked := 0

	for len(queue) > 0 && statesChecked < maxStates {
		current := queue[0]
		queue = queue[1:]
		statesChecked++

		// Check if this state is solved
		if g.isStateWon(current.bottles) {
			return true
		}

		// Don't go too deep
		if current.depth >= maxDepth {
			continue
		}

		// Generate state signature for visited check
		signature := g.getStateSignature(current.bottles)
		if visited[signature] {
			continue
		}
		visited[signature] = true

		// Try all possible moves
		for from := 0; from < len(current.bottles); from++ {
			for to := 0; to < len(current.bottles); to++ {
				if from == to {
					continue
				}

				newBottles := g.copyState(current.bottles)
				if g.tryPourInState(newBottles, from, to) {
					queue = append(queue, GameState{
						bottles: newBottles,
						depth:   current.depth + 1,
					})
				}
			}
		}
	}

	// If we checked many states and found no solution, it might be unsolvable
	// But for performance, we assume it's solvable if basic constraints pass
	return statesChecked < maxStates
}

// Helper functions for the search algorithm
func (g *WaterBottleGame) isStateWon(bottles [][]Color) bool {
	nonEmptyCount := 0
	for _, bottle := range bottles {
		if len(bottle) == 0 {
			continue
		}

		if len(bottle) != g.M {
			return false
		}

		// Check if single color
		if len(bottle) > 0 {
			firstColor := bottle[0]
			for _, color := range bottle {
				if color != firstColor {
					return false
				}
			}
		}
		nonEmptyCount++
	}

	return nonEmptyCount == g.N-g.J
}

func (g *WaterBottleGame) getStateSignature(bottles [][]Color) string {
	var parts []string
	for _, bottle := range bottles {
		var bottleStr strings.Builder
		for _, color := range bottle {
			bottleStr.WriteString(fmt.Sprintf("%d,", color))
		}
		parts = append(parts, bottleStr.String())
	}
	return strings.Join(parts, "|")
}

func (g *WaterBottleGame) copyState(bottles [][]Color) [][]Color {
	newBottles := make([][]Color, len(bottles))
	for i, bottle := range bottles {
		newBottles[i] = make([]Color, len(bottle))
		copy(newBottles[i], bottle)
	}
	return newBottles
}

func (g *WaterBottleGame) tryPourInState(bottles [][]Color, from, to int) bool {
	fromBottle := bottles[from]
	toBottle := bottles[to]

	if len(fromBottle) == 0 {
		return false // Cannot pour from empty bottle
	}

	if len(toBottle) >= g.M {
		return false // Target bottle is full
	}

	// Get the top color from source bottle
	topColor := fromBottle[len(fromBottle)-1]

	// Check if we can pour to target bottle
	if len(toBottle) > 0 && toBottle[len(toBottle)-1] != topColor {
		return false // Top colors don't match
	}

	// Count how many consecutive top colors we can pour
	fromIndex := len(fromBottle) - 1
	for fromIndex >= 0 && fromBottle[fromIndex] == topColor {
		fromIndex--
	}
	fromIndex++ // Now fromIndex points to the first occurrence of topColor from top

	availableAmount := len(fromBottle) - fromIndex
	targetSpace := g.M - len(toBottle)
	pourAmount := min(availableAmount, targetSpace)

	if pourAmount <= 0 {
		return false
	}

	// Perform the pour
	for i := 0; i < pourAmount; i++ {
		bottles[to] = append(bottles[to], topColor)
	}
	bottles[from] = bottles[from][:len(bottles[from])-pourAmount]

	return true
}

// validateReverseSteps verifies that all reverse steps can be undone to restore solved state
func (g *WaterBottleGame) validateReverseSteps() error {
	// Save current state
	currentState := g.copyGameState()
	
	// Apply reverse steps in reverse order (forward direction)
	fmt.Printf("🔍 验证 %d 步逆向操作的可还原性...\n", len(g.reverseSteps))
	
	for i := len(g.reverseSteps) - 1; i >= 0; i-- {
		move := g.reverseSteps[i]
		
		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败", 
				len(g.reverseSteps)-i, move.To, move.From)
		}
		
		if (i+1)%50 == 0 {
			fmt.Printf("   验证进度: %d/%d\n", len(g.reverseSteps)-i, len(g.reverseSteps))
		}
	}
	
	// Check if we're back to solved state
	if !g.IsWon() {
		return fmt.Errorf("还原后的状态不是完美解状态")
	}
	
	// Restore the generated initial state
	g.restoreGameState(currentState)
	return nil
}

// copyGameState creates a deep copy of the current game state
func (g *WaterBottleGame) copyGameState() [][]Color {
	state := make([][]Color, len(g.bottles))
	for i, bottle := range g.bottles {
		state[i] = make([]Color, len(bottle))
		copy(state[i], bottle)
	}
	return state
}

// restoreGameState restores the game to a previous state
func (g *WaterBottleGame) restoreGameState(state [][]Color) {
	for i, bottleState := range state {
		g.bottles[i] = make(Bottle, len(bottleState))
		copy(g.bottles[i], bottleState)
	}
	
	// Recalculate empty count
	g.emptyCount = 0
	for _, bottle := range g.bottles {
		if len(bottle) == 0 {
			g.emptyCount++
		}
	}
}

// GetReverseSteps returns the recorded reverse steps for analysis  
func (g *WaterBottleGame) GetReverseSteps() []Move {
	return g.reverseSteps
}

// getColorName returns the Chinese name of a color for display
func getColorName(color Color) string {
	names := []string{"红", "蓝", "绿", "黄", "橙", "紫", "棕", "黑", "白", "粉"}
	if int(color) < len(names) {
		return names[color]
	}
	return fmt.Sprintf("色%d", color)
}
