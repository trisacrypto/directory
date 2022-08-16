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
    wordWrap: 'break-word',
    boxShadow: 'rgb(0 0 0 / 20%) 0px 3px 1px -2px, rgb(0 0 0 / 14%) 0px 2px 2px 0px, rgb(0 0 0 / 12%) 0px 1px 5px 0px',
    _active: {
      boxShadow: 'rgb(0 0 0 / 20%) 0px 5px 5px -3px, rgb(0 0 0 / 14%) 0px 8px 10px 1px, rgb(0 0 0 / 12%) 0px 3px 14px 2px'
    },
    transition: 'background-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, border-color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms, color 250ms cubic-bezier(0.4, 0, 0.2, 1) 0ms'
  },
  sizes: {},
  variants: {
    outline: {
      border: '2px solid blue'
    },
    primary: {
      color: 'white',
      background: 'blue',
      boxShadow: 'rgb(0 0 0 / 20%) 0px 3px 1px -2px, rgb(0 0 0 / 14%) 0px 2px 2px 0px, rgb(0 0 0 / 12%) 0px 1px 5px 0px',
      _focus: {
        boxShadow: 'none'
      },
      _hover: {
        background: '#4389ac',
        _disabled: {
          background: "blue"
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
