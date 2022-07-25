import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Input: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    field: {
      fontFamily: 'Open Sans, Roboto, sans-serif'
    }
  },
  sizes: {
    lg: {
      field: {
        borderRadius: 0,
        backgroundColor: '#E5EDF1 !important'
      }
    }
  },
  variants: {},
  defaultProps: {
    size: 'md',
    type: 'text'
  }
};

export default Input;
