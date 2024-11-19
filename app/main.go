package main

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "os/exec"
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

    //input.AddTextArea("Input", "", 0, 20, maxInputLength, nil)
    input.AddTextArea("Input", "", 0, 5, maxInputLength, nil)
    input.AddButton("Submit", submitInput)
    input.AddButton("Clear", clearInput)
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

func clearInput() {
    input.GetFormItem(0).(*tview.TextArea).SetText("", true)
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
    // 100 is "d" as a `rune`, the `EventKey` stores all ASCII characters like that.
    if event.Rune() == 100 && event.Modifiers() == tcell.ModAlt {
        clearInput()
        return event
    }
    return event
}
