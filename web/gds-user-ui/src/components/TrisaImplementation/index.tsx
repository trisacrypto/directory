import { Heading, HStack, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import FormLayout from 'layouts/FormLayout';
// import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import TrisaForm from 'components/TrisaImplementation/TrisaImplementationForm';

const TrisaImplementation: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  // const toast = useToast();

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
      <TrisaForm />
    </Stack>
  );
};

export default TrisaImplementation;
