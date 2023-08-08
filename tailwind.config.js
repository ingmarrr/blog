/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "*.{html, go}",
    "view/*.{html, js, go}",
    "pkg/*.{html, js, go}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Molengo', 'monospace', 'sans-serif']
      },
    },
  },
  plugins: [],
}
