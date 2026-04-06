module.exports = {
  root: true,
  env: {
    browser: true,
    es2022: true,
    node: true,
  },
  parser: "vue-eslint-parser",
  parserOptions: {
    parser: "@typescript-eslint/parser",
    ecmaVersion: "latest",
    sourceType: "module",
    extraFileExtensions: [".vue"],
  },
  extends: [
    "eslint:recommended",
    "plugin:vue/vue3-recommended",
    "plugin:@typescript-eslint/recommended",
  ],
  rules: {
    "vue/multi-word-component-names": "off",
  },
  overrides: [
    {
      files: ["src/views/tickets/TicketDetailView.vue"],
      rules: {
        // VULN-03 (vulnerable baseline): intentional v-html for stored-XSS coursework; re-enable on secure branch.
        "vue/no-v-html": "off",
      },
    },
  ],
};
