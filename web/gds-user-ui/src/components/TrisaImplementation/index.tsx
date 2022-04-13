import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import TrisaImplementationForm from 'components/TrisaImplementationForm';
import FormLayout from 'layouts/FormLayout';

const TrisaImplementation: React.FC = () => {
  return (
    <Stack spacing={7} pt={8}>
      <HStack>
        <Heading size="md">Section 4: TRISA Implementation</Heading>
        <Box>
          <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please
          specify the details of your endpoint for certificate issuance. Please specify the TestNet
          endpoint and the MainNet endpoint. The TestNet endpoint and the MainNet endpoint must be
          different.
        </Text>
      </FormLayout>
      <TrisaImplementationForm
        type="TestNet"
        name="trisa_endpoint_testnet"
        headerText="TRISA Endpoint: TestNet"
      />
      <TrisaImplementationForm
        type="MainNet"
        name="trisa_endpoint_mainnet"
        headerText="TRISA Endpoint: MainNet"
      />
    </Stack>
  );
};

export default TrisaImplementation;
