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
  theme: "monokai",
  readOnly: "true",
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
    inEditor.getDoc().setValue(`
CREATE TABLE user
(
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    city VARCHAR(255),
    country VARCHAR(255),
    date_of_birth DATE,
    email VARCHAR(255),
    name VARCHAR(100),
    postal_code VARCHAR(5),
    surname VARCHAR(100),
    test_alter_type INT,
    test_alter_rename INT
);

ALTER TABLE user ADD example VARCHAR(255);

ALTER TABLE user DROP COLUMN country;

ALTER TABLE user ALTER COLUMN test_alter_type TYPE VARCHAR(255);

ALTER TABLE user RENAME COLUMN test_alter_rename TO test_alter_rename_new;

CREATE TABLE book
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    user_id BIGINT REFERENCES user(id),
    test_alter_reference_id BIGINT REFERENCES user(id)
);`)
});
