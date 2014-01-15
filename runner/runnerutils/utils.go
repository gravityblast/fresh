package runnerutils

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var logFilePath string

func init() {
	root := os.Getenv("RUNNER_WD")
	tmpPath := os.Getenv("RUNNER_TMP_PATH")
	fileName := os.Getenv("RUNNER_BUILD_LOG")
	logFilePath = filepath.Join(root, tmpPath, fileName)
}

// Returns true if a build error file exists in the tmp folder.
func HasErrors() bool {
	if _, err := os.Stat(logFilePath); err == nil {
		return true
	}

	return false
}

// It renders an error page with the build error message.
func RenderError(w http.ResponseWriter) {
	data := map[string]interface{}{
		"Output": readErrorFile(),
	}

	w.Header().Set("Content-Type", "text/html")
	tpl := template.Must(template.New("ErrorPage").Parse(buildPageTpl))
	tpl.Execute(w, data)
}

func readErrorFile() string {
	file, err := os.Open(logFilePath)
	if err != nil {
		return ""
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	bytes, _ := ioutil.ReadAll(reader)

	return string(bytes)
}

const buildPageTpl string = `
  <html>
    <head>
      <title>Traffic Panic</title>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
      <style>
      html, body{ padding: 0; margin: 0; }
      header { background: #C52F24; color: white; border-bottom: 2px solid #9C0606; }
      h1 { padding: 10px 0; margin: 0; }
      .container { margin: 0 20px; }
      .output { height: 300px; overflow-y: scroll; border: 1px solid #e5e5e5; padding: 10px; }
      </style>
    </head>
  <body>
    <header>
      <div class="container">
        <h1>Build Error</h1>
      </div>
    </header>

    <div class="container">
      <pre class="output">{{ .Output }}</pre>
    </div>
  </body>
  </html>
  `
