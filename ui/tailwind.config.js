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
        ink: "rgba(36, 30, 32, 0.9)",
        light: "rgba(36, 30, 32, 0.7)",
        "gray-std": "rgb(241 245 249)",
      },
    },
  },
  plugins: [
    plugin(function ({ addVariant }) {
      addVariant("focin", ["&:focus", "&:focus-within"])
      addVariant("focgrin", ":merge(.group):focus-within &")
    }),
  ],
}
