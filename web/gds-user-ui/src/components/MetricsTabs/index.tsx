import React, { Suspense } from 'react';
import { Tabs, TabList, Tab, TabPanels, TabPanel, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import Metrics from 'components/Metrics';
import TrisaDetail from 'components/OrganizationProfile/TrisaDetail';
import { getOverview } from 'modules/dashboard/overview/service';

import { useAsync } from 'react-use';
import { handleError } from 'utils/utils';
import TrisaOrganizationProfile from 'components/TrisaOrganizationProfile';
function MetricsTabs() {
  const { error, value } = useAsync(getOverview);

  if (error) {
    handleError(error);
    return null;
  }

  return (
    <Tabs mt={8}>
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
          <Metrics data={value?.data?.mainnet} type="Mainnet" />
          <TrisaDetail data={value?.data?.mainnet?.member_details} type={'MainNet'} />
          <Suspense fallback="">
            <TrisaOrganizationProfile
              type="mainnet"
              status={value?.data?.mainnet?.member_details.status}
            />
          </Suspense>
        </TabPanel>
        <TabPanel p={0}>
          <Metrics data={value?.data?.testnet} type="Testnet" />
          <TrisaDetail data={value?.data?.testnet?.member_details} type={'TestNet'} />
          <TrisaOrganizationProfile
            type="testnet"
            status={value?.data?.testnet?.member_details.status}
          />
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
}

export default MetricsTabs;
