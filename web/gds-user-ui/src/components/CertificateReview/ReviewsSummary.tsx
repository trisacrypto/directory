import React, { useState } from 'react';
import { Stack, HStack, Heading, Text, Box, Button, useColorModeValue } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import BasicDetailsReview from './BasicDetailsReview';
import ContactsReview from './ContactsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import TrixoReview from './TrixoReview';
import { CgExport } from 'react-icons/cg';

import { downloadRegistrationData } from 'modules/dashboard/registration/utils';
import { handleError } from 'utils/utils';

const ReviewsSummary: React.FC = () => {
  const [isLoadingExport, setIsLoadingExport] = useState(false);
  const handleExport = () => {
    const downloadData = async () => {
      try {
        setIsLoadingExport(true);
        await downloadRegistrationData();
      } catch (error) {
        handleError(error, 'Error while downloading registration data');
      } finally {
        setIsLoadingExport(false);
      }
    };
    downloadData();
  };

  return (
    <Stack spacing={7}>
      <HStack pt={10} justifyContent={'space-between'}>
        <Heading size="md" data-testid="review">
          <Trans id="Review">Review</Trans>
        </Heading>
        <Box>
          <Button
            bg={useColorModeValue('black', 'white')}
            _hover={{
              bg: useColorModeValue('black', 'white')
            }}
            color={useColorModeValue('white', 'black')}
            onClick={handleExport}
            isLoading={isLoadingExport}
            leftIcon={<CgExport />}>
            <Trans id="Export Data">Export Data</Trans>
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
