# Frontend Setup

## Tools used (and why)

- **Vue 3**: Used to build the frontend UI in reusable components. It is simple, fast, and good for building screens like login, dashboard, and ticket pages.
- **TypeScript**: Used for type safety. It catches many mistakes early and makes the code easier to maintain as the project grows.
- **Vite**: Used as the dev server and build tool. It starts quickly and gives a fast development experience.
- **Tailwind CSS**: Used for styling. It helps build clean UI faster with utility classes and keeps styling consistent.
- **ESLint (Vue + TypeScript rules)**: Used to enforce code quality and consistency, and to catch common errors before submission.
- **vite-plugin-checker**: Used to run TypeScript/Vue checks during development and show issues early in terminal/overlay.

## Commands used to initialize this app

```bash
# 1) Scaffold Vue + TypeScript app
npm create vite@6.2.0 . -- --template vue-ts

# 2) Install dependencies
npm install

# 3) Install Tailwind packages
npm install -D tailwindcss @tailwindcss/vite

# 4) Install linting + checker tools
npm install -D eslint@8.57.0 @typescript-eslint/parser@6.21.0 @typescript-eslint/eslint-plugin@6.21.0 eslint-plugin-vue@9.32.0 vue-eslint-parser@9.4.3 @types/node vite-plugin-checker@0.11.0
```

## Commands to run

```bash
# Start dev server
npm run dev

# Run linting
npm run lint

# Auto-fix lint issues
npm run lint:fix

# Build production bundle
npm run build
```
