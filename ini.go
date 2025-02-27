package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

func usage() {
	fmt.Printf("Usage: %v --file <ini_file> --read section[.key] || --write section.key=value [--show-section] [--show-export]\n",
		filepath.Base(os.Args[0]))
}

// ini_print formats the output depending on the OS and --show-export flag.
func iniPrint(key string, value string, osExport bool) {
	if osExport {
		if runtime.GOOS != "windows" {
			fmt.Printf("export %s=%q\n", key, value) // Using %q to ensure proper formatting for shell variables.
		} else {
			fmt.Printf("setx %s %q\n", key, value)
		}
	} else {
		fmt.Printf("%s=%s\n", key, value)
	}
}

func main() {
	iniFile := flag.String("file", "", "INI file to process.")
	read := flag.String("read", "", "Section or section.key to read.")
	write := flag.String("write", "", "Section.key=value to write.")
	showSection := flag.Bool("show-section", false, "Show section name when reading.")
	showExport := flag.Bool("show-export", false, "Format output for shell export.")

	flag.Parse()

	// Validate required flags
	if *iniFile == "" {
		usage()
		os.Exit(1)
	}
	if *write != "" && *read != "" {
		fmt.Fprintln(os.Stderr, "Error: Cannot use --read and --write together.")
		usage()
		os.Exit(1)
	}

	// Load INI file if it exists
	var cfg *ini.File
	var err error

	if _, err := os.Stat(*iniFile); err == nil {
		cfg, err = ini.Load(*iniFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to read INI file: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg = ini.Empty() // Create a new INI structure if file doesn't exist.
	}

	// Handle writing to the INI file
	if *write != "" {
		parts := strings.SplitN(*write, ".", 2)
		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, "Error: Invalid format. Use <section>.<key>=<value>")
			os.Exit(1)
		}
		sectionName := parts[0]
		kv := strings.SplitN(parts[1], "=", 2)
		if len(kv) != 2 {
			fmt.Fprintln(os.Stderr, "Error: Invalid format. Use <section>.<key>=<value>")
			os.Exit(1)
		}
		key, value := kv[0], kv[1]

		cfg.Section(sectionName).Key(key).SetValue(value)
		err = cfg.SaveTo(*iniFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to save INI file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated [%s] %s=%s\n", sectionName, key, value)
	}

	// Handle reading from the INI file
	if *read != "" {
		parts := strings.SplitN(*read, ".", 2)
		sectionName := parts[0]

		// Check if section exists
		section, err := cfg.GetSection(sectionName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Section [%s] not found\n", sectionName)
			os.Exit(1)
		}

		// Show section header if requested
		if *showSection && !*showExport {
			fmt.Printf("[%s]\n", sectionName)
		}

		// If only the section name is given, list all keys
		if len(parts) == 1 {
			for _, key := range section.Keys() {
				iniPrint(key.Name(), key.String(), *showExport)
			}
		} else {
			// If section.key is specified
			key, err := section.GetKey(parts[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Key '%s' not found in section [%s]\n", parts[1], sectionName)
				os.Exit(1)
			}
			iniPrint(key.Name(), key.String(), *showExport)
		}
	}
}
