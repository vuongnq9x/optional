package example

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vuongnq9x/optional"
)

func main() {
	// Ví dụ cơ bản
	fmt.Println("=== Basic Usage ===")
	someValue := optional.Some(42)
	noneValue := optional.None[int]()

	fmt.Printf("Some value: %s\n", someValue.String())
	fmt.Printf("None value: %s\n", noneValue.String())

	// Sử dụng OrElse
	fmt.Println("\n=== OrElse ===")
	fmt.Printf("Some.OrElse(0): %d\n", someValue.OrElse(0))
	fmt.Printf("None.OrElse(100): %d\n", noneValue.OrElse(100))

	// Map transformation
	fmt.Println("\n=== Map ===")
	stringOpt := optional.Map(someValue, func(x int) string {
		return "Number: " + strconv.Itoa(x)
	})
	fmt.Printf("Mapped: %s\n", stringOpt.String())

	// Chain operations
	fmt.Println("\n=== Chain Operations ===")
	result := optional.Map(someValue, func(x int) int { return x * 2 })
	result = optional.Map(result, func(x int) int { return x + 10 })
	fmt.Printf("42 * 2 + 10 = %s\n", result.String())

	// JSON serialization
	fmt.Println("\n=== JSON ===")
	data, _ := json.Marshal(someValue)
	fmt.Printf("JSON: %s\n", string(data))

	noneData, _ := json.Marshal(noneValue)
	fmt.Printf("None JSON: %s\n", string(noneData))

	// Practical example: parsing user input
	fmt.Println("\n=== Practical Example ===")
	userInputs := []string{"42", "invalid", "100"}

	for _, input := range userInputs {
		parsed := parseNumber(input)
		parsed.IfPresent(func(num int) {
			fmt.Printf("Parsed '%s' as %d\n", input, num)
		})

		if parsed.IsEmpty() {
			fmt.Printf("Failed to parse '%s'\n", input)
		}
	}
}

func parseNumber(s string) optional.Optional[int] {
	if num, err := strconv.Atoi(s); err == nil {
		return optional.Some(num)
	}
	return optional.None[int]()
}
