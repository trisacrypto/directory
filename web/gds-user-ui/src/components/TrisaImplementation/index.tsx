import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import TrisaImplementationForm from 'components/TrisaImplementationForm';
import FormLayout from 'layouts/FormLayout';

const TrisaImplementation: React.FC = () => {
  return (
    <Stack spacing={7}>
      <HStack>
        <Heading size="md">Section 4: TRISA Implementation</Heading>
        <Box>
          <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please
          specify the details of your endpoint for certificate issuance.
        </Text>
      </FormLayout>
      <TrisaImplementationForm headerText="TRISA Endpoint: TestNet" />
      <TrisaImplementationForm headerText="TRISA Endpoint: MainNet" />
    </Stack>
  );
};

export default TrisaImplementation;
