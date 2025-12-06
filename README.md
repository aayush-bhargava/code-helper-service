# Code Helper Service

## Run
go run main.go

## URL
http://localhost:3000

## Routes
POST /run
POST /autofix
POST /help

### Auto-Fix Rules
- Adds missing semicolons
- Fixes indentation after `{` and before `}`
- Collapses multiple spaces into single space
- Appends `}` if unmatched opening braces exist

### Help Keywords
| Keywords | Help Message |
|----------|--------------|
| loop, for, while | Loops repeat code, ensure exit condition |
| if, else, condition | Use correct comparison operators |
| function, method | Keep functions single-purpose |
| error, debug | Test small chunks, add logs |
| variable, name, naming | Use meaningful variable names |
| bracket, brace, parenthesis | Keep braces aligned consistently |
