module.exports = {
  stories: [
    "../src/**/*.stories.mdx",
    "../src/**/*.stories.@(js|jsx|ts|tsx)"
  ],
  addons: ["@chakra-ui/storybook-addon"],

  framework: "@storybook/react",
  "core": {
    "builder": "webpack5"
  }
}