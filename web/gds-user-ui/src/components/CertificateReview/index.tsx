import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';

import BasicDetailsReview from './BasicDetailsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import ContactsReview from './ContactsReview';
import TrixoReview from './TrixoReview';
import FormLayout from 'layouts/FormLayout';
import { RootStateOrAny, useSelector } from 'react-redux';
import useCertificateStepper from 'hooks/useCertificateStepper';

import ReviewSubmit from 'components/ReviewSubmit';
import { ReportHandler } from 'web-vitals';

interface TCertificateReviewProps {}

const CertificateReview = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    await null;
    console.log('handleSubmitRegister', network);
  };

  return (
    <>
      {!hasReachSubmitStep ? (
        <Stack spacing={7}>
          <HStack pt={10}>
            <Heading size="md"> Review </Heading>
            <Box>{/* <Icon as={InfoIcon} color="#F29C36" w={7} h={7} /> (not saved) */}</Box>
          </HStack>
          <FormLayout>
            <Text>
              Please review the information provided, edit as needed, and submit to complete the
              registration form. After the information is reviewed, you will be contacted to verify
              details. Once verified, your TestNet certificate will be issued.
            </Text>
          </FormLayout>
          <BasicDetailsReview />
          <LegalPersonReview />
          <ContactsReview />
          <TrisaImplementationReview />
          <TrixoReview />
        </Stack>
      ) : (
        <ReviewSubmit onSubmitHandler={handleSubmitRegister} />
      )}
    </>
  );
};

export default CertificateReview;
