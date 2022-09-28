import React, { Suspense } from 'react';
import { Box, Heading } from '@chakra-ui/react';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import { useNavigate } from 'react-router-dom';
import { t } from '@lingui/macro';
import MetricsTabs from 'components/MetricsTabs';
import TrisaOrganizationProfile from 'components/TrisaOrganizationProfile';

const Overview: React.FC = () => {
  const navigate = useNavigate();

  return (
    <>
      <Heading marginBottom="30px">Overview</Heading>
      <Suspense fallback="">
        <NeedsAttention
          text={t`Start Certificate Registration`}
          buttonText={'Start'}
          onClick={() => navigate('/dashboard/certificate/registration')}
        />
      </Suspense>
      <Suspense fallback="">
        <NetworkAnnouncements />
      </Suspense>
      <Box fontSize={'md'} mx={'auto'} w={'100%'}>
        <MetricsTabs />
      </Box>
      {/* </Sentry.ErrorBoundary> */}
      <Suspense fallback="">
        <TrisaOrganizationProfile />
      </Suspense>
    </>
  );
};

export default Overview;
