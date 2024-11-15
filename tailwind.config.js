/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './web/**/*.{templ,go,js,html}'
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["light", "dark", "cupcake"]
  }
}

