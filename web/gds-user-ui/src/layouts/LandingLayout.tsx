import React from 'react';
import { Flex, FlexProps } from '@chakra-ui/react';
import LandingHeader from 'components/Header/LandingHeader';
import LandingFooter from 'components/Footer/LandingFooter';

type LandingLayoutProp = {
  children: React.ReactNode;
};

export default function LandingLayout(props: LandingLayoutProp): JSX.Element {
  return (
    <Flex direction="column" align="center" maxW={'100%'} m="0 auto" fontFamily={'Open Sans'}>
      <LandingHeader />
      {props.children}
      <LandingFooter />
    </Flex>
  );
}
