import { Box, Heading, Stack, Icon, HStack } from '@chakra-ui/react';
import BasicDetailsForm from 'components/BasicDetailsForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';

const BasicDetails: React.FC = (props) => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  return (
    <Stack spacing={7} mt="2rem">
      <HStack>
        <Heading size="md">
          <Trans id={'Section 1: Basic Details'}>Section 1: Basic Details</Trans>
        </Heading>{' '}
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <Box w={{ base: '100%' }}>
        <BasicDetailsForm />
      </Box>
    </Stack>
  );
};

export default BasicDetails;
