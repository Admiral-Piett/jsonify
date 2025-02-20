# jsonify

Paste partial or complete JSON data into the input field.  Key names with dot.notation are accepted.  We will then 
parse that into true json and display the output.

![Demo Gif](static/demo.gif)

## Warning
We try to do some clean up for you in terms of quoting things, correcting single quotes VS double quotes, etc. 
**HOWEVER** we're not perfect.  Not having your commas or colon fields placed correctly will wreak havoc.  
We need those as anchors to figure out what's going on, and what your structure is.  So missing or misplacing 
those characters is fair game.

### For Now and Your Sanity's Sake
You will want to try to strip both `:` and `'` out of your keys and values.  
Otherwise, it'll almost certainly chop them up.

## Hot Keys
- `Cmd + q` - Exit

#### Select Input Area
- `Cmd + Enter` - Submits input
- `Cmd + Backspace` - Clears input

#### Select Output Area
- `Cmd + c` - Copied output
