import { Tabs, TabList, Tab, TabPanels, TabPanel, Box, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import Metrics from 'components/Metrics';

function NetworkMetricsTabs({ data }: any) {
  return (
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
            <Metrics data={data?.testnet} type="Testnet" />
          </TabPanel>
          <TabPanel p={0}>
            <Metrics data={data?.mainnet} type="Mainnet" />
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Box>
  );
}

export default NetworkMetricsTabs;
