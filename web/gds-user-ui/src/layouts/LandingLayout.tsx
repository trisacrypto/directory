import React from 'react';
import { Box, Container, Stack } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
type LandingLayoutProp = {
  children?: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  return (
    <Stack
      align="center"
      justifyContent="space-between"
      minW={'100%'}
      m="0 auto"
      spacing={0}
      fontFamily={'Open Sans'}
      position={'relative'}
      minHeight={'100vh'}>
      <LandingHeader />
      <Box maxW="7xl">{props.children}</Box>
      <Footer />
    </Stack>
  );
}
