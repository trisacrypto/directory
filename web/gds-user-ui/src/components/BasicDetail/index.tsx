import React, { useState } from 'react';
import { Box, Heading, Stack, HStack, useToast } from '@chakra-ui/react';
import BasicDetailsForm from 'components/BasicDetailsForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus, handleError, format2ShortDate } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';
import FileUploader from 'components/FileUpload';
import { validationSchema } from 'modules/dashboard/certificate/lib';
import { postRegistrationValue } from 'modules/dashboard/registration/utils';
import { getRegistrationData } from 'modules/dashboard/registration/service';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import MinusLoader from 'components/Loader/MinusLoader';
import { StepEnum } from 'types/enums';
interface BasicDetailProps {
  onChangeRegistrationState?: any;
}
const BasicDetails: React.FC<BasicDetailProps> = ({ onChangeRegistrationState }) => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);
  const toast = useToast();

  const { updateStateFromFormValues, setRegistrationValue } = useCertificateStepper();
  const { isFetchingCertificateStep, certificateStep } = useFetchCertificateStep({
    key: StepEnum.BASIC
  });

  const [isLoadingDefaultValue, setIsLoadingDefaultValue] = useState(false);
  // create a hook to handle file uploaded
  const handleFileUploaded = (file: any) => {
    // console.log('[handleFileUploaded]', file);
    setIsLoadingDefaultValue(true);
    const reader = new FileReader();
    reader.onload = async (ev: any) => {
      // if file is empty
      // if (!ev.target.result) {
      //   setIsLoadingDefaultValue(false);
      //   return;
      // }

      const data = JSON.parse(ev.target.result);
      try {
        const validationData = await validationSchema[0].validate(data);
        const updatedCertificate: any = await postRegistrationValue(validationData);
        // console.log('[updatedCertificate]', updatedCertificate);
        // console.log('[updatedCertificate Status]', updatedCertificate.status);
        if (updatedCertificate.status === 200) {
          const getValue = await getRegistrationData();
          const values = {
            ...getValue.data,
            established_on: getValue?.data?.established_on
              ? format2ShortDate(getValue?.data?.established_on)
              : ''
          };
          onChangeRegistrationState(values);
          setRegistrationValue(values);
          updateStateFromFormValues(values.state);
        }
      } catch (e: any) {
        if (e.name === 'ValidationError') {
          toast({
            title: 'Invalid file',
            description: e.message || 'your json file is invalid',
            status: 'error',
            duration: 5000,
            isClosable: true,
            position: 'top-right'
          });
          handleError(e, `[Invalid file], it's missing some required fields : ${e.message}`);
        }
      } finally {
        setIsLoadingDefaultValue(false);
      }
    };

    reader.readAsText(file);
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
          <FileUploader onReadFileUploaded={handleFileUploaded} />
        </Box>
      </HStack>
      <Stack w={{ base: '100%' }}>
        {isLoadingDefaultValue || isFetchingCertificateStep ? (
          <MinusLoader text={'Loading data ...'} />
        ) : (
          <BasicDetailsForm isLoading={isFetchingCertificateStep} data={certificateStep?.form} />
        )}
      </Stack>
    </Stack>
  );
};

export default BasicDetails;
