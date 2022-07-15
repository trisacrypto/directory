import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Select: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    fontFamily: 'Open Sans, Roboto, sans-serif',
    border: 'none'
  },
  sizes: {},
  variants: {},
  defaultProps: {
    size: 'md',
    type: 'text'
  }
};

export default Select;
