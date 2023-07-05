import { Heading, HStack, useColorModeValue, Button } from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
const MemberTableHeader = () => {
  const exportHandler = () => {
    console.log('modalHandler');
  };
  return (
    <HStack justify={'space-between'} mb={'10'}>
      <Heading size="md" color={'black'}>
        <Trans>Member List</Trans>
      </Heading>
      <Button
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
  );
};

export default MemberTableHeader;
