import React, { useState, useEffect } from 'react';
import { SimpleDashboardLayout } from 'layouts';
import {
  Box,
  Heading,
  HStack,
  VStack,
  useToast,
  Text,
  Link,
  Flex,
  useDisclosure
} from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { FormProvider, useForm, useFormState } from 'react-hook-form';
import useAuth from 'hooks/useAuth';
import { yupResolver } from '@hookform/resolvers/yup';
import { DevTool } from '@hookform/devtools';
import { RootStateOrAny, useSelector } from 'react-redux';
import _ from 'lodash';
import { hasStepError } from 'utils/utils';
import HomeButton from 'components/ui/HomeButton';
import ConfirmationResetFormModal from 'components/Modal/ConfirmationResetFormModal';
import { fieldNamesPerSteps, validationSchema } from './lib';

import {
  loadDefaultValueFromLocalStorage,
  setCertificateFormValueToLocalStorage
} from 'utils/localStorageHelper';
const fieldNamesPerStepsEntries = () => Object.entries(fieldNamesPerSteps);
import { isProdEnv } from 'application/config';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
const Certificate: React.FC = () => {
  const [, updateState] = React.useState<any>();
  const forceUpdate = React.useCallback(() => updateState({}), []);
  const [isResetForm, setIsResetForm] = useState<boolean>(false);

  const { nextStep, previousStep } = useCertificateStepper();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);
  const steps: number = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [isResetModalOpen, setIsResetModalOpen] = useState<boolean>(false);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const { isUserAuthenticated } = useAuth();
  const toast = useToast();
  const current = currentStep === lastStep ? lastStep - 1 : currentStep;
  function getCurrentStepValidationSchema() {
    return validationSchema[current - 1];
  }
  const resolver = yupResolver(getCurrentStepValidationSchema());

  const methods = useForm({
    defaultValues: loadDefaultValueFromLocalStorage(),
    resolver,
    mode: 'onChange'
  });

  const { formState, reset } = methods;

  const dirtyFields = formState.dirtyFields;

  function getFieldValue(name: string) {
    return _.get(methods.getValues(), name);
  }

  function isFormCompleted() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.every((n: any) => !!getFieldValue(n));
  }

  function getCurrentFormValue() {
    // console.log('current', current);
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
          title: `Please fill in all required fields before proceeding`,
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
      if (window.confirm('Some elements required for registration are missing; continue anyway?')) {
        nextStep({
          isFormCompleted: isFormCompleted(),
          errors: methods.formState.errors,
          formValues: getCurrentFormValue()
        });
      }
    } else {
      setCertificateFormValueToLocalStorage(methods.getValues());
      nextStep({
        isFormCompleted: isFormCompleted(),
        formValues: getCurrentFormValue()
      });
    }
  }
  const handlePreviousStep = () => {
    setCertificateFormValueToLocalStorage(methods.getValues());

    previousStep();
  };

  const handleResetForm = () => {
    // open confirmation modal
    setIsResetModalOpen(true);
  };
  const onChangeModalState = (value: boolean) => {
    setIsResetModalOpen(value);
  };
  const onChangeResetForm = (value: boolean) => {
    setIsResetForm(value);
  };

  useEffect(() => {
    if (isResetModalOpen) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isResetModalOpen]);

  useEffect(() => {
    reset({ ...loadDefaultValueFromLocalStorage() });
    setIsResetForm(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isResetForm]);

  return (
    // <DashboardLayout>
    //   <CertificateLayout>
    //     <BasicDetails />
    //   </CertificateLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <>
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(handleNextStepClick)}>
            <Flex justifyContent={'space-between'}>
              <Heading size="lg" mb="24px" className="heading">
                <Trans id="Certificate Registration">Certificate Registration</Trans>
              </Heading>
              <Box>{!isUserAuthenticated && <HomeButton link={'/'} />}</Box>
            </Flex>

            <VStack spacing={3}>
              <Card maxW="100%" bg={'white'}>
                <Card.Body>
                  <Text>
                    <Trans id="This multi-section form is an important step in the registration and certificate issuance process. The information you provide will be used to verify the legal entity that you represent and, where appropriate, will be available to verified TRISA members to facilitate compliance decisions. If you need guidance, see the">
                      This multi-section form is an important step in the registration and
                      certificate issuance process. The information you provide will be used to
                      verify the legal entity that you represent and, where appropriate, will be
                      available to verified TRISA members to facilitate compliance decisions. If you
                      need guidance, see the
                    </Trans>{' '}
                    <Link isExternal href="/getting-started" color={'blue'} fontWeight={'bold'}>
                      <Trans id="Getting Started Help Guide">Getting Started Help Guide</Trans>.{' '}
                    </Link>
                  </Text>
                  <Text pt={4}>
                    <Trans id="To assist in completing the registration form, the form is divided into multiple sections">
                      To assist in completing the registration form, the form is divided into
                      multiple sections
                    </Trans>
                    .{' '}
                    <Text as={'span'} fontWeight={'bold'}>
                      <Trans id="No information is sent until you complete Section 6 - Review & Submit">
                        No information is sent until you complete Section 6 - Review & Submit
                      </Trans>
                      .{' '}
                    </Text>
                  </Text>
                </Card.Body>
              </Card>

              <Box width={'100%'}>
                <TestNetCertificateProgressBar />
                {!isProdEnv ? <DevTool control={methods.control} /> : null}
              </Box>
              <HStack width="100%" spacing={8} justifyContent={'center'} py={6}>
                {!hasReachSubmitStep && (
                  <>
                    <FormButton
                      onClick={handlePreviousStep}
                      isDisabled={currentStep === 1}
                      borderRadius={5}
                      type="button"
                      w="100%"
                      _active={{ outline: 'none' }}
                      _focus={{ outline: 'none' }}
                      maxW="13rem">
                      <Trans id="Save & Previous">Save & Previous</Trans>
                    </FormButton>
                    <FormButton
                      borderRadius={5}
                      w="100%"
                      _active={{ outline: 'none' }}
                      _focus={{ outline: 'none' }}
                      maxW="13rem"
                      backgroundColor="#FF7A59"
                      type="submit"
                      _hover={{ backgroundColor: '#f07253' }}>
                      {currentStep === lastStep ? t`Next` : t`Save & Next`}
                    </FormButton>
                    {/* add review button when reach to final step */}

                    <FormButton
                      borderRadius={5}
                      w="100%"
                      maxW="13rem"
                      onClick={handleResetForm}
                      type="button">
                      <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
                    </FormButton>
                  </>
                )}
              </HStack>
            </VStack>
          </form>
        </FormProvider>
        {isResetModalOpen && (
          <ConfirmationResetFormModal
            isOpen={isOpen}
            onClose={onClose}
            onChangeState={onChangeModalState}
            onRefeshState={forceUpdate}
            onChangeResetState={onChangeResetForm}
          />
        )}
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
