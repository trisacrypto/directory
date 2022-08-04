import { extendTheme } from '@chakra-ui/react';
import colors from './colors';
import fontSizes from './fontSizes';
import breakpoints from './breakpoints';
import Button from './components/Button';
import Input from './components/Input';
import Select from './components/Select';
import Heading from './components/Heading';

const config = {
  cssVarPrefix: 'ck'
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
    Heading
  }
});

export default theme;
