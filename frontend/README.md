# Frontend Setup

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
