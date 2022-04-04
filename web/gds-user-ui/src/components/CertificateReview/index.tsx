import { Box, Heading, HStack, Icon, Stack, Text, useToast } from '@chakra-ui/react';

import BasicDetailsReview from './BasicDetailsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import ContactsReview from './ContactsReview';
import TrixoReview from './TrixoReview';
import FormLayout from 'layouts/FormLayout';
import { RootStateOrAny, useSelector } from 'react-redux';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { hasStepError } from 'utils/utils';
import ReviewSubmit from 'components/ReviewSubmit';
import { registrationRequest } from 'modules/dashboard/certificate/service';
const CertificateReview = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const toast = useToast();
  const steps: number = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    try {
      await registrationRequest(network, steps);
    } catch (err: any) {
      toast({
        title: 'Error',
        description: err.response.message,
        status: 'error',
        duration: 5000,
        isClosable: true
      });
    }
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
