package main

import (
    "cogentcore.org/core/colors"
    "cogentcore.org/core/core"
    "cogentcore.org/core/events"
    "cogentcore.org/core/styles"
    "cogentcore.org/core/styles/units"
    "cogentcore.org/core/texteditor"
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
        formattedText, err := parse(inputEditor.Buffer.String())
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
