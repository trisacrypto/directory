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
import LegalPerson from 'components/LegalPerson';

const ReviewsSummary: React.FC = () => {
  const [isLoadingExport, setIsLoadingExport] = useState(false);
  const [basicDetail, setBasicDetail] = React.useState<any>({});
  const [_, setLegalPerson] = React.useState<any>({});
  const [contacts, setContacts] = React.useState<any>({});
  const [trisa, setTrisa] = React.useState<any>({});
  const [trixo, setTrixo] = React.useState<any>({});

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

  // load value from trtl and set it to review pages
  useEffect(() => {
    const fetchData = async () => {
      try {
        const getStepperData = await getRegistrationDefaultValue();
        const basicDetailData = {
          website: getStepperData.website,
          established_on: getStepperData.established_on,
          vasp_categories: getStepperData.vasp_categories,
          business_category: getStepperData.business_category
        };
        const trisaData = {
          mainnet: getStepperData.mainnet,
          testnet: getStepperData.testnet
        };
        setBasicDetail(basicDetailData);
        setLegalPerson(getStepperData.entity);
        setContacts(getStepperData.contacts);
        setTrisa(trisaData);
        setTrixo(getStepperData.trixo);
      } catch (error) {
        handleError(error, 'Error while getting registration data');
      }
    };
    fetchData();
  }, []);

  return (
    <Stack spacing={7}>
      <HStack pt={10} justifyContent={'space-between'}>
        <Heading size="md" data-testid="review">
          <Trans id="Review">Review</Trans>
        </Heading>
        <Box>
          <Button bg={'black'} onClick={handleExport} isLoading={isLoadingExport}>
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
      <BasicDetailsReview data={basicDetail} />
      <LegalPersonReview data={LegalPerson} />
      <ContactsReview data={contacts} />
      <TrisaImplementationReview data={trisa} />
      <TrixoReview data={trixo} />
    </Stack>
  );
};

export default ReviewsSummary;
