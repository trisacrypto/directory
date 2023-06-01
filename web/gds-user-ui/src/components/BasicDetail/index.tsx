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
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';

interface BasicDetailProps {
  onChangeRegistrationState?: any;
}
const BasicDetails: React.FC<BasicDetailProps> = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  const { setInitialState, currentState, nextStep, getIsDirtyStateByStep } =
    useCertificateStepper();
  const { isFetchingCertificateStep, certificateStep, wasCertificateStepFetched } =
    useFetchCertificateStep({
      key: StepEnum.BASIC
    });

  const {
    updateCertificateStep,
    updatedCertificateStep,
    isUpdatingCertificateStep,
    wasCertificateStepUpdated,
    reset
  } = useUpdateCertificateStep();

  const { isFileLoading, handleFileUpload } = useUploadFile();

  const isDirty = getIsDirtyStateByStep(StepEnum.BASIC);

  if (wasCertificateStepFetched) {
    const { stepper } = Store.getState();
    if (!stepper?.steps) {
      // init stepper
      setInitialState(certificateStep?.form);
    }
  }

  if (wasCertificateStepUpdated) {
    reset();
    nextStep(updatedCertificateStep);
  }

  const handleNextStepClick = (values: any) => {
    if (!isDirty) {
      nextStep(certificateStep);
    } else {
      const payload = {
        step: StepEnum.BASIC,
        form: {
          ...values,
          state: currentState()
        } as any
      };
      updateCertificateStep(payload);
    }
  };

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
        {isFileLoading || isFetchingCertificateStep || isUpdatingCertificateStep ? (
          <MinusLoader text={'Loading data ...'} />
        ) : (
          <BasicDetailsForm
            isLoading={isFetchingCertificateStep}
            data={certificateStep?.form}
            onNextStepClick={handleNextStepClick}
          />
        )}
      </Stack>
    </Stack>
  );
};

export default BasicDetails;
