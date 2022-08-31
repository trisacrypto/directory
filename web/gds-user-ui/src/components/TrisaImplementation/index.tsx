import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import TrisaImplementationForm from 'components/TrisaImplementationForm';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';

const TrisaImplementation: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

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
