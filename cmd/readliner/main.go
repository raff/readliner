package main

import {
    "github.com/raff/readliner"
}

func main() {
	rl := readliner.New("> ", "")
	defer rl.Close()

	scanner := bufio.NewScanner(rl)

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}

		fmt.Println("scanned:", text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read error:", err)
	}
}