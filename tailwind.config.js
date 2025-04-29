/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ", // catch all .templ files
    "./**/*.{html,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  // plugins: [require("daisyui")],
  // daisyui: {
  //   themes: ["light", "dark"],
  // },
};
