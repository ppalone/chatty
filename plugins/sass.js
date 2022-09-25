import path from "node:path"
import fs from "node:fs"
import sass from "sass"

export default function plugin() {
  return {
    name: "rollup-plugin-sass",
    buildStart() {
      this.addWatchFile(path.resolve(__dirname, "public", "styles.sass"))
    },
    async generateBundle() {
      try {
        const result = await sass.compileAsync(path.resolve(__dirname, "public", "styles.sass"))
        fs.writeFileSync(path.join(__dirname, "build", "style.css"), result.css)
      } catch (err) {
        console.error(err)
      }
    }
  }
}