# Regenerate with: vhs examples/demo/customer-intake.tape
# (run from repo root, with a `pipectl` binary built and on PATH)
Output website/docs/public/demo.gif

Set Shell "bash"
Set FontSize 19
Set Width 1120
Set Height 980
Set Padding 28
Set TypingSpeed 55ms
Set Theme {"name":"pipectl","black":"#0f172a","red":"#e0665c","green":"#4fbf86","yellow":"#e0a83f","blue":"#35c1e0","purple":"#8de3f5","cyan":"#0891b2","white":"#e6edf3","brightBlack":"#45525f","brightRed":"#e0665c","brightGreen":"#4fbf86","brightYellow":"#e0a83f","brightBlue":"#35c1e0","brightPurple":"#8de3f5","brightCyan":"#35c1e0","brightWhite":"#ffffff","background":"#0f172a","foreground":"#e6edf3","cursor":"#35c1e0","selection":"#223047"}

Hide
Type "cd examples/demo && export PS1=$'\n$ ' && clear"
Enter
Show

Type "cat customers.csv"
Enter
Sleep 2500ms

Type "cat customer-intake.yaml"
Enter
Sleep 3500ms

Type "pipectl run customer-intake.yaml < customers.csv"
Sleep 500ms
Enter
Sleep 5s
