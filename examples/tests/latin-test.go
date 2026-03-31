package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/corrreia/govfd"
	"github.com/corrreia/govfd/types"
)

type TestResult struct {
	TestNumber     int    `json:"test_number"`
	Category       string `json:"category"`
	InputText      string `json:"input_text"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description"`
	Worked         bool   `json:"worked"`
}

type CategoryStats struct {
	Total       int     `json:"total"`
	Passed      int     `json:"passed"`
	SuccessRate float64 `json:"success_rate"`
}

type TestSummary struct {
	TestSuite     string                   `json:"test_suite"`
	Description   string                   `json:"description"`
	Timestamp     string                   `json:"timestamp"`
	TotalTests    int                      `json:"total_tests"`
	Passed        int                      `json:"passed"`
	Failed        int                      `json:"failed"`
	Errors        int                      `json:"errors"`
	SuccessRate   float64                  `json:"success_rate"`
	CategoryStats map[string]CategoryStats `json:"category_stats"`
	TestResults   []TestResult             `json:"test_results"`
}

func main() {
	fmt.Println("🌍 === LATIN CHARACTER TEST SUITE ===")
	fmt.Println("🎯 Tests only Latin-based characters that work!")
	fmt.Println("✨ Portuguese, Spanish, French, German, Italian")
	fmt.Println()

	// Connect to VFD
	display, err := govfd.OpenModel("/dev/ttyUSB0", types.ModelEpsonDMD110)
	if err != nil {
		fmt.Printf("❌ Error opening display: %v\n", err)
		return
	}
	defer display.Close()

	// Clear display
	display.Clear()
	display.WriteText("GoVFD Latin Test")
	time.Sleep(2 * time.Second)

	// Define test cases - only Latin characters that should work
	tests := []struct {
		category    string
		text        string
		expected    string
		description string
	}{
		// Portuguese - CP860
		{"🇵🇹 Portuguese", "ação", "ação", "Portuguese with tilde and cedilla"},
		{"🇵🇹 Portuguese", "São Paulo", "São Paulo", "Portuguese proper noun"},
		{"🇵🇹 Portuguese", "coração", "coração", "Portuguese with tilde"},

		// Spanish - CP850
		{"🇪🇸 Spanish", "niño", "niño", "Spanish with tilde n"},
		{"🇪🇸 Spanish", "España", "España", "Spanish with tilde n"},

		// French - CP850
		{"🇫🇷 French", "café", "café", "French with acute accent"},
		{"🇫🇷 French", "résumé", "résumé", "French with acute accents"},
		{"🇫🇷 French", "naïve", "naïve", "French with diaeresis"},

		// German - CP850
		{"🇩🇪 German", "Müller", "Müller", "German umlaut u"},
		{"🇩🇪 German", "größe", "größe", "German umlaut o and eszett"},

		// Italian - CP850
		{"🇮🇹 Italian", "città", "città", "Italian with grave accent"},
		{"🇮🇹 Italian", "università", "università", "Italian with grave accents"},

		// Euro symbol - CP858
		{"💰 Euro", "€19.99", "€19.99", "Euro symbol with price"},
	}

	var results []TestResult
	categoryCount := make(map[string]int)
	categoryPassed := make(map[string]int)

	fmt.Printf("📺 Display: 20x2 VFD\n")
	fmt.Printf("🧪 Total tests: %d\n", len(tests))
	fmt.Printf("⏰ Each test waits for YOUR confirmation!\n\n")

	for i, test := range tests {
		fmt.Printf("🧪 Test %d/%d [%s]\n", i+1, len(tests), test.category)
		fmt.Printf("   Input:    '%s'\n", test.text)
		fmt.Printf("   Expected: '%s'\n", test.expected)
		fmt.Printf("   Note:     %s\n", test.description)

		// Clear and send test text
		display.Clear()
		display.WriteText(fmt.Sprintf("%d: %s", i+1, test.category))
		display.SetCursor(1, 2) // Second line
		display.WriteText(test.text)

		// Get user confirmation
		worked := getUserConfirmation(test.category + ": " + test.text)

		// Record result
		result := TestResult{
			TestNumber:     i + 1,
			Category:       test.category,
			InputText:      test.text,
			ExpectedOutput: test.expected,
			Description:    test.description,
			Worked:         worked,
		}
		results = append(results, result)

		// Update category stats
		categoryCount[test.category]++
		if worked {
			categoryPassed[test.category]++
		}

		fmt.Println()
	}

	// Display final message
	display.Clear()
	display.WriteText("Latin Test Done!")

	// Print JSON summary
	printJSONSummary(results, categoryCount, categoryPassed)
}

