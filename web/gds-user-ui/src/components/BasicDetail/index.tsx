import { Box, Heading, Stack, Icon, HStack } from '@chakra-ui/react';
import BasicDetailsForm from 'components/BasicDetailsForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';

type TBasicDetailsProps = {};

const BasicDetails: React.FC<TBasicDetailsProps> = (props) => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  return (
    <Stack
      spacing={5}
      w="100%"
      paddingX="39px"
      paddingY="27px"
      border="3px solid #E5EDF1"
      mt="2rem"
      bg={'white'}
      borderRadius="md">
      <HStack>
        <Heading size="md">Section 1: Basic Details</Heading>{' '}
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <Box w={{ base: '100%', lg: '715px' }}>
        <BasicDetailsForm />
      </Box>
    </Stack>
  );
};

export default BasicDetails;
