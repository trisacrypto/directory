import { extendTheme } from '@chakra-ui/react';
import colors from './colors';
import fontSizes from './fontSizes';
import fonts from './fonts';
import breakpoints from './breakpoints';
import * as Components from './components';

import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools';

const config = {
  cssVarPrefix: 'ck',
  initialColorMode: 'light',
  useSystemColorMode: false
};

const theme = extendTheme({
  colors,
  fontSizes,
  fonts,
  breakpoints,
  config,
  components: {
    ...Components
  },
  styles: {
    global: (props: StyleFunctionProps) => ({
      h1: mode('red', 'green')(props)
    })
  }
});

// remove default theme boxShadow
theme.shadows.outline = 'none';

export default theme;
