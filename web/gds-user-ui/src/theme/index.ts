import { extendTheme } from '@chakra-ui/react';
import colors from './colors';
import fontSizes from './fontSizes';
import breakpoints from './breakpoints';
import * as overridenComponents from './components';
import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools';

const config = {
  cssVarPrefix: 'ck',
  initialColorMode: 'light',
  useSystemColorMode: true
};

const theme = extendTheme({
  colors,
  fontSizes,
  breakpoints,
  config,
  components: {
    ...overridenComponents
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
