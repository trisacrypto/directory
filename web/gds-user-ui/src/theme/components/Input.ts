import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Input: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    fontFamily: 'Open Sans, Roboto, sans-serif'
  },
  sizes: {},
  variants: {},
  defaultProps: {
    size: 'md',
    type: 'text'
  }
};

export default Input;
