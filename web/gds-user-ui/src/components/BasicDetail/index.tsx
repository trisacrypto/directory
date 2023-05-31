import React from 'react';
import { Box, Heading, Stack, HStack } from '@chakra-ui/react';
import BasicDetailsForm from 'components/BasicDetailsForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';
import FileUploader from 'components/FileUpload';

import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import MinusLoader from 'components/Loader/MinusLoader';
import { StepEnum } from 'types/enums';
import Store from 'application/store';
import useUploadFile from 'hooks/useUploadFile';
interface BasicDetailProps {
  onChangeRegistrationState?: any;
}
const BasicDetails: React.FC<BasicDetailProps> = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  const { setInitialState } = useCertificateStepper();
  const { isFetchingCertificateStep, certificateStep, wasCertificateStepFetched } =
    useFetchCertificateStep({
      key: StepEnum.BASIC
    });

  const { isFileLoading, handleFileUpload } = useUploadFile();

  if (wasCertificateStepFetched) {
    const { stepper } = Store.getState();
    if (!stepper?.steps) {
      // init stepper
      setInitialState(certificateStep?.form);
    }
  }

  return (
    <Stack spacing={7} mt="2rem">
      <HStack justifyContent={'space-between'}>
        <Box display={'flex'}>
          <Heading size="md" pr={3} ml={2}>
            <Trans id={'Section 1: Basic Details'}>Section 1: Basic Details</Trans>
          </Heading>{' '}
          {stepStatus ? <SectionStatus status={stepStatus} /> : null}
        </Box>
        <Box>
          <FileUploader onReadFileUploaded={handleFileUpload} />
        </Box>
      </HStack>
      <Stack w={{ base: '100%' }}>
        {isFileLoading || isFetchingCertificateStep ? (
          <MinusLoader text={'Loading data ...'} />
        ) : (
          <BasicDetailsForm isLoading={isFetchingCertificateStep} data={certificateStep?.form} />
        )}
      </Stack>
    </Stack>
  );
};

export default BasicDetails;
