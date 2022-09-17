import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
export const Button: ComponentStyleConfig = {
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
      background: 'blue',
      _focus: {
        boxShadow: 'none'
      },
      _hover: {
        background: '#4389ac',
        _disabled: {
          background: 'blue'
        }
      }
    },
    secondary: {
      color: 'white',
      background: 'orange',
      _hover: {
        background: '#c85d42'
      }
    }
  },
  defaultProps: {
    size: 'md',
    variant: 'primary',
    type: 'button'
  }
};
