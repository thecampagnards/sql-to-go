const inEditor = CodeMirror.fromTextArea(document.getElementById("in"), {
  mode: "sql",
  lineNumbers: true,
  theme: "monokai"
});
inEditor.setSize("45%", null)
inEditor.save()

var outEditor = CodeMirror.fromTextArea(document.getElementById("out"), {
  mode: "go",
  lineNumbers: true,
  theme: "monokai"
});
outEditor.setSize("45%", null)
outEditor.save()

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    inEditor.on("change", (event) => {
      try {
        const r = window.parse(event.doc.getValue());
        console.log(r)
        outEditor.getDoc().setValue(r);
      } catch (err) {
        document.getElementById("error").textContent = err;
        console.log("caught error:", err);
      }
    });
});
