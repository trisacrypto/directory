import type { ComponentStyleConfig } from '@chakra-ui/theme';

const Modal: ComponentStyleConfig = {
  baseStyle: {
    dialog: {
      border: '1px solid #000',
      borderRadius: 'sm'
    }
  },
  sizes: {},
  variants: {},
  defaultProps: {
    isCentered: true
  }
};

export default Modal;
