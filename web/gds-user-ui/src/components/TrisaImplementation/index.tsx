import { Heading, HStack, Stack, Text, useToast } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import TrisaImplementationForm from 'components/TrisaImplementationForm';
import FormLayout from 'layouts/FormLayout';
import { useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';

const TrisaImplementation: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);
  const {
    watch,
    trigger,
    register,
    formState: { errors }
  } = useFormContext();
  const toast = useToast();

  const testnetEndpoint = watch('testnet.endpoint');

  useEffect(() => {
    register('tempField', {
      shouldUnregister: true
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    trigger(['testnet.endpoint', 'mainnet.endpoint']);
  }, [testnetEndpoint, trigger]);

  useEffect(() => {
    if (errors?.tempField?.message) {
      toast({
        position: 'top-right',
        title: errors?.tempField?.message,
        status: 'error',
        duration: 5000,
        isClosable: true
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [errors?.tempField?.message]);

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
    </Stack>
  );
};

export default TrisaImplementation;
