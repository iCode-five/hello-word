package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Demo function to show how to use the water bottle game
func runWaterBottleDemo() {
	fmt.Println("🎮 欢迎来到水瓶分色游戏！")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Get custom parameters from user
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("📝 请设置游戏参数：")
	fmt.Println()

	// Get N (total bottles)
	var N, M, J, K int
	var err error

	for {
		fmt.Print("🍶 请输入总瓶子数量 N (建议 4-10，最大20): ")
		if !scanner.Scan() {
			return
		}
		N, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || N < 3 || N > 20 {
			fmt.Println("❌ 请输入 3-20 之间的数字")
			continue
		}
		break
	}

	// Get M (bottle capacity)
	for {
		fmt.Print("📏 请输入每个瓶子的容量 M (建议 3-6): ")
		if !scanner.Scan() {
			return
		}
		M, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || M < 2 || M > 10 {
			fmt.Println("❌ 请输入 2-10 之间的数字")
			continue
		}
		break
	}

	// Get J (empty bottles)
	for {
		fmt.Printf("🫗 请输入空瓶子数量 J (建议 1-%d): ", N-2)
		if !scanner.Scan() {
			return
		}
		J, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || J < 1 || J >= N {
			fmt.Printf("❌ 请输入 1-%d 之间的数字\n", N-1)
			continue
		}
		break
	}

	// Get K (number of colors)
	totalWater := (N - J) * M
	maxColors := totalWater / M               // Maximum possible colors (each needs at least M units)
	recommendedMaxColors := maxColors * 2 / 3 // Leave some room for randomness
	if maxColors > 8 {
		maxColors = 8 // Limit for visual clarity
	}
	if recommendedMaxColors < 2 {
		recommendedMaxColors = 2
	}

	for {
		fmt.Printf("🎨 请输入颜色种类数 K (建议 2-%d, 最大%d): ", recommendedMaxColors, maxColors)
		if !scanner.Scan() {
			return
		}
		K, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || K < 2 {
			fmt.Println("❌ 颜色数至少需要2种")
			continue
		}
		if K > maxColors {
			fmt.Printf("❌ 颜色数太多！总水量%d，每种颜色至少需要%d单位，最多只能有%d种颜色\n", totalWater, M, maxColors)
			continue
		}
		if K > recommendedMaxColors {
			fmt.Printf("⚠️  颜色数较多，可能生成困难。建议不超过%d种。确定要使用%d种颜色吗？(y/n): ", recommendedMaxColors, K)
			if !scanner.Scan() {
				return
			}
			if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
				continue
			}
		}
		break
	}

	// Choose generation method
	var generationMethod string
	for {
		fmt.Println("🎲 请选择初始状态生成方式：")
		fmt.Println("  1. 逆向生成（保证有解，推荐）")
		fmt.Println("  2. 纯随机生成（可能无解，更有挑战性）")
		fmt.Print("请输入选择 (1/2): ")
		if !scanner.Scan() {
			return
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "1" {
			generationMethod = "reverse"
			break
		} else if choice == "2" {
			generationMethod = "random"
			break
		} else {
			fmt.Println("❌ 请输入 1 或 2")
			continue
		}
	}

	// Get reverse steps (difficulty) only for reverse generation
	var reverseSteps int
	if generationMethod == "reverse" {
		suggestedSteps := N * K * M / 4 // Suggested based on complexity
		if suggestedSteps < 10 {
			suggestedSteps = 10
		}

		for {
			fmt.Printf("🎯 请输入逆序步数（游戏难度）(建议 %d, 范围 5-1000): ", suggestedSteps)
			if !scanner.Scan() {
				return
			}
			reverseSteps, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || reverseSteps < 5 || reverseSteps > 1000 {
				fmt.Println("❌ 请输入 5-1000 之间的数字")
				continue
			}
			break
		}
	}

	fmt.Println()
	if generationMethod == "reverse" {
		fmt.Printf("✅ 游戏参数设置完成：%d个瓶子，每个容量%d，%d个空瓶，%d种颜色\n", N, M, J, K)
		fmt.Printf("🔄 使用逆向生成，%d步逆序\n", reverseSteps)
	} else {
		fmt.Printf("✅ 游戏参数设置完成：%d个瓶子，每个容量%d，%d个空瓶，%d种颜色\n", N, M, J, K)
		fmt.Println("🎲 使用纯随机生成")
	}
	fmt.Println("正在生成游戏初始状态...")

	// Create game with user parameters
	game1, err := NewWaterBottleGame(N, M, J, K)
	if err != nil {
		fmt.Printf("❌ 创建游戏失败: %v\n", err)
		return
	}

	// Generate initial state based on chosen method
	if generationMethod == "reverse" {
		err = game1.generateInitialStateWithSteps(reverseSteps)
	} else {
		err = game1.generateRandomState()
	}

	if err != nil {
		fmt.Printf("❌ 生成初始状态失败: %v\n", err)
		return
	}

	fmt.Println("🎯 初始状态生成完成！")
	game1.PrintState()

	// Show initial move status
	if !game1.IsWon() {
		game1.PrintMoveStatus()
	}

	// Interactive mode
	fmt.Println("\n=== 🎮 开始游戏！===")
	fmt.Println("游戏目标：通过倒水让每个瓶子都装满单一颜色的水")
	fmt.Println("数字代表颜色：0=红色 🔴, 1=蓝色 🔵, 2=绿色 🟢, 3=黄色 🟡")
	fmt.Println()
	fmt.Println("📋 可用命令：")
	fmt.Println("  倒水 <源瓶子> <目标瓶子>     - 例如：倒水 0 3 （从0号瓶倒到3号瓶）")
	fmt.Println("  状态                       - 查看当前游戏状态和可能移动")
	fmt.Println("  检查                       - 单独检查可能的移动")
	fmt.Println("  新游戏 <瓶数> <容量> <空瓶数> <颜色数> [生成方式] - 创建新游戏")
	fmt.Println("    生成方式: random(随机) 或 reverse(逆向，默认)")
	fmt.Println("  退出                       - 结束游戏")
	fmt.Println()
	fmt.Println("💡 提示：只能倒到空瓶或者顶层颜色相同的瓶子里")
	fmt.Println()

	currentGame := game1

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "quit", "exit", "q", "退出":
			fmt.Println("👋 感谢游戏！再见！")
			return

		case "state", "s", "状态":
			currentGame.PrintState()
			if !currentGame.IsWon() {
				currentGame.PrintMoveStatus()
			}

		case "check", "c", "检查", "移动":
			currentGame.PrintMoveStatus()

		case "pour", "p", "倒水":
			if len(parts) != 3 {
				fmt.Println("❌ 用法：倒水 <源瓶子号> <目标瓶子号>")
				fmt.Println("   例如：倒水 0 3")
				continue
			}

			from, err1 := strconv.Atoi(parts[1])
			to, err2 := strconv.Atoi(parts[2])

			if err1 != nil || err2 != nil {
				fmt.Println("❌ 瓶子编号必须是数字")
				continue
			}

			success, moved := currentGame.Pour(from, to)
			if success {
				fmt.Printf("✅ 成功从 %d 号瓶倒了 %d 单位水到 %d 号瓶\n", from, moved, to)
				currentGame.PrintState()

				if currentGame.IsWon() {
					fmt.Println("🎉🎉🎉 恭喜！你赢了！所有瓶子都是单色满瓶！🎉🎉🎉")
				} else {
					// Check for possible moves after each successful move
					currentGame.PrintMoveStatus()
				}
			} else {
				fmt.Printf("❌ 无法从 %d 号瓶倒水到 %d 号瓶\n", from, to)
				fmt.Println("💡 检查：源瓶是否有水？目标瓶是否满了？顶层颜色是否匹配？")
			}

		case "new", "n", "新游戏":
			if len(parts) < 5 || len(parts) > 6 {
				fmt.Println("❌ 用法：新游戏 <瓶子数> <容量> <空瓶数> <颜色数> [生成方式]")
				fmt.Println("   例如：新游戏 5 4 2 3        （默认逆向生成）")
				fmt.Println("   例如：新游戏 5 4 2 3 random （纯随机生成）")
				fmt.Println("   例如：新游戏 5 4 2 3 reverse（逆向生成）")
				continue
			}

			N, err1 := strconv.Atoi(parts[1])
			M, err2 := strconv.Atoi(parts[2])
			J, err3 := strconv.Atoi(parts[3])
			K, err4 := strconv.Atoi(parts[4])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				fmt.Println("❌ 所有参数必须是数字")
				continue
			}

			// Determine generation method
			genMethod := "reverse" // default
			if len(parts) == 6 {
				method := strings.ToLower(parts[5])
				if method == "random" || method == "r" {
					genMethod = "random"
				} else if method == "reverse" || method == "rev" {
					genMethod = "reverse"
				} else {
					fmt.Println("❌ 生成方式必须是 'random' 或 'reverse'")
					continue
				}
			}

			newGame, err := NewWaterBottleGame(N, M, J, K)
			if err != nil {
				fmt.Printf("❌ 创建游戏失败: %v\n", err)
				continue
			}

			// Generate based on method
			if genMethod == "random" {
				fmt.Println("🎲 使用纯随机生成...")
				err = newGame.generateRandomState()
			} else {
				fmt.Println("🔄 使用逆向生成...")
				err = newGame.generateInitialState()
			}

			if err != nil {
				fmt.Printf("❌ 生成初始状态失败: %v\n", err)
				continue
			}

			currentGame = newGame
			fmt.Printf("✅ 新游戏创建成功！（%s生成）\n",
				map[string]string{"random": "随机", "reverse": "逆向"}[genMethod])
			currentGame.PrintState()

		case "help", "h", "帮助":
			fmt.Println("📋 可用命令：")
			fmt.Println("  倒水 <源瓶子> <目标瓶子>     - 例如：倒水 0 3")
			fmt.Println("  状态                       - 查看当前游戏状态和可能移动")
			fmt.Println("  检查                       - 单独检查可能的移动")
			fmt.Println("  新游戏 <瓶数> <容量> <空瓶数> <颜色数> [生成方式] - 创建新游戏")
			fmt.Println("    生成方式: random(随机) 或 reverse(逆向，默认)")
			fmt.Println("    例如：新游戏 5 4 2 3 random")
			fmt.Println("  帮助                       - 显示此帮助信息")
			fmt.Println("  退出                       - 结束游戏")

		default:
			fmt.Printf("❓ 未知命令：%s\n", command)
			fmt.Println("💡 输入 '帮助' 查看可用命令")
		}
	}
}

