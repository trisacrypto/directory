import { Box, Text, Button as CkButton } from '@chakra-ui/react';

export default function Button(props: any) {
  return (
    <Box w="80%" pt={7}>
      <CkButton w="full" colorScheme="red" variant="outline">
        Start trial
      </CkButton>
    </Box>
  );
}
