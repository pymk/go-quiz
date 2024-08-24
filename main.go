package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const defaultFilePath = "./data/problems.csv"

// main is the entry point of the program.
// It orchestrates the flow of a quiz application that reads questions from a
// CSV file, prompts the user for answers, and calculates the score.
//
// The function performs the following steps:
// 1. Gets the file path for the CSV file containing quiz questions.
// 2. Reads the CSV file, extracting headers and records.
// 3. Iterates through the records, prompting the user for answers to each question.
// 4. Calculates and displays the user's score.
//
// Note:
//   - The program assumes the CSV file is formatted with questions in the first column
//     and correct answers in the second column.
//   - The score is calculated as a percentage of correct answers out of total questions.
func main() {
	filePath, err := getFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Using filepath:", filePath)

	records, _, err := readCSV(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Number of records: %d\n", len(records))

	// Pre-allocate to improve performance
	userAnswers := make([]string, 0, len(records))
	correctAnswers := make([]string, 0, len(records))

	for _, row := range records {
		fmt.Printf("%s?\n", row[0])
		answer, err := recordAnswer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error recording answer: %v\n", err)
			continue
		}
		userAnswers = append(userAnswers, answer)
		correctAnswers = append(correctAnswers, row[1])
	}

	userPoints := calculateScore(userAnswers, correctAnswers)
	userScore := float64(userPoints) / float64(len(records)) * 100

	fmt.Printf("You got %d (%.1f%%) correct!\n", userPoints, userScore)

	os.Exit(0)
}

// getFilePath prompts the user for a file path and returns the validated, absolute path.
//
// The function uses a global variable 'defaultFilePath' which should be defined elsewhere.
//
// Returns:
//   - string: The validated file path. This will be the absolute path to the file.
//   - error: An error if any step of the process fails.
func getFilePath() (string, error) {
	fmt.Printf("Enter file path [%s]: ", defaultFilePath)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}
		// If input stream ends without providing any data (i.e. Ctrl+D)
		return "", fmt.Errorf("no input provided")
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultFilePath, nil
	}

	expandedPath, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("error expanding path: %w", err)
	}

	if _, err := os.Stat(expandedPath); errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("file does not exist: %v", expandedPath)
	}

	return expandedPath, nil
}

// readCSV reads a CSV file and returns its contents as a slice of string slices,
// along with the headers.
//
// Parameters:
//   - filePath: a string representing the path to the CSV file to be read.
//
// Returns:
//   - [][]string: a slice of string slices, where each inner slice represents a row
//     from the CSV file (excluding the header row).
//   - []string: a slice of strings representing the headers from the first row of the CSV file.
//   - error: an error if any step of the reading process fails.
//
// Note:
//   - This function assumes that the CSV file has at a header row.
//   - The expected CSV schema is: question | answer
func readCSV(filePath string) ([][]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading headers: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading records: %w", err)
	}

	return records, headers, nil
}

// recordAnswer prompts the user for input and returns the entered string.
//
// This function reads a single line of text from standard input.
//
// Returns:
//   - answer: a string containing the user's input, with leading and trailing whitespace removed.
//   - err: an error if the input operation fails or if no input is provided.
func recordAnswer() (answer string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}
		// If input stream ends without providing any data (i.e. Ctrl+D)
		return "", fmt.Errorf("no input provided")
	}
	answer = scanner.Text()
	return answer, nil
}

// calculateScore compares user answers to correct answers and returns the number of correct responses.
// It takes two parameters:
//   - userAnswers: a slice of strings representing the user's answers
//   - correctAnswers: a slice of strings representing the correct answers
//
// The function assumes that both slices have the same length and correspond to each other.
// It returns an integer representing the number of correct answers.
func calculateScore(userAnswers, correctAnswers []string) int {
	userPoints := 0
	for i, v := range userAnswers {
		if v == correctAnswers[i] {
			userPoints++
		}
	}
	return userPoints
}
