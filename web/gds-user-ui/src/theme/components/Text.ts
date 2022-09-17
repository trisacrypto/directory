import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
export const Text: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: ({ colorMode }) => ({
    color: colorMode === 'light' ? 'black' : '#EDF2F7'
  }),
  sizes: {},
  variants: {},
  defaultProps: {}
};
