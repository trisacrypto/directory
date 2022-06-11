import { Box, Text, Button } from '@chakra-ui/react';

export default function Logo(props: any) {
  return (
    <Box w="80%" pt={7}>
      <Button w="full" colorScheme="red" variant="outline">
        Start trial
      </Button>
    </Box>
  );
}
