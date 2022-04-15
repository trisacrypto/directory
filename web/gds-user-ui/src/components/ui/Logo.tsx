import React from 'react';
import { Box, Text } from '@chakra-ui/react';
import TrisaLogo from 'assets/trisa_logo.svg';

export default function Logo(props: any) {
  return (
    <Box {...props}>
      <img src={TrisaLogo} alt="Trisa logo" />
    </Box>
  );
}
