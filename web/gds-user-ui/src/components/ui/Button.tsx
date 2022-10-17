import { Box, Button as CkButton } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

export default function Button() {
  return (
    <Box w="80%" pt={7}>
      <CkButton w="full" colorScheme="red" variant="outline">
        <Trans id="Start trial">Start trial</Trans>
      </CkButton>
    </Box>
  );
}
