import fs from "node:fs"
import path from "node:path"
import pug from "pug"

export default function plugin() {
  return {
    name: "rollup-plugin-pug",
    buildStart() {
      this.addWatchFile(path.join(__dirname, "index.pug"))
    },
    generateBundle() {
      const compile = pug.compileFile(path.join(__dirname, "index.pug"))
      const contents = compile()
      fs.writeFileSync(path.join(__dirname, "index.html"), contents)
    },
  }
}