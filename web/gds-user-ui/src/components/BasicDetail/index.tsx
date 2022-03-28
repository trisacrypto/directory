import { Box, Heading, Stack, Icon, HStack } from '@chakra-ui/react';
import { CheckCircleIcon } from '@chakra-ui/icons';

import BasicDetailsForm from 'components/BasicDetailsForm';

type TBasicDetailsProps = {};

const BasicDetails: React.FC<TBasicDetailsProps> = () => {
  return (
    <Stack
      spacing={5}
      w="100%"
      paddingX="39px"
      paddingY="27px"
      border="3px solid #E5EDF1"
      mt="2rem"
      borderRadius="md">
      <HStack>
        <Heading size="md">Section 1: Basic Details</Heading>{' '}
        <Box>
          <Icon as={CheckCircleIcon} color="green.300" /> (saved)
        </Box>
      </HStack>
      <Box w={{ base: '100%', lg: '715px' }}>
        <BasicDetailsForm />
      </Box>
    </Stack>
  );
};

export default BasicDetails;
