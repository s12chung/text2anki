module.exports = {
  env: { browser: true, es2020: true },
  extends: [
    "eslint:all",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/strict",
    "plugin:react-hooks/recommended",
    "plugin:prettier/recommended",
  ],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
    project: ["./tsconfig.json"],
  },
  plugins: ["react-refresh"],
  rules: {
    "@typescript-eslint/no-unsafe-assignment": "off",
    "@typescript-eslint/non-nullable-type-assertion-style": "off",
    "func-style": ["error", "declaration", { allowArrowFunctions: true }],
    "no-magic-numbers": "off",
    "no-shadow": "off",
    "react-refresh/only-export-components": "error",
    "sort-imports": "off",
  },
};
