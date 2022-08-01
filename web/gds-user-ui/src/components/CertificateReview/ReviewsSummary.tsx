import React, { useEffect, useState } from 'react';
import { Stack, HStack, Heading, Text, Box, Button } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import BasicDetailsReview from './BasicDetailsReview';
import ContactsReview from './ContactsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import TrixoReview from './TrixoReview';

import {
  getRegistrationDefaultValue,
  downloadRegistrationData
} from 'modules/dashboard/registration/utils';
import { handleError } from 'utils/utils';

const ReviewsSummary: React.FC = () => {
  const [isLoadingExport, setIsLoadingExport] = useState(false);
  const handleExport = () => {
    const downladData = async () => {
      try {
        setIsLoadingExport(true);
        await downloadRegistrationData();
      } catch (error) {
        handleError(error, 'Error while downloading registration data');
      } finally {
        setIsLoadingExport(false);
      }
    };
    downladData();
  };

  return (
    <Stack spacing={7}>
      <HStack pt={10} justifyContent={'space-between'}>
        <Heading size="md" data-testid="review">
          <Trans id="Review">Review</Trans>
        </Heading>
        <Box>
          <Button bg={'black'} onClick={handleExport} isLoading={isLoadingExport}>
            Export Data{' '}
          </Button>
        </Box>
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
};

export default ReviewsSummary;
