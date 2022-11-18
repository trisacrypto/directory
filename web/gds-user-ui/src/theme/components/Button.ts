import type { ComponentStyleConfig } from '@chakra-ui/theme';

const Button: ComponentStyleConfig = {
  baseStyle: {
    textTransform: 'capitalize',
    borderRadius: '5px',
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

export default Button;
