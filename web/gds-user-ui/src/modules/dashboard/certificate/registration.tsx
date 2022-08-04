import React, { useState, useEffect } from 'react';
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
import { hasStepError } from 'utils/utils';
import HomeButton from 'components/ui/HomeButton';
import ConfirmationResetFormModal from 'components/Modal/ConfirmationResetFormModal';
import { fieldNamesPerSteps, validationSchema } from './lib';
import { getRegistrationDefaultValues } from 'modules/dashboard/certificate/lib';
import FileUploader from 'components/FileUpload';
import {
  getRegistrationDefaultValue,
  postRegistrationValue,
  setRegistrationDefaultValue
} from 'modules/dashboard/registration/utils';

const fieldNamesPerStepsEntries = () => Object.entries(fieldNamesPerSteps);
import { isProdEnv } from 'application/config';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';

const Certificate: React.FC = () => {
  const [, updateState] = React.useState<any>();
  const forceUpdate = React.useCallback(() => updateState({}), []);
  const [isResetForm, setIsResetForm] = useState<boolean>(false);
  const [isLoadingDefaultValue, setIsLoadingDefaultValue] = useState(true);
  const textColor = useColorModeValue('black', '#EDF2F7');
  const backgroundColor = useColorModeValue('white', '#171923');

  const { nextStep, previousStep } = useCertificateStepper();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);
  const steps: number = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [isResetModalOpen, setIsResetModalOpen] = useState<boolean>(false);
  const [registrationData, setRegistrationData] = useState<any>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isLoadingRegistration, setIsLoadingRegistration] = useState<boolean>(false);
  const [isLoadingRegistrationDefaultValue, setIsLoadingRegistrationDefaultValue] =
    useState<boolean>(false);
  const [shouldFillForm, setShouldFillForm] = useState<boolean>(false);
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
  console.log('[registrationData from state]', registrationData);
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
      console.log('handleNextStepClick', methods.getValues());
      postRegistrationValue(methods.getValues());
      nextStep({
        isFormCompleted: isFormCompleted(),
        formValues: getCurrentFormValue()
      });
    }
  }
  const handlePreviousStep = () => {
    postRegistrationValue(methods.getValues());
    previousStep();
  };

  const isDefaultValue = () => {
    return _.isEqual(registrationData, getRegistrationDefaultValues());
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
  // handle file upload by extract json file , validate schema  and post to server
  const handleFileUploaded = (file: any) => {
    console.log('[handleFileUploaded]', file);
    setIsLoading(true);
    const reader = new FileReader();
    reader.onload = async (ev: any) => {
      const data = JSON.parse(ev.target.result);
      console.log('[read data]', data);
      try {
        const validationData = await validationSchema[0].validate(data);
        console.log('[validationData]', validationData);
        if (validationData.error) {
          console.log('[validationData.error]', validationData.error);
          toast({
            position: 'top',
            title: `Invalid file format`,
            status: 'error',
            isClosable: true,
            containerStyle: {
              width: '800px',
              maxWidth: '100%'
            }
          });
          setIsLoading(false);
        }
        if (validationData.value) {
          console.log('[validationData.value]', validationData.value);

          postRegistrationValue(validationData.value);
          setIsLoading(false);
          setShouldFillForm(true);
        }

        reader.readAsText(file);
      } catch (e: any) {
        toast({
          position: 'top',
          title: `Invalid file format`,
          description: e.message || 'your json file is invalid',
          status: 'error',
          isClosable: true,
          containerStyle: {
            width: '800px',
            maxWidth: '100%'
          }
        });
        setIsLoading(false);
      }
    };
  };
  // handle reset modal
  useEffect(() => {
    if (isResetModalOpen) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isResetModalOpen]);

  useEffect(() => {
    const defaultValue =
      Object.keys(registrationData).length > 0 ? registrationData : getRegistrationDefaultValues();
    reset(defaultValue);
    setIsResetForm(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [registrationData, isResetForm]);

  // load default value from trtl
  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await getRegistrationDefaultValue();
        console.log('[getRegistrationData]', data);
        setRegistrationData(data);
      } catch (error) {
        console.log('[getRegistrationData]', error);
      } finally {
        setIsLoadingDefaultValue(false);
      }
    };
    fetchData();
  }, []);
  // should choose to fill form or import file when value is default
  useEffect(() => {
    console.log('[isDefaultValue]', isDefaultValue());
    if (!isDefaultValue()) {
      setShouldFillForm(true);
    }
  }, [registrationData]);
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
                <TestNetCertificateProgressBar />
                {!isProdEnv ? <DevTool control={methods.control} /> : null}
              </Box>
              <Stack width="100%" direction={'row'} spacing={8} justifyContent={'center'} py={6}>
                {!hasReachSubmitStep && (
                  <>
                    <Button onClick={handlePreviousStep} isDisabled={currentStep === 1}>
                      <Trans id="Save & Previous">Save & Previous</Trans>
                    </Button>
                    <Button type="submit" variant="secondary">
                      {currentStep === lastStep ? t`Next` : t`Save & Next`}
                    </Button>
                    {/* add review button when reach to final step */}

                    <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
                      <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
                    </Button>
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
            onRefeshState={forceUpdate}
            onChangeResetState={onChangeResetForm}
          />
        )}
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