func getUserConfirmation(testName string) bool {
	fmt.Println("   📺 LOOK AT YOUR DISPLAY NOW!")
	fmt.Println("   🎯 SUCCESS = Characters display correctly (native)")
	fmt.Println("   ❌ FAILURE = Garbled/wrong characters")
	fmt.Printf("   ❓ Did characters display correctly on the VFD? (Y/n): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" || input == "y" || input == "yes" {
		fmt.Printf("   ✅ CONFIRMED: %s worked!\n", testName)
		return true
	} else {
		fmt.Printf("   ❌ FAILED: %s didn't work as expected\n", testName)
		return false
	}
}

func printJSONSummary(results []TestResult, categoryCount, categoryPassed map[string]int) {
	totalTests := len(results)
	totalPassed := 0
	for _, result := range results {
		if result.Worked {
			totalPassed++
		}
	}

	successRate := float64(totalPassed) / float64(totalTests) * 100

	// Build category stats
	categoryStats := make(map[string]CategoryStats)
	for category, total := range categoryCount {
		passed := categoryPassed[category]
		rate := float64(passed) / float64(total) * 100
		categoryStats[category] = CategoryStats{
			Total:       total,
			Passed:      passed,
			SuccessRate: rate,
		}
	}

	summary := TestSummary{
		TestSuite:     "Latin Character Test",
		Description:   "Test suite for Latin-based characters (Portuguese, Spanish, French, German, Italian)",
		Timestamp:     "2025-01-XX",
		TotalTests:    totalTests,
		Passed:        totalPassed,
		Failed:        totalTests - totalPassed,
		Errors:        0,
		SuccessRate:   successRate,
		CategoryStats: categoryStats,
		TestResults:   results,
	}

	fmt.Println("======================================================================")
	fmt.Println("📊 LATIN CHARACTER TEST RESULTS - JSON REPORT")
	fmt.Println("======================================================================")
	fmt.Printf("📈 SUMMARY STATISTICS:\n")
	fmt.Printf("   Total Tests: %d\n", totalTests)
	fmt.Printf("   ✅ Passed: %d (%.1f%%)\n", totalPassed, successRate)
	fmt.Printf("   ❌ Failed: %d (%.1f%%)\n", totalTests-totalPassed, 100-successRate)
	fmt.Printf("   💥 Errors: 0 (0.0%%)\n\n")

	fmt.Printf("📊 CATEGORY BREAKDOWN:\n")
	for category, stats := range categoryStats {
		fmt.Printf("   %s: %d/%d (%.1f%%)\n", category, stats.Passed, stats.Total, stats.SuccessRate)
	}

	fmt.Printf("\n📋 DETAILED JSON RESULTS:\n")
	fmt.Println("```json")

	jsonData, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(jsonData))

	fmt.Println("```")
	fmt.Printf("\n🎯 Latin character test completed!\n")
	fmt.Printf("📊 Results: %d/%d tests passed (%.1f%% success rate)\n", totalPassed, totalTests, successRate)
	fmt.Printf("📋 Full JSON report saved above for analysis!\n")

	if successRate >= 80 {
		fmt.Printf("🏆 EXCELLENT! Latin character encoding is working well!\n")
	} else if successRate >= 60 {
		fmt.Printf("⚠️  GOOD: Most characters work, some may need adjustment.\n")
	} else {
		fmt.Printf("❌ ISSUES: Several encoding problems detected.\n")
	}
}
