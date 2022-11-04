import { extendTheme } from '@chakra-ui/react';
import colors from './colors';
import fontSizes from './fontSizes';
import breakpoints from './breakpoints';
import Button from './components/Button';
import Input from './components/Input';
import Select from './components/Select';
import Heading from './components/Heading';
import { mode, StyleFunctionProps } from '@chakra-ui/theme-tools';
import Table from './components/Table';

const config = {
  cssVarPrefix: 'ck',
  initialColorMode: 'light',
  useSystemColorMode: false
};

const theme = extendTheme({
  colors,
  fontSizes,
  breakpoints,
  config,
  components: {
    Button,
    Input,
    Select,
    Heading,
    Table
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
