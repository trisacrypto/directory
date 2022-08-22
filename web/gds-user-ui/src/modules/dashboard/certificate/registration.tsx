import React, { useState, useEffect, useCallback } from 'react';
import { SimpleDashboardLayout } from 'layouts';
import {
  Box,
  Heading,
  VStack,
  useToast,
  Text,
  Link,
  Flex,
  useDisclosure,
  Button,
  HStack,
  Stack,
  useColorModeValue,
  chakra
} from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { FormProvider, useForm } from 'react-hook-form';
import Loader from 'components/Loader';
import { userSelector } from 'modules/auth/login/user.slice';
import { yupResolver } from '@hookform/resolvers/yup';
import { DevTool } from '@hookform/devtools';
import { RootStateOrAny, useSelector } from 'react-redux';
import _ from 'lodash';
import { hasStepError, handleError } from 'utils/utils';
import HomeButton from 'components/ui/HomeButton';
import ConfirmationResetFormModal from 'components/Modal/ConfirmationResetFormModal';
import { fieldNamesPerSteps, validationSchema } from './lib';
import { getRegistrationDefaultValues } from 'modules/dashboard/certificate/lib';
import Store from 'application/store';
import {
  postRegistrationValue,
  getRegistrationAndStepperData
} from 'modules/dashboard/registration/utils';

