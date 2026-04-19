/** @type {import('tailwindcss').Config} */
export const content = [
  './app/**/*.{ts,tsx}', // Covers app directory (e.g., layout.tsx, page.tsx)
  './src/**/*.{ts,tsx,css}', // Covers src directory, including globals.css
];
export const theme = {
  extend: {},
};
export const plugins = [];