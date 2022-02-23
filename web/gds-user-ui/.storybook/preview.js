import customTheme from '../src/utils/theme';

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  chakra: {
    customTheme,
  },
}