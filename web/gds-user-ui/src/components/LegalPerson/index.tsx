import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import CountryOfRegistration from 'components/CountryOfRegistration';
import FormLayout from 'layouts/FormLayout';
import NameIdentifiers from '../NameIdentifiers';
import NationalIdentification from '../NameIdentification';
import Address from 'components/Addresses';
import { useFormContext } from 'react-hook-form';

type LegalPersonProps = {};

const LegalPerson: React.FC<LegalPersonProps> = () => {
  const {} = useFormContext();

  return (
    <Stack spacing={7} mt="2rem">
      <HStack>
        <Heading size="md">Section 2: Legal Person</Heading>
        <Box>
          <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          Please enter the information that identify your organization as a Legal Person. This form
          represents the IVMS 101 data structure for legal persons and is strongly suggested for use
          as KYC or CDD information exchanged in TRISA transfers.
        </Text>
      </FormLayout>
      <NameIdentifiers />
      <Address />
      <CountryOfRegistration />
      <NationalIdentification />
    </Stack>
  );
};

export default LegalPerson;
