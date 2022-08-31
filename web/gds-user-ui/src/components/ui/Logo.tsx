import React from 'react';
import { Box, BoxProps, Text } from '@chakra-ui/react';
import TrisaLogo from 'assets/trisa_logo.svg';

export default function Logo(props: BoxProps) {
  return (
    <Box {...props}>
      <img src={TrisaLogo} alt="Trisa logo" />
    </Box>
  );
}
