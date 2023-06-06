/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: 'source-han-sans-cjk-ko,"Open Sans","Gill Sans MT","Gill Sans",Corbel,Arial,sans-serif',
        'ko-sans': 'gowun-batang,"Open Sans","Gill Sans MT","Gill Sans",Corbel,Arial,sans-serif'
      },
      colors: {
        'ink': 'rgba(36, 30, 32, 0.9)'
      },
      maxWidth: {
        '8xl': '90rem'
      }
    },
  },
  plugins: [],
}
