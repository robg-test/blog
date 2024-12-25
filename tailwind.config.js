/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './web/**/*.{templ,go,js,html}'
  ],
  plugins: [
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
  daisyui: {
    themes: ["retro", "dark", "retro"]
  }
}

