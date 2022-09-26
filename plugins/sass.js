import path from "node:path"
import fs from "node:fs"
import sass from "sass"

export default function plugin() {
  return {
    name: "rollup-plugin-sass",
    buildStart() {
      this.addWatchFile(path.join(__dirname, "public", "styles.sass"))
    },
    generateBundle() {
      const result = sass.compile(path.join(__dirname, "public", "styles.sass"))
      if (!fs.existsSync(path.join(__dirname, "build"))) {
        fs.mkdirSync(path.join(__dirname, "build"))
      }
      fs.writeFileSync(path.join(__dirname, "build", "style.css"), result.css)
    }
  }
}