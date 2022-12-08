import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Table: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    // fontFamily: 'Open Sans, serif'
  },
  sizes: {},
  variants: {
    simple: {
      th: {
        color: '#000',
        textTransform: 'capitalize',
        fontWeight: 700,
        fontSize: 'sm'
      },
      td: {
        fontSize: 14
      }
    }
  },
  defaultProps: {}
};

export default Table;
