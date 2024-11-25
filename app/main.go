package main

import (
    "cogentcore.org/core/colors"
    "cogentcore.org/core/core"
    "cogentcore.org/core/events"
    "cogentcore.org/core/styles"
    "cogentcore.org/core/styles/units"
    "cogentcore.org/core/texteditor"
    "fmt"
    "os/exec"
)

func main() {
    b := core.NewBody("JSONify")
    splits := core.NewSplits(b)

    input := core.NewFrame(splits)
    input.Styler(splitEntryStyler)

    inputHeader := core.NewText(input).SetText("Input")
    inputHeader.Styler(headerStyler)

    inputEditor := texteditor.NewEditor(input)
    inputEditor.Styler(inputTextEditorStyler)

    inputBtnFrame := core.NewFrame(input)
    inputBtnFrame.Styler(btnRowStyler)

    inputSubmitBtn := core.NewButton(inputBtnFrame)
    inputSubmitBtn.Styler(func(s *styles.Style) {
        s.Background = colors.Scheme.Success.Base
    })

    inputClearBtn := core.NewButton(inputBtnFrame)
    inputClearBtn.Styler(func(s *styles.Style) {
        s.Background = colors.Scheme.Warn.Base
    })

    output := core.NewFrame(splits)
    output.Styler(splitEntryStyler)

    outputHeader := core.NewText(output)
    outputHeader.SetText("Output")
    outputHeader.Styler(headerStyler)

    outputTextAreaContainer := core.NewFrame(output)
    outputTextAreaContainer.Styler(inputTextEditorStyler)
    outputTextArea := core.NewText(outputTextAreaContainer)
    outputTextArea.Styler(outputTextAreaStyler)

    outputBtnFrame := core.NewFrame(output)
    outputBtnFrame.Styler(func(s *styles.Style) {
        s.Grow.Set(1, 0)
        s.Justify.Content = styles.End
    })

    outputCopyBtn := core.NewButton(outputBtnFrame)

    // --- Event Handlers
    inputSubmitBtn.SetText("Submit").OnClick(func(e events.Event) {
        fmt.Println("submit")
        inputText := string(inputEditor.Buffer.Text())

        formattedText, err := parse(inputText)
        if err != nil {
            core.ErrorDialog(b, err, "Unable to parse input")
            return
        }

        outputTextArea.SetText(formattedText)
        outputTextArea.Update()
    })
    inputClearBtn.SetText("Clear").OnClick(func(e events.Event) {
        inputEditor.Buffer.SetText([]byte(""))
    })
    outputCopyBtn.SetText("Copy").OnClick(func(e events.Event) {
        err := setToClipboard(outputTextArea.Text)
        if err != nil {
            core.ErrorSnackbar(b, err, "Error copying content")
        }
        core.MessageSnackbar(b, "Output copied")
    })
    b.RunMainWindow()
}

func splitEntryStyler(s *styles.Style) {
    s.Direction = styles.Column
    s.Overflow.Set(styles.OverflowHidden)
}

func headerStyler(s *styles.Style) {
    s.Font.Size.Set(1.5, units.UnitRem)
    s.Font.Weight = styles.WeightSemiBold
}

func inputTextEditorStyler(s *styles.Style) {
    s.Grow.Set(1, 1)
}

func outputTextAreaStyler(s *styles.Style) {
    s.Text.WhiteSpace = styles.WhiteSpacePre
}

func btnRowStyler(s *styles.Style) {
    s.Direction = styles.Row
}

//func main() {
//    app.SetInputCapture(hotKeyParser)
//
//    outputView := tview.NewTextView()
//    outputView.SetDynamicColors(true).SetBorder(true).SetTitle("Output").SetTitleAlign(tview.AlignLeft)
//    outputView.SetSize(20, 0)
//
//    output.AddFormItem(outputView)
//    output.AddButton("Back", back)
//    output.AddButton("Copy", func() {
//        text := output.GetFormItem(0).(*tview.TextView).GetText(false)
//        setToClipboard(text)
//    })
//
//    //input.AddTextArea("Input", "", 0, 20, maxInputLength, nil)
//    input.AddTextArea("Input", "", 0, 5, maxInputLength, nil)
//    input.AddButton("Submit", submitInput)
//    input.AddButton("Clear", clearInput)
//    input.AddButton("Quit", func() {
//        app.Stop()
//    })
//    input.SetBorder(true).SetTitle("JSONify").SetTitleAlign(tview.AlignLeft)
//
//    pages.AddPage("Input", input, true, true)
//    pages.AddPage("Output", output, true, false)
//
//    if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
//        panic(err)
//    }
//}

//func submitInput() {
//    text := input.GetFormItem(0).(*tview.TextArea).GetText()
//    w := output.GetFormItem(0).(*tview.TextView).BatchWriter()
//    parse(w, text)
//    pages.SwitchToPage("Output")
//}
//
//func back() {
//    pages.SwitchToPage("Input")
//}
//
//func clearInput() {
//    input.GetFormItem(0).(*tview.TextArea).SetText("", true)
//}
//
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

//
//func hotKeyParser(event *tcell.EventKey) *tcell.EventKey {
//    if event.Key() == tcell.KeyEnter && event.Modifiers() == tcell.ModAlt {
//        submitInput()
//        return event
//    }
//    if event.Key() == tcell.KeyDEL && event.Modifiers() == tcell.ModAlt {
//        back()
//        return event
//    }
//    // 100 is "d" as a `rune`, the `EventKey` stores all ASCII characters like that.
//    if event.Rune() == 100 && event.Modifiers() == tcell.ModAlt {
//        clearInput()
//        return event
//    }
//    return event
//}
