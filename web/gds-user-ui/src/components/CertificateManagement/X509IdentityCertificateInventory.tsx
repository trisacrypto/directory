import {
  Box,
  Button,
  Heading,
  Stack,
  Tab,
  Table,
  TableCaption,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Tbody,
  Th,
  Thead,
  Tr
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from '../../layouts/FormLayout';
import StatisticCard from './StatisticCard';
import X509IdentityCertificateInventoryDataTable from './X509IdentityCertificateInventoryDataTable';
import X509IdentityCertificateInventoryStatistics from './X509IdentityCertificateInventoryStatistics';
import X509TableRows from './X509TableRows';

function X509IdentityCertificateInventory() {
  return (
    <>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        <Trans id="X.509 Identity Certificate Inventory">
          X.509 Identity Certificate Inventory
        </Trans>
      </Heading>
      <Box>
        <Tabs isFitted>
          <TabList bg={'#E5EDF1'} border={'1px solid rgba(0, 0, 0, 0.29)'}>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>
              <Trans id="MainNet Certificates">MainNet Certificates</Trans>
            </Tab>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>
              <Trans id="TestNet Certificates">TestNet Certificates</Trans>
            </Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <Stack spacing={5}>
                <X509IdentityCertificateInventoryStatistics />
                <Box>
                  <X509IdentityCertificateInventoryDataTable />
                </Box>
              </Stack>
            </TabPanel>
            <TabPanel>
              <p>two!</p>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        <Trans id="Sealing Certificate Inventory">Sealing Certificate Inventory</Trans>
      </Heading>
    </>
  );
}

export default X509IdentityCertificateInventory;