// Example of programmatic game solving (basic strategy)
func demonstrateBasicSolver() {
	fmt.Println("\n=== Basic Solver Demonstration ===")

	game, err := NewWaterBottleGame(4, 3, 1, 2)
	if err != nil {
		fmt.Printf("Error creating game: %v\n", err)
		return
	}

	err = game.generateInitialState()
	if err != nil {
		fmt.Printf("Error generating initial state: %v\n", err)
		return
	}

	fmt.Println("Initial state:")
	game.PrintState()

	// Try a few strategic moves (scale with game complexity)
	moves := 0
	maxMoves := max(20, game.N*5)

	for !game.IsWon() && moves < maxMoves {
		moved := false

		// Strategy: Try to consolidate same colors
		for from := 0; from < game.N && !moved; from++ {
			for to := 0; to < game.N && !moved; to++ {
				if from != to {
					success, amount := game.Pour(from, to)
					if success {
						fmt.Printf("Move %d: Poured %d units from bottle %d to bottle %d\n",
							moves+1, amount, from, to)
						game.PrintState()
						moved = true
						moves++
					}
				}
			}
		}

		if !moved {
			fmt.Println("No more moves possible")
			break
		}
	}

	if game.IsWon() {
		fmt.Printf("🎉 Solved in %d moves! 🎉\n", moves)
	} else {
		fmt.Printf("Could not solve within %d moves\n", maxMoves)
	}
}

// Main function to run the interactive demo
func main() {
	// Run the interactive demo
	runWaterBottleDemo()

	// Or run the basic solver demonstration
	// demonstrateBasicSolver()
}
