import type { ComponentStyleConfig } from '@chakra-ui/theme';

// You can also use the more specific type for
// a single part component: ComponentSingleStyleConfig
const Heading: ComponentStyleConfig = {
  // The styles all button have in common
  baseStyle: {
    fontFamily: 'Roboto Slab, serif'
  },
  sizes: {},
  variants: {},
  defaultProps: {}
};

export default Heading;
