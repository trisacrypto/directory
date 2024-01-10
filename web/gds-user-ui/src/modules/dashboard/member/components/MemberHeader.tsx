import { Heading, HStack, useColorModeValue, Button, Stack } from '@chakra-ui/react';

import { Trans } from '@lingui/macro';
import { useExportMembers } from '../hooks/useExportMembers';

const MemberHeader = () => {
  const { isLoading, exportHandler, isDisabled } = useExportMembers();

  return (
    <Stack width={'100%'}>
      <HStack justify={'space-between'} mb="24px">
        <Heading size="md" color={'black'}>
          <Trans>Member List</Trans>
        </Heading>
        <Button
          isLoading={isLoading}
          disabled={isDisabled}
          minW="100px"
          onClick={exportHandler}
          bg={useColorModeValue('black', 'white')}
          _hover={{
            bg: useColorModeValue('black', 'white')
          }}
          color={useColorModeValue('white', 'black')}>
          <Trans>Export</Trans>
        </Button>
      </HStack>
    </Stack>
  );
};

export default MemberHeader;
