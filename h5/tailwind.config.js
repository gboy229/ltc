/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        gold: {
          50: '#fdf9f0',
          100: '#fbf1dd',
          200: '#f0d9ad',
          300: '#e6c384',
          400: '#d9a75d',
          500: '#c99145',
          600: '#b07933',
          700: '#936027',
          800: '#744a20',
          900: '#573719',
        },
        ink: {
          900: '#1f1a13',
          800: '#33291d',
          700: '#55462f',
          600: '#6f5c42',
          500: '#826b4b',
          400: '#a08a6a',
        },
      },
      fontFamily: {
        sans: [
          '"PingFang SC"',
          '"Hiragino Sans GB"',
          '"Microsoft YaHei"',
          'ui-sans-serif',
          'system-ui',
          'sans-serif',
        ],
      },
      boxShadow: {
        'gold-soft': '0 6px 20px rgba(160, 115, 45, 0.08)',
        'gold-btn': '0 6px 16px rgba(166, 116, 44, 0.28)',
      },
    },
  },
  plugins: [],
}