const fieldNamesPerStepsEntries = () => Object.entries(fieldNamesPerSteps);
import { isProdEnv } from 'application/config';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import {
  getCurrentStep,
  getSteps,
  getCurrentState,
  getLastStep,
  getTestNetSubmittedStatus,
  getMainNetSubmittedStatus
} from 'application/store/selectors/stepper';
const Certificate: React.FC = () => {
  const [, updateState] = React.useState<any>();
  const forceUpdate = React.useCallback(() => updateState({}), []);
  const [isResetForm, setIsResetForm] = useState<boolean>(false);
  const textColor = useColorModeValue('black', '#EDF2F7');
  const backgroundColor = useColorModeValue('white', '#171923');

  const { nextStep, previousStep, setInitialState, currentState } = useCertificateStepper();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const currentStep: number = useSelector(getCurrentStep);
  const currentStateValue = useSelector(getCurrentState);
  const lastStep: number = useSelector(getLastStep);
  const steps: number = useSelector(getSteps);
  const isTestNetSubmitted: boolean = useSelector(getTestNetSubmittedStatus);
  const isMainNetSubmitted: boolean = useSelector(getMainNetSubmittedStatus);
  const [isResetModalOpen, setIsResetModalOpen] = useState<boolean>(false);
  const [registrationData, setRegistrationData] = useState<any>([]);
  const [isLoadingDefaultValue, setIsLoadingDefaultValue] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const hasReachSubmitStep: boolean = useSelector(
    (state: RootStateOrAny) => state.stepper.hasReachSubmitStep
  );
  const { isLoggedIn } = useSelector(userSelector);
  const toast = useToast();
  const current = currentStep === lastStep ? lastStep - 1 : currentStep;
  function getCurrentStepValidationSchema() {
    return validationSchema[current - 1];
  }
  const resolver = yupResolver(getCurrentStepValidationSchema());
  // console.log('[registrationData from state]', registrationData);
  const methods = useForm({
    defaultValues: registrationData,
    resolver,
    mode: 'onChange'
  });

  const { formState, reset } = methods;

  function getFieldValue(name: string) {
    return _.get(methods.getValues(), name);
  }

  function isFormCompleted() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.every((n: any) => !!getFieldValue(n));
  }
  // check if the form is submitted or not
  const isFormSubmitted = () => {
    if (isTestNetSubmitted && isMainNetSubmitted) {
      return true;
    }
    if (isTestNetSubmitted || isMainNetSubmitted) {
      return true;
    }
    return false;
  };
  function getCurrentFormValue() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.reduce((acc, n) => ({ ...acc, [n]: getFieldValue(n) }), {});
  }

  function hasErroredField() {
    const fieldsNames = fieldNamesPerStepsEntries()[current - 1][1];
    return fieldsNames.some((n: any) => methods.getFieldState(n).error);
  }

  // if fields if filled

  function handleNextStepClick() {
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
      nextStep({
        isFormCompleted: isFormCompleted(),
        formValues: getCurrentFormValue(),
        values: methods.getValues(),
        registrationValues: registrationData,
        isDirty: methods.formState.isDirty,
        setRegistrationState: setRegistrationData
      });
    }
  }
  const handlePreviousStep = () => {
    previousStep({
      isDirty: methods.formState.isDirty,
      registrationValues: registrationData,
      values: methods.getValues()
    });
  };

  const isDefaultValue = () => {
    return _.isEqual(registrationData, getRegistrationDefaultValues());
  };

  const updateCurrentStep = useCallback(() => {
    const updatedState = Store.getState().stepper;
    return updatedState.currentStep;
  }, []);

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

  const resetForm = () => {
    reset(getRegistrationDefaultValues());
  };

  useEffect(() => {
    resetForm();
    setIsResetForm(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isResetForm, registrationData]);

  // handle reset modal
  useEffect(() => {
    if (isResetModalOpen) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isResetModalOpen]);

  // set registration data value
  useEffect(() => {
    if (registrationData) {
      reset(registrationData);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [registrationData]);

  // load default value from trtl
  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await getRegistrationAndStepperData();
        setRegistrationData(data.registrationData);
        // console.log('[registrationData]', data.registrationData);
        // console.log('[registrationData from state]', data.stepperData);
        console.log('[called from useEffect]');
        setInitialState(data.stepperData);
      } catch (error) {
        handleError(error, 'failed when trying to fetch [getRegistrationAndStepperData]');
      } finally {
        setIsLoadingDefaultValue(false);
      }
    };
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <SimpleDashboardLayout>
      <>
        <FormProvider {...methods}>
          <Flex justifyContent={'space-between'}>
            <Heading size="lg" mb="24px" className="heading">
              <Trans id="Certificate Registration">Certificate Registration</Trans>
            </Heading>
            <Box>{!isLoggedIn && <HomeButton link={'/'} />}</Box>
          </Flex>
          <Stack my={3}>
            <Card maxW="100%" bg={backgroundColor} color={textColor}>
              <Card.Body>
                <Text>
                  <Trans id="This multi-section form is an important step in the registration and certificate issuance process. The information you provide will be used to verify the legal entity that you represent and, where appropriate, will be available to verified TRISA members to facilitate compliance decisions. If you need guidance, see the">
                    This multi-section form is an important step in the registration and certificate
                    issuance process. The information you provide will be used to verify the legal
                    entity that you represent and, where appropriate, will be available to verified
                    TRISA members to facilitate compliance decisions. If you need guidance, see the
                  </Trans>{' '}
                  <Link isExternal href="/getting-started" color={'link'} fontWeight={'bold'}>
                    <Trans id="Getting Started Help Guide">Getting Started Help Guide</Trans>.{' '}
                  </Link>
                </Text>
                <Text pt={4}>
                  <Trans id="To assist in completing the registration form, the form is divided into multiple sections">
                    To assist in completing the registration form, the form is divided into multiple
                    sections
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
          </Stack>

          <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
            <VStack spacing={3}>
              <Box width={'100%'}>
                <TestNetCertificateProgressBar onSetRegistrationState={setRegistrationData} />
                {!isProdEnv ? <DevTool control={methods.control} /> : null}
              </Box>
              <Stack width="100%" direction={'row'} spacing={8} justifyContent={'center'} py={6}>
                {!hasReachSubmitStep && (
                  <>
                    {/* {!isFormSubmitted() && ( */}
                    <Button onClick={handlePreviousStep} isDisabled={currentStep === 1}>
                      <Trans id="Save & Previous">Save & Previous</Trans>
                    </Button>
                    {/* )} */}
                    <Button type="submit" variant="secondary">
                      {currentStep === lastStep ? t`Save & Next` : t`Save & Next`}
                    </Button>
                    {/* add review button when reach to final step */}

                    {/* {!isFormSubmitted() && ( */}
                    <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
                      <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
                    </Button>
                    {/* )} */}
                  </>
                )}
              </Stack>
            </VStack>
          </chakra.form>
        </FormProvider>

        {isResetModalOpen && (
          <ConfirmationResetFormModal
            isOpen={isOpen}
            onClose={onClose}
            onChangeState={onChangeModalState}
            onRefreshState={forceUpdate}
            onReset={reset}
            onChangeResetState={onChangeResetForm}
          />
        )}
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
