import { extendTheme, ThemeConfig } from '@chakra-ui/react';
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
    link: '#1F4CED',
    blue_light: '#D8EAF6'
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
    global: {
      // This could also be "div, p"
      '*, *::before, ::after': {
        wordWrap: 'normal'
      }
      // 'input[type="search"]': {
      //   bg: 'red',
      //   color: 'red'
      // },
      // 'input[type="search"]::-webkit-search-cancel-button': {
      //   WebkitAppearance: 'none',
      //   height: '1em',
      //   width: '1em',
      //   borderRadius: '50em',
      //   background:
      //     'url(https://pro.fontawesome.com/releases/v5.10.0/svgs/solid/times-circle.svg) no-repeat 50% 50%',
      //   backgroundSize: 'contain',
      //   opacity: 0,
      //   pointerEvents: 'none'
      // }
    }
  }
});

export default customTheme;
