package main

import (
    "encoding/json"
    "fmt"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "os/exec"
    "strings"
)

const maxInputLength = 100_000

var app = tview.NewApplication()

var input = tview.NewForm()
var output = tview.NewForm()
var pages = tview.NewPages()

func main() {
    app.SetInputCapture(hotKeyParser)

    outputView := tview.NewTextView()
    outputView.SetDynamicColors(true).SetBorder(true).SetTitle("Output").SetTitleAlign(tview.AlignLeft)
    outputView.SetSize(20, 0)

    output.AddFormItem(outputView)
    output.AddButton("Back", back)
    output.AddButton("Copy", func() {
        text := output.GetFormItem(0).(*tview.TextView).GetText(false)
        setToClipboard(text)
    })

    input.AddTextArea("Input", "", 0, 20, maxInputLength, nil)
    input.AddButton("Submit", submitInput)
    input.AddButton("Clear", func() {
        input.GetFormItem(0).(*tview.TextArea).SetText("", true)
    })
    input.AddButton("Quit", func() {
        app.Stop()
    })
    input.SetBorder(true).SetTitle("JSONify").SetTitleAlign(tview.AlignLeft)

    pages.AddPage("Input", input, true, true)
    pages.AddPage("Output", output, true, false)

    if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
        panic(err)
    }
}

func submitInput() {
    text := input.GetFormItem(0).(*tview.TextArea).GetText()
    w := output.GetFormItem(0).(*tview.TextView).BatchWriter()
    parse(w, text)
    pages.SwitchToPage("Output")
}

func back() {
    pages.SwitchToPage("Input")
}

func sanitizeText(text string) string {
    text = strings.TrimSpace(text)
    if !strings.HasPrefix(text, "{") {
        text = "{" + text
    }
    if !strings.HasSuffix(text, "}") {
        if strings.HasSuffix(text, ",") {
            text = text[:len(text)-1]
        }
        text = text + "}"
    }
    return text
}

func parse(w tview.TextViewWriter, text string) {
    defer w.Close()
    w.Clear()

    data := &map[string]interface{}{}

    text = sanitizeText(text)

    err := json.Unmarshal([]byte(text), data)
    if err != nil {
        fmt.Fprintf(w, "Invalid input: \n%s", err.Error())
        return
    }

    result := make(map[string]interface{})
    for path, value := range *data {
        if path == "" {
            continue
        }
        setNestedValue(result, path, value)
    }

    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        fmt.Fprintf(w, "Unable to marshall result: \n%s", err.Error())
        return
    }
    jsonStr := string(jsonData)
    fmt.Fprintf(w, "%s", jsonStr)
}

func setNestedValue(data map[string]interface{}, path string, value interface{}) {
    keys := strings.Split(path, ".")
    current := data

    for _, key := range keys[:len(keys)-1] {
        if _, ok := current[key]; !ok {
            current[key] = make(map[string]interface{})
        }
        current = current[key].(map[string]interface{})
    }

    current[keys[len(keys)-1]] = value
}

func setToClipboard(content string) error {
    cmd := exec.Command("pbcopy")
    in, err := cmd.StdinPipe()
    if err != nil {
        return err
    }
    if err := cmd.Start(); err != nil {
        return err
    }
    if _, err := in.Write([]byte(content)); err != nil {
        return err
    }
    if err := in.Close(); err != nil {
        return err
    }
    return cmd.Wait()
}

func hotKeyParser(event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyEnter && event.Modifiers() == tcell.ModAlt {
        submitInput()
        return event
    }
    if event.Key() == tcell.KeyDEL && event.Modifiers() == tcell.ModAlt {
        back()
        return event
    }
    return event
}
