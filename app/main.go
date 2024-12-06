package main

import (
    "cogentcore.org/core/colors"
    "cogentcore.org/core/core"
    "cogentcore.org/core/events"
    "cogentcore.org/core/styles"
    "cogentcore.org/core/styles/abilities"
    "cogentcore.org/core/styles/states"
    "cogentcore.org/core/styles/units"
    "cogentcore.org/core/texteditor"
    "github.com/Admiral-Piett/jsonify/app/parse"
    "os/exec"
)

const CmdEnterKeyChord = "Meta+ReturnEnter"
const CmdBackspaceKeyChord = "Meta+Backspace"
const CmdCKeyChord = "Meta+C"

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
    outputTextAreaContainer.Styler(outputTextContainerStyler)
    outputTextArea := core.NewText(outputTextAreaContainer)
    outputTextArea.Styler(outputTextAreaStyler)

    outputBtnFrame := core.NewFrame(output)
    outputBtnFrame.Styler(func(s *styles.Style) {
        s.Grow.Set(1, 0)
        s.Justify.Content = styles.End
    })

    outputCopyBtn := core.NewButton(outputBtnFrame)

    // --- Event Handlers
    // TODO - clean this uuuup, do we really need to be so javascripty?
    onSubmit := func(e events.Event) {
        formattedText, err := parse.Parse(inputEditor.Buffer.String())
        if err != nil {
            core.ErrorDialog(b, err, "Unable to parse input")
            return
        }

        outputTextArea.SetText(formattedText)
        outputTextArea.Update()
    }
    onClear := func(e events.Event) {
        inputEditor.Buffer.SetText([]byte(""))
    }
    onCopy := func(e events.Event) {
        err := setToClipboard(outputTextArea.Text)
        if err != nil {
            core.ErrorSnackbar(b, err, "Error copying content")
        }
        core.MessageSnackbar(b, "Output copied")
    }

    inputSubmitBtn.SetText("Submit").OnClick(onSubmit)
    inputClearBtn.SetText("Clear").OnClick(onClear)
    outputCopyBtn.SetText("Copy").OnClick(onCopy)

    // Hot Keys
    inputEditor.OnKeyChord(func(e events.Event) {
        if CmdEnterKeyChord == e.KeyChord() {
            onSubmit(e)
        }
        if CmdBackspaceKeyChord == e.KeyChord() {
            onClear(e)
        }
    })

    // We've got to do this on `outputTextAreaContainer` and `outputTextArea` because they may push each
    //other out of the way, and we don't care which they're on really, just that they're in the right-ish place.
    outputTextAreaContainer.OnKeyChord(func(e events.Event) {
        if CmdCKeyChord == e.KeyChord() {
            onCopy(e)
        }
    })
    outputTextArea.OnKeyChord(func(e events.Event) {
        if CmdCKeyChord == e.KeyChord() {
            onCopy(e)
        }
    })

    // --- Run
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

func outputTextContainerStyler(s *styles.Style) {
    s.Grow.Set(1, 1)
    s.SetAbilities(true, abilities.Selectable, abilities.Focusable, abilities.Hoverable)
}

func outputTextAreaStyler(s *styles.Style) {
    s.Grow.Set(1, 1)
    s.Abilities.SetFlag(true, abilities.Focusable)
    s.MaxBoxShadow = styles.BoxShadow1()
    if s.Is(states.Hovered) {
        s.BoxShadow = s.MaxBoxShadow
    }
    s.Text.WhiteSpace = styles.WhiteSpacePre
}

func btnRowStyler(s *styles.Style) {
    s.Direction = styles.Row
}

// TODO - can you use widgetevents.WidgetBase.Clipboard() instead?
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
