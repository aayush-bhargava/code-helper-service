package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RunResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type AutoFixResponse struct {
	OriginalCode string `json:"original_code"`
	FixedCode    string `json:"fixed_code"`
}

type HelpResponse struct {
	Query           string   `json:"query"`
	MatchedKeywords []string `json:"matched_keywords"`
	Tips            []string `json:"tips"`
	Message         string   `json:"message"`
}

type helpRule struct {
	Keywords []string
	Message  string
}

var helpRules = []helpRule{
	{[]string{"loop", "for", "while"}, "Loops repeat code. Ensure exit conditions to avoid infinite loops."},
	{[]string{"if", "else", "condition"}, "If conditions run code only when true. Check comparison operators."},
	{[]string{"function", "method"}, "Functions should be single-purpose. Use parameters over globals."},
	{[]string{"error", "debug"}, "Use print or logs and test in small steps to debug."},
	{[]string{"variable", "name", "naming"}, "Use meaningful variable names like totalPrice or count."},
	{[]string{"bracket", "brace", "parenthesis"}, "Align brackets using consistent indentation."},
}

func main() {
	app := fiber.New()
	app.Static("/", "./public")
	app.Post("/run", handleRun)
	app.Post("/autofix", handleAutoFix)
	app.Post("/help", handleHelp)
	log.Println("http://localhost:3000")
	app.Listen(":3000")
}

func handleRun(c *fiber.Ctx) error {
	code := strings.TrimSpace(string(c.Body()))
	if code == "" {
		return c.JSON(RunResponse{Error: "No code provided"})
	}
	if strings.Contains(strings.ToLower(code), "error") {
		return c.JSON(RunResponse{Error: "Simulated compiler error"})
	}
	lines := strings.Split(code, "\n")
	var out []string
	for i, l := range lines {
		out = append(out, "Line "+intToString(i+1)+": "+strings.TrimSpace(l))
	}
	return c.JSON(RunResponse{Output: strings.Join(out, "\n")})
}

func handleAutoFix(c *fiber.Ctx) error {
	orig := string(c.Body())
	return c.JSON(AutoFixResponse{
		OriginalCode: orig,
		FixedCode:    fix(orig),
	})
}

func handleHelp(c *fiber.Ctx) error {
	q := strings.ToLower(strings.TrimSpace(string(c.Body())))
	if q == "" {
		return c.JSON(HelpResponse{Message: "Empty query"})
	}
	var m []string
	var t []string
	for _, r := range helpRules {
		for _, k := range r.Keywords {
			if strings.Contains(q, k) {
				m = append(m, k)
				t = append(t, r.Message)
				break
			}
		}
	}
	if len(t) == 0 {
		return c.JSON(HelpResponse{Query: q, Message: "No matches"})
	}
	return c.JSON(HelpResponse{Query: q, MatchedKeywords: m, Tips: t, Message: "Matched"})
}

func fix(code string) string {
	ls := strings.Split(code, "\n")
	space := regexp.MustCompile(`\s+`)
	ind := 0
	var out []string
	for _, raw := range ls {
		l := strings.TrimSpace(raw)
		if l == "" {
			out = append(out, "")
			continue
		}
		if strings.HasPrefix(l, "}") && ind > 0 {
			ind--
		}
		l = space.ReplaceAllString(l, " ")
		if needSemicolon(l) {
			l = l + ";"
		}
		out = append(out, strings.Repeat(" ", ind*4)+l)
		ind += strings.Count(l, "{") - strings.Count(l, "}")
		if ind < 0 {
			ind = 0
		}
	}
	for ind > 0 {
		ind--
		out = append(out, strings.Repeat(" ", ind*4)+"}")
	}
	return strings.Join(out, "\n")
}

func needSemicolon(l string) bool {
	l = strings.TrimSpace(l)
	if l == "" {
		return false
	}
	c := l[len(l)-1]
	if c == ';' || c == '{' || c == '}' || c == ':' || c == ',' {
		return false
	}
	x := strings.ToLower(l)
	if strings.HasPrefix(x, "if") || strings.HasPrefix(x, "for") || strings.HasPrefix(x, "while") || strings.HasPrefix(x, "function") {
		return false
	}
	return true
}

func intToString(n int) string {
	d := "0123456789"
	if n == 0 {
		return "0"
	}
	var o []byte
	for n > 0 {
		o = append([]byte{d[n%10]}, o...)
		n /= 10
	}
	return string(o)
}
