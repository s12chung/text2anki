/** @type {import('tailwindcss').Config} */
const plugin = require("tailwindcss/plugin")

export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: 'source-han-sans-cjk-ko,"Open Sans","Gill Sans MT","Gill Sans",Corbel,Arial,sans-serif',
      },
      colors: {
        ink: "rgb(36, 30, 32)",
        light: "rgba(36, 30, 32, 0.7)",
        faded: "rgba(36, 30, 32, 0.3)",
        "gray-std": "rgb(241 245 249)",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    plugin(function ({ addVariant }) {
      addVariant("child", "& > *")
    }),
  ],
}
