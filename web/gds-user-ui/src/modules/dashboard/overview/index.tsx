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
import { getMetrics } from './overview.service';
import { useLocation, useNavigate } from 'react-router-dom';
import { colors } from 'utils/theme';
import OverviewLoader from 'components/ContentLoader/Overview';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import OrganizationProfile from 'components/OrganizationProfile';
const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const { user, getUser } = useAuth();
  const navigate = useNavigate();
  useEffect(() => {
    (async () => {
      try {
        const response = await getMetrics();
        console.log('response', response);
        setResult(response.data);
      } catch (e: any) {
        if (e.response.status === 401) {
          navigate('/auth/login?redirect=/dashboard/overview&q=unauthorized');
        }
        if (e.response.status === 403) {
          navigate('/auth/login?redirect=/dashboard/overview&q=token_expired');
        }

        console.log(e);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  return (
    <DashboardLayout>
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
              <Tabs mt={'10'} variant="enclosed">
                <TabList border={'1px solid #eee'} pb={5}>
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
                  <TabPanel p={0} border="1px solid #E5EDF1">
                    <Metrics data={result?.testnet} type="Testnet" />
                  </TabPanel>
                  <TabPanel p={0} border="1px solid #E5EDF1">
                    <Metrics data={result?.mainnet} type="Mainnet" />
                  </TabPanel>
                </TabPanels>
              </Tabs>
            </Box>
          </Box>
          {/* </Sentry.ErrorBoundary> */}
          <OrganizationProfile data={result} />
        </>
      )}
    </DashboardLayout>
  );
};

export default Overview;
