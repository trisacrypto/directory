import React, { useState, useEffect } from 'react';
import * as Sentry from '@sentry/react';
import { Box, Heading, Text, Tabs, TabList, TabPanels, Tab, TabPanel } from '@chakra-ui/react';
import OverviewLoader from 'components/ContentLoader/Overview';
import DashboardLayout from 'layouts/DashboardLayout';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import Metrics from 'components/Metrics';
import { getMetrics, getAnnouncementsData } from './service';
import { useNavigate } from 'react-router-dom';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import OrganizationalDetail from 'components/OrganizationProfile/OrganizationalDetail';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import TrisaDetail from 'components/OrganizationProfile/TrisaDetail';
import TrisaImplementation from 'components/OrganizationProfile/TrisaImplementation';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { handleError } from 'utils/utils';
import useFetchAttention from 'hooks/useFetchAttention';
const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const [announcements, setAnnouncements] = React.useState<any>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [stepperData, setStepperData] = React.useState<any>({});
  const [trisaData, setTrisaData] = React.useState<any>({});
  const { attentionResponse, attentionError, attentionLoading } = useFetchAttention();
  const navigate = useNavigate();

  // console.log('[attentionResponse]', attentionResponse);
  useEffect(() => {
    (async () => {
      try {
        const [metrics, getAnnouncements] = await Promise.all([
          getMetrics(),
          getAnnouncementsData()
        ]);
        if (metrics.status === 200) {
          setResult(metrics.data);
        }
        if (getAnnouncements.status === 200) {
          setAnnouncements(getAnnouncements.data.announcements);
        }
      } catch (e: any) {
        handleError(e);
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
          {attentionResponse && Object.keys(attentionResponse).length > 0 && (
            <NeedsAttention
              loading={attentionLoading}
              error={attentionError}
              data={attentionResponse.messages}
              text={t`Start Certificate Registration`}
              buttonText={'Start'}
              onClick={() => navigate('/dashboard/certificate/registration')}
            />
          )}

          {announcements.length > 0 && <NetworkAnnouncements datas={announcements} />}

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
                    <Metrics data={result?.mainnet} type="Mainnet" />

                    <TrisaDetail data={result?.mainnet?.member_details} type={'MainNet'} />
                  </TabPanel>
                  <TabPanel p={0}>
                    <Metrics data={result?.testnet} type="Testnet" />

                    <TrisaDetail data={result?.testnet?.member_details} type={'TestNet'} />
                  </TabPanel>
                </TabPanels>
              </Tabs>
            </Box>
          </Box>
          {/* </Sentry.ErrorBoundary> */}
          <OrganizationalDetail data={stepperData} />
          <TrisaImplementation data={trisaData} />
        </>
      )}
    </>
  );
};

export default Overview;
