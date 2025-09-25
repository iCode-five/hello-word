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
	From   int   // Source bottle index
	To     int   // Target bottle index
	Amount int   // Amount of water moved
	Color  Color // Color of water moved
}

// WaterBottleGame represents the game state
type WaterBottleGame struct {
	bottles      []Bottle // All bottles in the game
	N            int      // Total number of bottles
	M            int      // Capacity of each bottle
	J            int      // Number of empty bottles
	K            int      // Number of different colors
	emptyCount   int      // Current number of empty bottles
	reverseSteps []Move   // Record of reverse operations for validation
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

// generateRandomState creates a completely random initial state (may not be solvable)
func (g *WaterBottleGame) generateRandomState() error {
	rand.Seed(time.Now().UnixNano())

	// Check if parameters are reasonable
	totalWater := (g.N - g.J) * g.M
	maxPossibleColors := totalWater / g.M
	if g.K > maxPossibleColors {
		return fmt.Errorf("参数不合理：总水量%d，每种颜色至少需要%d单位，最多只能有%d种颜色，但要求%d种",
			totalWater, g.M, maxPossibleColors, g.K)
	}

	fmt.Println("🎲 正在进行纯随机生成...")

	// Create a pool of all water units with correct color distribution
	waterPool := g.createColorPool()

	// Shuffle the water pool randomly
	g.shuffleColorPool(waterPool)

	// Distribute water randomly into bottles
	return g.distributeWaterRandomly(waterPool)
}

// createColorPool creates a pool of water units with balanced color distribution
func (g *WaterBottleGame) createColorPool() []Color {
	totalWater := (g.N - g.J) * g.M

	// Calculate how many units each color should have
	baseUnitsPerColor := totalWater / g.K
	extraUnits := totalWater % g.K

	waterPool := make([]Color, 0, totalWater)

	for colorID := 0; colorID < g.K; colorID++ {
		unitsForThisColor := baseUnitsPerColor
		if colorID < extraUnits {
			unitsForThisColor++
		}

		// Add this color to the pool
		for i := 0; i < unitsForThisColor; i++ {
			waterPool = append(waterPool, Color(colorID))
		}
	}

	fmt.Printf("   💧 创建水池：总共%d单位水，%d种颜色\n", len(waterPool), g.K)

	// Print color distribution
	colorCounts := make(map[Color]int)
	for _, color := range waterPool {
		colorCounts[color]++
	}

	fmt.Print("   🎨 颜色分布：")
	for colorID := 0; colorID < g.K; colorID++ {
		fmt.Printf("%s×%d ", getColorName(Color(colorID)), colorCounts[Color(colorID)])
	}
	fmt.Println()

	return waterPool
}

// shuffleColorPool randomly shuffles the water pool using Fisher-Yates algorithm
func (g *WaterBottleGame) shuffleColorPool(pool []Color) {
	fmt.Println("   🔀 随机打乱水池...")

	for i := len(pool) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		pool[i], pool[j] = pool[j], pool[i]
	}
}

// distributeWaterRandomly distributes shuffled water into bottles randomly
func (g *WaterBottleGame) distributeWaterRandomly(waterPool []Color) error {
	fmt.Println("   🍶 随机分配水到瓶子...")

	// Clear all bottles first
	for i := range g.bottles {
		g.bottles[i] = make(Bottle, 0, g.M)
	}
	g.emptyCount = g.J

	// Randomly choose which bottles to fill (leaving J empty)
	bottlesToFill := make([]int, 0, g.N-g.J)
	for i := 0; i < g.N-g.J; i++ {
		bottlesToFill = append(bottlesToFill, i)
	}

	// Shuffle bottle order
	for i := len(bottlesToFill) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		bottlesToFill[i], bottlesToFill[j] = bottlesToFill[j], bottlesToFill[i]
	}

	waterIndex := 0

	// Fill bottles completely randomly
	for _, bottleIdx := range bottlesToFill {
		// Fill this bottle to capacity
		for unit := 0; unit < g.M && waterIndex < len(waterPool); unit++ {
			g.bottles[bottleIdx] = append(g.bottles[bottleIdx], waterPool[waterIndex])
			waterIndex++
		}
	}

	// Verify we used all water
	if waterIndex != len(waterPool) {
		return fmt.Errorf("水分配错误：应该分配%d单位，实际分配%d单位", len(waterPool), waterIndex)
	}

	fmt.Printf("   ✅ 随机分配完成！填充了%d个瓶子，保留%d个空瓶\n", g.N-g.J, g.J)

	// Analyze the generated state
	g.analyzeRandomState()

	return nil
}

