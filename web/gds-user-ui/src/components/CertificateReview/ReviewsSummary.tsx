import React, { useEffect, useState } from 'react';
import { Stack, HStack, Heading, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import BasicDetailsReview from './BasicDetailsReview';
import ContactsReview from './ContactsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import TrixoReview from './TrixoReview';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep } from 'utils/localStorageHelper';
const ReviewsSummary: React.FC = () => {
  // const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  // const [registrationValues, setRegistrationValue] = React.useState<any>({});
  // useEffect(() => {
  //   const fetchData = async () => {
  //     const getStepperData = await getRegistrationDefaultValue();

  //     setRegistrationValue(getStepperData);
  //   };
  //   fetchData();
  // }, []);

  return (
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
};

export default ReviewsSummary;
