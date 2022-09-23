import fs from "node:fs"
import path from "node:path"

function css() {
  return {
    name: "rollup-plugin-css",
    generateBundle() {
      fs.copyFileSync(
        path.resolve("./public/style.css"),
        path.resolve("./build/style.css")
      )
    }
  }
}

export default {
  input: "./public/index.js",
  output: {
    dir: "build",
  },
  plugins: [
    css(),
  ]
}