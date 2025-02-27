# INI File Manager

## Overview
This command-line tool allows reading and writing INI files efficiently. It provides options to read sections or specific keys, write new key-value pairs, and format output for shell export.

## Installation

Ensure you have Go installed, then build the application:

```sh
# Clone or download the repository
cd path/to/project

go build -o ini-manager ini.go
```

## Usage

```sh
Usage: ini-manager --file <ini_file> --read section[.key] || --write section.key=value [--show-section] [--show-export]
```

### Options

| Flag            | Description |
|----------------|-------------|
| `--file`       | Path to the INI file to process. **Required**. |
| `--read`       | Read a section or a specific `section.key`. |
| `--write`      | Write a new value to the INI file in `section.key=value` format. |
| `--show-section` | Display the section header when reading keys. |
| `--show-export` | Format output for shell export (Linux/macOS: `export KEY=value`, Windows: `setx KEY value`). |

## Examples

### Reading Values

#### Read all keys in a section:
```sh
ini-manager --file config.ini --read database
```
_Output:_
```ini
host=localhost
port=5432
user=admin
```

#### Read a specific key:
```sh
ini-manager --file config.ini --read database.port
```
_Output:_
```sh
port=5432
```

#### Read a key and format output for shell export:
```sh
ini-manager --file config.ini --read database.user --show-export
```
_Output (Linux/macOS):_
```sh
export user="admin"
```
_Output (Windows):_
```sh
setx user "admin"
```

### Writing Values

#### Update or add a new key:
```sh
ini-manager --file config.ini --write database.password=securepass
```
_Output:_
```sh
Updated [database] password=securepass
```

## Error Handling
- If the INI file does not exist, the tool creates a new one.
- If a section or key is not found, an error message is displayed.
- `--read` and `--write` cannot be used together.

## Author
Konstantinos Patronas

