const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  mode: 'jit',
  purge: {
    enabled: process.env.NODE_ENV === 'production',
    // classes that are generated dynamically, e.g. `rounded-${size}` and must
    // be kept
    safeList: [],
    content: [
      './index.html',
      './src/**/*.{vue,js,ts}',
      // etc.
    ],
  },
  theme: {
    extend: {
      fontFamily: {
        // Prefer locally installed Inter; otherwise fall back to system stack
        sans: ['Inter', 'Inter var', ...defaultTheme.fontFamily.sans],
      },
    },
  },
}
