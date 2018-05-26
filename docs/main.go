package main

import (
	"bytes"
	"github.com/xyproto/splash"
	"io/ioutil"
	"time"
)

var styles = []string{
	"abap",
	"algol",
	"algol_nu",
	"api",
	"arduino",
	"autumn",
	"borland",
	"bw",
	"colorful",
	"dracula",
	"emacs",
	"friendly",
	"fruity",
	"github",
	"igor",
	"lovelace",
	"manni",
	"monokai",
	"monokailight",
	"murphy",
	"native",
	"paraiso-dark",
	"paraiso-light",
	"pastie",
	"perldoc",
	"pygments",
	"rainbow_dash",
	"rrt",
	"solarized-dark",
	"solarized-dark256",
	"solarized-light",
	"swapoff",
	"tango",
	"trac",
	"vim",
	"vs",
	"xcode",
}

const (
	title         = "Chroma Style Gallery"
	simpleCSS     = "body { font-family: sans-serif; margin: 4em; } .chroma { padding: 1em; } #main-headline { border-bottom: 3px solid red; margin-bottom: 2em; } a { color: #1E385B; } a:visited { color: #1E385B; } a:hover { color: #1D6142; }"
	sampleContent = `<code><pre>
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
</pre></code>`
)

func footer() string {
	var buf bytes.Buffer
	buf.WriteString("<small>Generated ")
	buf.WriteString(time.Now().UTC().Format(time.RFC3339)[:10])
	buf.WriteString(" by <a href=\"https://github.com/xyproto\">xyproto</a> using <a href=\"https://github.com/xyproto/splash\">splash</a>.")
	buf.WriteString("</small></body></html>")
	return buf.String()
}

func main() {

	for i, styleName := range styles {

		// Generate a HTML document for the current style name
		var inputBuffer bytes.Buffer
		inputBuffer.WriteString("<!doctype html><html><head><title>")
		inputBuffer.WriteString(styleName)
		inputBuffer.WriteString("</title><style>")
		inputBuffer.WriteString(simpleCSS)
		inputBuffer.WriteString("</style></head><body>")
		inputBuffer.WriteString("<h1>")
		inputBuffer.WriteString(styleName)
		inputBuffer.WriteString("</h1>")
		inputBuffer.WriteString(sampleContent)

		// Button to the previous style, if possible
		if i > 0 {
			prevName := styles[i-1]
			inputBuffer.WriteString("<button onClick=\"location.href='" + prevName + ".html'\">Prev</button>")
		} else {
			inputBuffer.WriteString("<button disabled='true'>Prev</button>")
		}

		// Button to the next style, if possible
		if i < (len(styles) - 1) {
			nextName := styles[i+1]
			inputBuffer.WriteString("<button onClick=\"location.href='" + nextName + ".html'\">Next</button>")
		} else {
			inputBuffer.WriteString("<button disabled='true'>Next</button>")
		}

		// Button to the overview
		inputBuffer.WriteString("<button onClick=\"location.href='index.html'\">Up</button>")

		// Button to the single page
		inputBuffer.WriteString("<button onClick=\"location.href='all.html'\">All</button>")

		inputBuffer.WriteString("</body></html>")

		// Highlight the source code in the HTML with the current style
		htmlBytes, err := splash.Splash(inputBuffer.Bytes(), styleName)
		if err != nil {
			panic(err)
		}

		// Write the HTML sample for this style name
		err = ioutil.WriteFile(styleName+".html", htmlBytes, 0644)
		if err != nil {
			panic(err)
		}

	}

	// Generate an Index file for viewing the different styles
	var buf bytes.Buffer
	buf.WriteString("<!doctype html><html><head><title>")
	buf.WriteString(title)
	buf.WriteString("</title><style>")
	buf.WriteString(simpleCSS)
	buf.WriteString("</style></head><body><h1 id='main-headline'>")
	buf.WriteString(title)
	buf.WriteString("</h1><ul>")
	for _, styleName := range styles {
		buf.WriteString("<li><a href=\"" + styleName + ".html\" alt=\"" + styleName + " style\">" + styleName + "</a></li>")
	}
	buf.WriteString("</ul>")
	buf.WriteString("<p><a href='all.html' alt='All styles on one page'>All styles on one page</a></p>")
	buf.WriteString(footer())
	err := ioutil.WriteFile("index.html", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	// Generate a single page with all the styles
	buf = bytes.Buffer{}
	buf.WriteString("<!doctype html><html><head><title>")
	buf.WriteString(title)
	buf.WriteString("</title><style>")
	css := "body { font-family: sans-serif; margin: 4em; } h1 { border-bottom: 3px solid red; margin-bottom: 2em; } pre { display: inline-block; margin: 2em; padding: 3em 5em 3em 2em; box-shadow: 5px 5px 5px rgba(68, 68, 68, 0.6); border-radius: 7px; border: 2px solid black;} #stylelink:link { text-decoration: none; color: black; } #stylelink:visited { color: black; } #stylelink:hover { color: #1D6142; } a { color: #1E385B; }"
	buf.WriteString(css)
	buf.WriteString("</style></head><body><h1>")
	buf.WriteString(title)
	buf.WriteString("</h1>")
	for _, styleName := range styles {
		buf.WriteString("<h2>")
		buf.WriteString("<a id='stylelink' href='" + styleName + ".html' alt='View only " + styleName + "'>" + styleName + "</a>")
		buf.WriteString("</h2>")

		htmlBytes, cssBytes, err := splash.Highlight([]byte(sampleContent), styleName, false)
		if err != nil {
			panic(err)
		}
		cssElements := bytes.Split(cssBytes, []byte("}"))
		for _, element := range cssElements {
			elements := bytes.Split(bytes.Replace(element, []byte(".chroma"), []byte(""), -1), []byte("{"))
			if len(elements) == 2 {
				className := bytes.TrimSpace(elements[0])
				if bytes.HasPrefix(className, []byte(".")) {
					className = className[1:]
				}
				style := bytes.TrimSpace(elements[1])
				if len(className) == 0 {
					htmlBytes = bytes.Replace(htmlBytes, []byte("<pre "), []byte("<pre style=\""+string(style)+"\" "), -1)
				} else {
					htmlBytes = bytes.Replace(htmlBytes, []byte("<span class=\""+string(className)+"\""), []byte("<span style=\""+string(style)+"\""), -1)
				}
			}
		}
		buf.Write(htmlBytes)
	}
	buf.WriteString("<p><a id='bottom' href='index.html' alt='Back to overview'>Back to overview</a></p>")
	buf.WriteString(footer())
	err = ioutil.WriteFile("all.html", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

}