// analyzeRandomState analyzes the randomly generated state
func (g *WaterBottleGame) analyzeRandomState() {
	fmt.Println("   📊 随机状态分析：")

	mixedBottles := 0
	singleColorBottles := 0

	for i, bottle := range g.bottles {
		if len(bottle) == 0 {
			continue
		}

		if g.isSingleColor(bottle) {
			singleColorBottles++
			if len(bottle) == g.M {
				fmt.Printf("      瓶子%d：✅ 已完成（单色满瓶）\n", i)
			} else {
				fmt.Printf("      瓶子%d：🟡 单色但未满\n", i)
			}
		} else {
			mixedBottles++
			fmt.Printf("      瓶子%d：🔴 混色瓶\n", i)
		}
	}

	fmt.Printf("   📈 统计：%d个混色瓶，%d个单色瓶\n", mixedBottles, singleColorBottles)

	if g.IsWon() {
		fmt.Println("   🎉 幸运！随机生成了一个已完成的状态！")
	} else if mixedBottles == 0 {
		fmt.Println("   🟡 生成了全单色状态，但可能有未满的瓶子")
	} else {
		fmt.Println("   🎯 生成了混合状态，需要玩家解决")
	}
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
		maxAttempts := 100 // Increase attempts per step
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
			fmt.Printf("   ⏹️  逆向操作已达到极限，实际完成 %d 步有效逆向操作\n", actualSteps)
			fmt.Printf("   📊 总共尝试了 %d 次操作，成功率 %.1f%%\n", totalAttempts, float64(actualSteps)/float64(totalAttempts)*100)
			fmt.Printf("   ✅ 当前状态已足够复杂，继续正常流程...\n")
			break
		}
	}

	successRate := float64(actualSteps) / float64(totalAttempts) * 100
	fmt.Printf("🎯 逆向生成完成！\n")
	fmt.Printf("   - 目标步数: %d 步\n", reverseSteps)
	fmt.Printf("   - 实际完成: %d 步 ", actualSteps)
	if actualSteps < reverseSteps {
		fmt.Printf("(已达到复杂度极限)\n")
	} else {
		fmt.Printf("(完全达成目标)\n")
	}
	fmt.Printf("   - 总尝试次数: %d 次\n", totalAttempts)
	fmt.Printf("   - 成功率: %.1f%%\n", successRate)

	// Step 3: Validate that we can restore the original state using the recorded steps
	if actualSteps > 0 {
		if err := g.validateReverseSteps(); err != nil {
			return fmt.Errorf("逆向步骤验证失败: %v", err)
		}
		fmt.Println("✅ 逆向步骤验证成功！所有操作都可以还原")
	} else {
		fmt.Println("ℹ️  没有执行逆向操作，保持完美解状态")
	}

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

	// Try multiple random combinations (scale with bottle count)
	maxAttempts := min(50, g.N*5)
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

			// Perform the reverse move and test if it can be immediately reversed
			if g.performSimplePour(sourceIdx, targetIdx, pourAmount) {
				// Now try to reverse this operation using forward game rules
				canReverse, actualMoved := g.Pour(targetIdx, sourceIdx)

				if canReverse && actualMoved == pourAmount {
					// Check if we're back to the original state
					if g.statesEqual(stateBefore, g.copyGameState()) {
						// Perfect! This reverse operation is valid
						// Restore to the state after reverse operation (before the test)
						g.restoreGameState(stateBefore)
						g.performSimplePour(sourceIdx, targetIdx, pourAmount)

						// Record the move
						move := Move{
							From:   sourceIdx,
							To:     targetIdx,
							Amount: pourAmount,
							Color:  topColor,
						}
						g.reverseSteps = append(g.reverseSteps, move)
						return true
					}
				}

				// If we can't properly reverse this move, restore original state
				g.restoreGameState(stateBefore)
				continue // Try next target or next attempt
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
		stepNum := len(g.reverseSteps) - i

		// Apply the reverse of this move (from To back to From)
		success, _ := g.Pour(move.To, move.From)
		if !success {
			fmt.Printf("   ❌ 第%d步还原失败: 从%d号瓶到%d号瓶\n", stepNum, move.To, move.From)
			fmt.Printf("      原始逆向操作: 从%d号瓶倒%d单位%s色水到%d号瓶\n",
				move.From, move.Amount, getColorName(move.Color), move.To)
			return fmt.Errorf("步骤 %d 无法还原: 从瓶子%d到瓶子%d失败",
				stepNum, move.To, move.From)
		}

		// Show all successful restoration steps in simple format
		fmt.Printf("倒水 %d %d\n", move.To, move.From)

		if stepNum%50 == 0 {
			fmt.Printf("   📊 验证进度: %d/%d\n", stepNum, len(g.reverseSteps))
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

// CheckPossibleMoves checks if there are any possible moves and returns detailed information
func (g *WaterBottleGame) CheckPossibleMoves() (bool, int, []string) {
	possibleMoves := 0
	moveDescriptions := make([]string, 0)

	for from := 0; from < g.N; from++ {
		for to := 0; to < g.N; to++ {
			if from != to {
				// Save current state
				originalState := g.copyGameState()

				success, moved := g.Pour(from, to)
				if success {
					possibleMoves++
					// Create move description
					fromBottle := originalState[from]
					toBottle := originalState[to]

					var fromDesc, toDesc string
					if len(fromBottle) == 0 {
						fromDesc = "空瓶"
					} else {
						topColor := fromBottle[len(fromBottle)-1]
						fromDesc = fmt.Sprintf("顶层%s色", getColorName(topColor))
					}

					if len(toBottle) == 0 {
						toDesc = "空瓶"
					} else {
						topColor := toBottle[len(toBottle)-1]
						toDesc = fmt.Sprintf("顶层%s色", getColorName(topColor))
					}

					moveDesc := fmt.Sprintf("从%d号瓶(%s)倒%d单位到%d号瓶(%s)",
						from, fromDesc, moved, to, toDesc)
					moveDescriptions = append(moveDescriptions, moveDesc)
				}

				// Restore state
				g.restoreGameState(originalState)
			}
		}
	}

	return possibleMoves > 0, possibleMoves, moveDescriptions
}

// PrintMoveStatus prints the current move status
func (g *WaterBottleGame) PrintMoveStatus() {
	hasMoves, moveCount, moveDescriptions := g.CheckPossibleMoves()

	fmt.Printf("\n🔍 移动状态检查：\n")
	if !hasMoves {
		fmt.Println("🚨 没有可用的移动！")
		if g.IsWon() {
			fmt.Println("🎉 游戏已完成！")
		} else {
			fmt.Println("💀 游戏陷入死局！")
			g.analyzeDeadlock()
		}
	} else {
		fmt.Printf("✅ 共有 %d 种可能的移动：\n", moveCount)

		// Show first few moves as examples
		maxShow := min(5, len(moveDescriptions))
		for i := 0; i < maxShow; i++ {
			fmt.Printf("  • %s\n", moveDescriptions[i])
		}

		if len(moveDescriptions) > maxShow {
			fmt.Printf("  • ... 还有 %d 种其他移动\n", len(moveDescriptions)-maxShow)
		}
	}
	fmt.Println()
}

// analyzeDeadlock analyzes why the game is in deadlock
func (g *WaterBottleGame) analyzeDeadlock() {
	fmt.Println("📊 死局分析：")

	// Check empty bottles
	if g.emptyCount == 0 {
		fmt.Println("  ❌ 没有空瓶子可以倒水")
	} else {
		fmt.Printf("  ✅ 还有 %d 个空瓶子\n", g.emptyCount)
	}

	// Check top colors
	topColors := make(map[Color][]int) // color -> bottle indices
	for i, bottle := range g.bottles {
		if len(bottle) > 0 {
			topColor := bottle[len(bottle)-1]
			topColors[topColor] = append(topColors[topColor], i)
		}
	}

	fmt.Printf("  📈 顶层颜色分布：\n")
	allDifferent := true
	for color, bottles := range topColors {
		if len(bottles) > 1 {
			allDifferent = false
			fmt.Printf("    %s色：瓶子 %v（可以互相倒水）\n", getColorName(color), bottles)
		} else {
			fmt.Printf("    %s色：瓶子 %v（孤立）\n", getColorName(color), bottles)
		}
	}

	if allDifferent && g.emptyCount == 0 {
		fmt.Println("  🚨 死局原因：所有瓶子顶层颜色都不同，且没有空瓶")
	}
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

	// Adjust separator length based on bottle count
	separatorLength := min(80, max(50, g.N*8))
	fmt.Println(strings.Repeat("━", separatorLength))

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

	fmt.Println(strings.Repeat("━", separatorLength))
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

// statesEqual compares two game states for equality
func (g *WaterBottleGame) statesEqual(state1, state2 [][]Color) bool {
	if len(state1) != len(state2) {
		return false
	}

	for i := range state1 {
		if len(state1[i]) != len(state2[i]) {
			return false
		}

		for j := range state1[i] {
			if state1[i][j] != state2[i][j] {
				return false
			}
		}
	}

	return true
}
