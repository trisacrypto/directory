import React from 'react';
import { Box, Text } from '@chakra-ui/react';

export default function Logo(props: any) {
  return (
    <Box {...props}>
      <Text fontSize="lg" fontWeight="bold">
        Trisa
      </Text>
    </Box>
  );
}
