import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
export const Box: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: ({ colorMode }) => ({
    color: colorMode === 'light' ? '#F7F8FC' : '#171923'
  }),
  sizes: {},
  variants: {},
  defaultProps: {}
};
