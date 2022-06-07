import customTheme from "../src/utils/theme";
import { LanguageProvider } from '../src/contexts/LanguageContext'

export const parameters = {
  chakra: {
    customTheme,
  },
};

export const decorators = [
  (Story) => (
    <LanguageProvider>
      <Story />
    </LanguageProvider>
  ),
];