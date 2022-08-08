import React, { useState, useEffect } from 'react';
import * as Sentry from '@sentry/react';
import { Box, Heading } from '@chakra-ui/react';
import OverviewLoader from 'components/ContentLoader/Overview';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import { getMetrics, getAnnouncementsData } from './service';
import { useNavigate } from 'react-router-dom';
import OrganizationalDetail from 'components/OrganizationProfile/OrganizationalDetail';
import TrisaDetail from 'components/OrganizationProfile/TrisaDetail';
import TrisaImplementation from 'components/OrganizationProfile/TrisaImplementation';
import NetworkMetricsTabs from '../../../components/NetworkMetricsTabs';
import { getRegistrationDefaultValue } from '../registration/utils';
import { handleError } from 'utils/utils';

const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const [announcements, setAnnouncements] = React.useState<any>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [stepperData, setStepperData] = React.useState<any>({});
  const [trisaData, setTrisaData] = React.useState<any>({});
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        console.log(1);
        const [metrics, getAnnouncements] = await Promise.all([
          getMetrics(),
          getAnnouncementsData()
        ]);
        console.log('[]', metrics);

        if (metrics.status === 200) {
          setResult(metrics.data);
        }
        if (getAnnouncements.status === 200) {
          setAnnouncements(getAnnouncements.data.announcements);
        }
      } catch (e: any) {
        if (e.response.status === 401) {
          navigate('/auth/login?from=/dashboard/overview&q=unauthorized');
        }
        if (e.response.status === 403) {
          navigate('/auth/login?from=/dashboard/overview&q=token_expired');
        }

        Sentry.captureException(e);
      } finally {
        setIsLoading(false);
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // load legal person & contact information
  useEffect(() => {
    const fetchRegistration = async () => {
      try {
        const registration = await getRegistrationDefaultValue();
        if (registration) {
          setStepperData(registration);
        }
      } catch (e: any) {
        handleError(e, '[Overview] fetchRegistration failed');
      }
    };
    fetchRegistration();

    const trisaDetailData = {
      mainnet: stepperData.mainnet,
      testnet: stepperData.testnet,
      organization: result.organization
    };

    setTrisaData(trisaDetailData);
  }, [result]);

  return (
    <>
      {isLoading ? (
        <OverviewLoader />
      ) : (
        <>
          <Heading marginBottom="30px">Overview</Heading>
          <NeedsAttention
            buttonText={'Start'}
            onClick={() => navigate('/dashboard/certificate/registration')}
          />
          {announcements.length > 0 && <NetworkAnnouncements data={announcements} />}

          <Box fontSize={'md'} mx={'auto'} w={'100%'}>
            <NetworkMetricsTabs data={result} />
          </Box>
          {/* </Sentry.ErrorBoundary> */}
          <OrganizationalDetail data={stepperData} />
          <TrisaDetail data={trisaData} />
          <TrisaImplementation data={trisaData} />
        </>
      )}
    </>
  );
};

export default Overview;
