import React from 'react';
import { Box, Stack, Flex } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
type LandingLayoutProp = {
  children?: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  return (
    <Stack
      // align="center"
      justifyContent="space-between"
      minW={'100%'}
      m="0 auto"
      spacing={0}
      fontFamily={'Open Sans'}
      position={'relative'}
      minHeight={'100vh'}>
      <LandingHeader />
      <Stack flexGrow={1}>{props.children}</Stack>
      <Footer />
    </Stack>
  );
}
