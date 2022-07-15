import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Button: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    textTransform: 'capitalize',
    borderRadius: 'md',
    fontFamily: 'Open Sans, Roboto, sans-serif',
    paddingX: 16,
    whiteSpace: 'normal',
    wordWrap: 'break-word'
  },
  sizes: {},
  variants: {
    outline: {
      border: '2px solid blue'
    },
    primary: {
      color: 'white',
      background: 'blue'
    },
    secondary: {
      color: 'white',
      background: 'orange'
    }
  },
  defaultProps: {
    size: 'md',
    variant: 'primary',
    type: 'button'
  }
};

export default Button;
