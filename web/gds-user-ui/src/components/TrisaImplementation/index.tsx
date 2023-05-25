import { Heading, HStack, Stack, Text, chakra } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import TrisaImplementationForm from 'components/TrisaImplementationForm';
import FormLayout from 'layouts/FormLayout';
// import { useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import { trisaImplementationValidationSchema } from 'modules/dashboard/certificate/lib/trisaImplementationValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';
import useCertificateStepper from 'hooks/useCertificateStepper';
import StepButtons from 'components/StepsButtons';

const TrisaImplementation: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  // const toast = useToast();

  const { previousStep, nextStep } = useCertificateStepper();

  const resolver = yupResolver(trisaImplementationValidationSchema);

  const methods = useForm({
    defaultValues: {},
    resolver,
    mode: 'onChange'
  });

  // const {
  //   register,
  //   watch,
  //   trigger,
  //   formState: { errors }
  // } = methods;

  const handleNextStepClick = () => {
    console.log('[] handleNextStep ', methods.getValues());
    nextStep();
  };

  // const testnetEndpoint = watch('testnet.endpoint');

  // useEffect(() => {
  //   register('tempField', {
  //     shouldUnregister: true
  //   });
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, []);

  // useEffect(() => {
  //   trigger(['testnet.endpoint', 'mainnet.endpoint']);
  // }, [testnetEndpoint, trigger]);

  // useEffect(() => {
  //   if (errors?.tempField?.message) {
  //     toast({
  //       position: 'top-right',
  //       title: errors?.tempField?.message as string,
  //       status: 'error',
  //       duration: 5000,
  //       isClosable: true
  //     });
  //   }
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, [errors?.tempField?.message]);

  return (
    <Stack spacing={7} pt={8}>
      <HStack>
        <Heading size="md" pr={3} ml={2}>
          <Trans id="Section 4: TRISA Implementation">Section 4: TRISA Implementation</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please specify the details of your endpoint for certificate issuance. Please specify the TestNet endpoint and the MainNet endpoint. The TestNet endpoint and the MainNet endpoint must be different.">
            Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please
            specify the details of your endpoint for certificate issuance. Please specify the
            TestNet endpoint and the MainNet endpoint. The TestNet endpoint and the MainNet endpoint
            must be different.
          </Trans>
        </Text>
      </FormLayout>

      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
          <TrisaImplementationForm
            type="TestNet"
            name="testnet"
            headerText={t`TRISA Endpoint: TestNet`}
          />
          <TrisaImplementationForm
            type="MainNet"
            name="mainnet"
            headerText={t`TRISA Endpoint: MainNet`}
          />
          <StepButtons handlePreviousStep={previousStep} handleNextStep={handleNextStepClick} />
        </chakra.form>
      </FormProvider>
    </Stack>
  );
};

export default TrisaImplementation;
