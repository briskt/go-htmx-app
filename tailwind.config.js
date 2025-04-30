/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ", // catch all .templ files
    "./**/*.{html,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
};
