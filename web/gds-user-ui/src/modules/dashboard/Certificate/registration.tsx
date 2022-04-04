import { SimpleDashboardLayout } from 'layouts';
import { Box, Heading, HStack, VStack, useToast } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { FormProvider, useForm, useFormState } from 'react-hook-form';

import { yupResolver } from '@hookform/resolvers/yup';
import { DevTool } from '@hookform/devtools';
import { RootStateOrAny, useSelector } from 'react-redux';
import _ from 'lodash';
import { hasStepError } from '../../../utils/utils';

import {
  fieldNamesPerSteps,
  certificateRegistrationValidationSchema,
  getCertificateRegistrationDefaultValue
} from './lib';

const fieldNamesPerStepsEntries = () => Object.entries(fieldNamesPerSteps);

const Certificate: React.FC = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);
  const steps: number = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const toast = useToast();
  const current = currentStep === lastStep ? lastStep - 1 : currentStep;

  const methods = useForm({
    defaultValues: getCertificateRegistrationDefaultValue(),
    resolver: yupResolver(certificateRegistrationValidationSchema),
    mode: 'onChange'
  });

  function getFieldValue(name: string) {
    return _.get(methods.getValues(), name);
  }

  function isFormCompleted() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.every((n) => getFieldValue(n).toString());
  }

  function getCurrentFormValue() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.reduce((acc, n) => ({ ...acc, [n]: getFieldValue(n) }), {});
  }

  function hasErroredField() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.some((n: any) => methods.getFieldState(n).error);
  }

  function handleNextStepClick() {
    if (currentStep === lastStep) {
      if (hasStepError(steps)) {
        toast({
          position: 'top',
          title: `Please fill all the required fields before submitting`,
          status: 'error',
          isClosable: true,
          containerStyle: {
            width: '800px',
            maxWidth: '100%'
          }
        });
      }
    }
    if (hasErroredField()) {
      // i think we should not use alert here , but we need to find a way to display the error message
      // eslint-disable-next-line no-alert
      if (window.confirm('Some requirement elements are missing, Would you like to continue?')) {
        nextStep({
          isFormCompleted: isFormCompleted(),
          errors: methods.formState.errors,
          formValues: getCurrentFormValue()
        });
      }
    } else {
      nextStep({
        isFormCompleted: isFormCompleted(),
        formValues: getCurrentFormValue()
      });
    }
  }
  const handlePreviousStep = () => {
    previousStep();
  };

  return (
    // <DashboardLayout>
    //   <CertificateLayout>
    //     <BasicDetails />
    //   </CertificateLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <>
        <Heading size="lg" mb="24px">
          Certificate Registration
        </Heading>
        <VStack spacing={3}>
          <Card maxW="100%" bg={'white'}>
            <Card.Body>
              This multi-section form is an important step in the registration and certificate
              issuance process. The information you provide will be used to verify the legal entity
              that you represent and, where appropriate, will be available to verified TRISA members
              to facilitate compliance decisions. To assist in completing the registration form, the
              form is broken into multiple sections. No information is sent until you complete
              Section 6 - Review & Submit.
            </Card.Body>
          </Card>
          <Box width={'100%'}>
            <FormProvider {...methods}>
              <TestNetCertificateProgressBar />
              <DevTool control={methods.control} /> {/* setting up the hook form dev tool */}
            </FormProvider>
          </Box>
          {!hasReachSubmitStep && (
            <HStack width="100%" spacing={4} justifyContent={'center'} pt={4}>
              <FormButton
                onClick={handlePreviousStep}
                isDisabled={currentStep === 1}
                borderRadius={5}
                type="button"
                w="100%"
                maxW="13rem">
                Previous
              </FormButton>
              <FormButton
                borderRadius={5}
                w="100%"
                maxW="13rem"
                backgroundColor="#FF7A59"
                _hover={{ backgroundColor: '#f07253' }}>
                Save & Next
              </FormButton>
              <FormButton onClick={handleNextStepClick} borderRadius={5} w="100%" maxW="13rem">
                {currentStep === lastStep ? 'Finish & submit' : 'Save & Continue Later'}
              </FormButton>
            </HStack>
          )}
        </VStack>
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
