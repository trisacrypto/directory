import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Table: ComponentStyleConfig = {
  // The styles all button have in common
  sizes: {},
  variants: {
    simple: {
      th: {
        color: '#000',
        textTransform: 'capitalize',
        fontWeight: 700,
        fontSize: 'sm',
        fontFamily: 'Open Sans, serif',
      },
      td: {
        fontSize: 'sm',
      }
    }
  },
  defaultProps: {}
};

export default Table;
