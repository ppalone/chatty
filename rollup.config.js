import sass from "./plugins/sass"
import pug from "./plugins/pug"

export default {
  input: "./public/index.js",
  output: {
    dir: "build",
  },
  plugins: [
    sass(),
    pug()
  ]
}