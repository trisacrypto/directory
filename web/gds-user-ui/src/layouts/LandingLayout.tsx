import React from 'react';
import { Flex } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
type LandingLayoutProp = {
  children?: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  return (
    <Flex
      direction="column"
      align="center"
      maxW={'100%'}
      bg={'white'}
      m="0 auto"
      fontFamily={'Open Sans'}
      position={'relative'}
      minHeight={'100vh'}>
      <LandingHeader />
      {props.children}
      <Footer />
    </Flex>
  );
}
