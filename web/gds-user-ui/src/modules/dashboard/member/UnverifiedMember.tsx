import React from 'react';
import { Box, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';

const UnverifiedMember = () => {
  return (
    <Box width={'100%'} bg={'#F7F8FC'} p="8" mx={'auto'} textAlign={'center'}>
      <Text fontSize="md" mb={2}>
        <Trans>
          Network directory member list not available because you are not a verified contact for
          this network. Please complete the registration process.
        </Trans>
      </Text>
    </Box>
  );
};

export default React.memo(UnverifiedMember);
