import { Stack, HStack, Heading, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import BasicDetailsReview from './BasicDetailsReview';
import ContactsReview from './ContactsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import TrixoReview from './TrixoReview';

const ReviewsSummary: React.FC = () => (
  <Stack spacing={7}>
    <HStack pt={10}>
      <Heading size="md" data-testid="review">
        <Trans id="Review">Review</Trans>
      </Heading>
    </HStack>
    <FormLayout>
      <Text>
        <Trans id="Please review the information provided, edit as needed, and submit to complete the registration form. After the information is reviewed, you will be contacted to verify details. Once verified, your TestNet certificate will be issued.">
          Please review the information provided, edit as needed, and submit to complete the
          registration form. After the information is reviewed, you will be contacted to verify
          details. Once verified, your TestNet certificate will be issued.
        </Trans>
      </Text>
    </FormLayout>
    <BasicDetailsReview />
    <LegalPersonReview />
    <ContactsReview />
    <TrisaImplementationReview />
    <TrixoReview />
  </Stack>
);

export default ReviewsSummary;
