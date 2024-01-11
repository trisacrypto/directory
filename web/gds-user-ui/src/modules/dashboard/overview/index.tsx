import React, { Suspense } from 'react';
import { Box, Heading } from '@chakra-ui/react';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import { useNavigate } from 'react-router-dom';
import { t } from '@lingui/macro';
import MetricsTabs from 'components/MetricsTabs';

import { APP_PATH } from 'utils/constants';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
const Overview: React.FC = () => {
  const { certificateStep, wasCertificateStepFetched } = useFetchCertificateStep({
    key: StepEnum.ALL
  });
  const { setInitialState } = useCertificateStepper();
  const navigate = useNavigate();
  // we need to add a default stepper

  if (wasCertificateStepFetched) {
    if (certificateStep) {
      setInitialState(certificateStep?.form);
    }
  }

  return (
    <>
      <Heading marginBottom="32px">Overview</Heading>
      <Suspense fallback="">
        <NeedsAttention
          text={t`Start Certificate Registration`}
          buttonText={'Start'}
          onClick={() => navigate(APP_PATH.DASH_CERTIFICATE_REGISTRATION)}
        />
      </Suspense>
      <Suspense fallback={t`Failed to load Network announcement , please reload`}>
        <NetworkAnnouncements />
      </Suspense>
      <Box fontSize={'md'} mx={'auto'} w={'100%'}>
        <MetricsTabs />
      </Box>
    </>
  );
};

export default Overview;
