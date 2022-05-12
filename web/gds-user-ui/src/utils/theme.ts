import { extendTheme, ThemeConfig } from '@chakra-ui/react';
import { mode } from '@chakra-ui/theme-tools';
const config: ThemeConfig = {
  useSystemColorMode: false
};

export const colors = {
  transparent: 'transparent',
  black: '#000',
  white: '#fff',
  font: ['Open Sans', 'sans-serif'],
  gray: {
    50: '#f7fafc',
    900: '#171923'
  },
  system: {
    blue: '#55ACD8',
    cyan: '#1BCE9F',
    orange: '#FF7A59',
    white: '#E3EBEF',
    gray: '#5B5858',
    green: '#0A864F',
    link: '#1F4CED'
  }
};

const fonts = {
  body: 'Open Sans, sans-serif'
};

const customTheme = extendTheme({
  colors,
  config,
  fonts,
  styles: {
    global: (props: any) => ({
      body: {
        bg: mode('#F7F8FC', 'dark')(props)
      }
    })
  }
});

export default customTheme;
