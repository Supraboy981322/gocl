# GoCL (Go Custom Language) 

a user-defined programming "language" transpiler

>[!WARNING]
>THIS IS ***SUPER*** EARLY, because of this, the repo is private, and the project can't do much yet

---

basically, you create a definitions file like this:
```gomn
["fn"] := "func"
["prim()"] := "main()"
["wr"] := |
  ["wr"] := "fmt"
  ["l"] := "Println"
|
```

which (in this example) transpiles the following code to Go:
```
fn prim() {
  wr.l("foo bar baz qux")
}
```

your definitions can be anything, you're not limited to transpiling to Go code 
