import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import ContactForm from 'components/ContactForm';
import FormLayout from 'layouts/FormLayout';

const Contacts: React.FC = () => {
  return (
    <Stack spacing={13} mt="2rem">
      <HStack>
        <Heading size="md">Section 3: Contacts</Heading>
        <Box>
          <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved)
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          Please supply contact information for representatives of your organization. All contacts
          will receive an email verification token and the contact email must be verified before the
          registration can proceed.
        </Text>
      </FormLayout>
      <ContactForm
        name="contacts.legal"
        title="Legal/ Compliance Contact (required)"
        description="Compliance officer or legal contact for requests about the compliance requirements and legal status of your organization."
      />
      <ContactForm
        name="contacts.technical"
        title="Technical Contact (required)"
        description="Primary contact for handling technical queries about the operation and status of your service participating in the TRISA network. Can be a group or admin email."
      />
      <ContactForm
        name="contacts.administrative"
        title="Administrative Contact (optional)"
        description="Administrative or executive contact for your organization to field high-level requests or queries. (Strongly recommended)"
      />
      <ContactForm
        name="contacts.billing"
        title="Billing Contact (optional)"
        description="Billing contact for your organization to handle account and invoice requests or queries relating to the operation of the TRISA network."
      />
    </Stack>
  );
};

export default Contacts;
