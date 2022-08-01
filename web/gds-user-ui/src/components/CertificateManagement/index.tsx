import DashboardLayout from '../../layouts/DashboardLayout';
import NeedsAttention from '../NeedsAttention';
import {
  Flex,
  Heading,
  Stack,
  Text,
  chakra,
  Box,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel
} from '@chakra-ui/react';
import CertificateDataTable from './CertificateDataTable';
import MainnetCertificates from './MainnetCertificates';
import TestnetCertificates from './TestnetCertificates';

function CertificateManagement() {
  return (
    <DashboardLayout>
      <Heading marginBottom="30px">Certificate Management</Heading>
      <NeedsAttention text={'Complete Certficate Registration'} buttonText={'Continue'} />

      <Flex
        border="1px solid #DFE0EB"
        fontFamily={'Open Sans'}
        bg={'white'}
        fontSize={'1rem'}
        p={5}
        mt={5}
        boxShadow={'0px 24px 50px rgba(55, 65, 81, 0.25)'}
        borderRadius={'10px'}>
        <Stack>
          <Heading fontSize={'1.2rem'}>TRISA Certificates</Heading>
          <Text>TRISA issues two types of certificates:</Text>
          <Text>
            (1) <chakra.span fontWeight={700}>X.509 Identity Certificates:</chakra.span> X.509
            Identity certificates are about <chakra.span fontWeight={700}>trust</chakra.span>. They
            prove the identity of the organization and are used to maintain a secure connection
            between TRISA members. X.509 Identity certificates must have an Endpoint and Common Name
            and may have Subject Alternative Names associated with them. TRISA’s X.509 Identity
            certificates are valid for one calendar year. TRISA’s X.509 Identity certificates expire
            after one year so member organizations must request a new X.509 Identity certificate
            upon expiration.
          </Text>
          <Text>
            (2) <chakra.span fontWeight={700}>Sealing Certificates:</chakra.span> Sealing
            certificates are for specific use cases such as{' '}
            <chakra.span fontWeight={700}>long-term data storage</chakra.span>. They are used to
            encrypt individual Secure Envelopes or batches of Secure Envelopes at the transactional
            level. Organizations may have multiple signing keys and multiple signing-key
            certificates for different clients, time periods, or use cases. In a compliance
            information transfer, the transactional sealing certificates are referred to as Envelope
            Seals. While an organization may use their X.509 Identity certificate as a sealing
            certificate, TRISA strongly recommends that transactional sealing certificates are
            different from X.509 Identity certificates for security purposes.
          </Text>
        </Stack>
      </Flex>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        X.509 Identity Certificate Inventory
      </Heading>
      <Box>
        <Tabs isFitted>
          <TabList bg={'#E5EDF1'} border={'1px solid rgba(0, 0, 0, 0.29)'}>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>MainNet Certificates</Tab>
            <Tab _selected={{ bg: '#60C4CA', fontWeight: 700 }}>TestNet Certificates</Tab>
          </TabList>
          <TabPanels>
            <TabPanel>
              <MainnetCertificates />
            </TabPanel>
            <TabPanel>
              <TestnetCertificates />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
      <Heading fontSize={'1.2rem'} fontWeight={700} p={5} my={5} bg={'#E5EDF1'}>
        Sealing Certificate Inventory
      </Heading>
      <Box>
        <CertificateDataTable />
      </Box>
    </DashboardLayout>
  );
}

export default CertificateManagement;
