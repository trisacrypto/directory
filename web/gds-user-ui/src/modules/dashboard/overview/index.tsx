import React, { useState, useEffect } from 'react';
import * as Sentry from '@sentry/react';
import {
  Box,
  Heading,
  VStack,
  Flex,
  Input,
  Stack,
  Text,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel
} from '@chakra-ui/react';
import Card from 'components/ui/Card';
import DashboardLayout from 'layouts/DashboardLayout';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import Metrics from 'components/Metrics';
import useAuth from 'hooks/useAuth';
import { getMetrics, getAnnouncementsData } from './service';
import { useLocation, useNavigate } from 'react-router-dom';
import { colors } from 'utils/theme';
import OverviewLoader from 'components/ContentLoader/Overview';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import OrganizationalDetail from 'components/OrganizationProfile/OrganizationalDetail';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import TrisaDetail from 'components/OrganizationProfile/TrisaDetail';
import TrisaImplementation from 'components/OrganizationProfile/TrisaImplementation';
const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const [announcements, setAnnouncements] = React.useState<any>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  // const { user, getUser } = useAuth();
  const [stepperData, setStepperData] = React.useState<any>({});
  const [trisaData, setTrisaData] = React.useState<any>({});

  const navigate = useNavigate();
  useEffect(() => {
    (async () => {
      try {
        const [metrics, getAnnouncements] = await Promise.all([
          getMetrics(),
          getAnnouncementsData()
        ]);
        if (metrics.status === 200) {
          setResult(metrics);
        }
        if (getAnnouncements.status === 200) {
          setAnnouncements(getAnnouncements.data.announcements);
        }

        console.log('[Overview] metrics', metrics);
        console.log('[announcements]', getAnnouncements);
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
    const getStepperData = loadDefaultValueFromLocalStorage();
    const trisaDetailData = {
      mainnet: getStepperData.trisa_endpoint_mainnet,
      testnet: getStepperData.trisa_endpoint_testnet,
      organization: result.organization
    };

    setTrisaData(trisaDetailData);

    setStepperData(getStepperData);
  }, [result]);

  return (
    <DashboardLayout>
      <Heading marginBottom="30px">Overview</Heading>
      <NeedsAttention />
      {announcements.length === 0 && <NetworkAnnouncements datas={announcements} />}
      {/* <Sentry.ErrorBoundary
      <Heading marginBottom="69px">Overview</Heading>
      {isLoading ? (
        <OverviewLoader />
      ) : (
        <>
          <NeedsAttention />
          <NetworkAnnouncements />
          {/* <Sentry.ErrorBoundary
        fallback={<Text color={'red'}>An error has occurred to load testnet metric</Text>}> */}
      {/* <Metrics data={result?.testnet} type="Testnet" />
      <Metrics data={result?.mainnet} type="Mainnet" /> */}
      <Box fontSize={'md'} mx={'auto'} w={'100%'}>
        <Box>
          <Tabs my={'10'}>
            <TabList>
              <Tab
                bg={'#E5EDF1'}
                sx={{ width: '100%' }}
                _focus={{ outline: 'none' }}
                _selected={{ bg: '#60C4CA', color: 'white', fontWeight: 'semibold' }}>
                <Text fontSize={['x-small', 'medium']}>
                  <Trans id="MainNet Network Metrics">MainNet Network Metrics</Trans>
                </Text>
              </Tab>
              <Tab
                bg={'#E5EDF1'}
                fontWeight={'bold'}
                sx={{ width: '100%' }}
                _focus={{ outline: 'none' }}
                _selected={{ bg: '#60C4CA', color: 'white', fontWeight: 'semibold' }}>
                <Text fontSize={['x-small', 'medium']}>
                  <Trans id="TestNet Network Metrics">TestNet Network Metrics</Trans>
                </Text>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel p={0}>
                <Metrics data={result?.testnet} type="Testnet" />
              </TabPanel>
              <TabPanel p={0}>
                <Metrics data={result?.mainnet} type="Mainnet" />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Box>
      </Box>
      {/* </Sentry.ErrorBoundary> */}
      <OrganizationalDetail data={stepperData} />
      <TrisaDetail data={trisaData} />
      <TrisaImplementation data={trisaData} />
    </DashboardLayout>
  );
};

export default Overview;
