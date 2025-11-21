# GoCL (Go Custom Language) 

a user-defined programming "language" transpiler

>[!WARNING]
>THIS IS ***SUPER*** EARLY and I haven't done much testing yet, better documentation will be made once I do more testing and feel that it's stable enough

---

basically, you create definitions in a `defs.gomn` file like this:
```gomn
["fn"] := "func"
["prim()"] := "main()"
["wr"] := |
  ["wr"] := "fmt"
  ["l"] := "Println"
|
```
<sub><a href="https://github.com/Supraboy981322/gomn">in-case you're wondering what `gomn` is</a></sub>

which (in this example) transpiles the following code to Go:
```
fn prim() {
  wr.l("foo bar baz qux")
}
```

your definitions can be anything, you're not limited to transpiling to Go code 
