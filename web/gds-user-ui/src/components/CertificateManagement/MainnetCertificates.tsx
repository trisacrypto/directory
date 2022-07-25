import {
  Stack,
  Box,
  Table,
  TableCaption,
  Heading,
  Button,
  Thead,
  Tr,
  Th,
  Tbody
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import Statistics from './Statistics';
import X509TableRows from './X509TableRows';

function MainnetCertificates() {
  return (
    <Stack spacing={5}>
      <Statistics />
      <Box>
        <FormLayout overflowX={'scroll'}>
          <Table variant="unstyled" css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
            <TableCaption placement="top" textAlign="start" p={0} m={0}>
              <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
                <Heading fontSize={'1.2rem'}>X.509 Identity Certificates</Heading>
                <Button>Request New Identity Certificate</Button>
              </Stack>
            </TableCaption>
            <Thead>
              <Tr>
                <Th>No</Th>
                <Th>Signature ID</Th>
                <Th>Issue Date</Th>
                <Th>Expiration Date</Th>
                <Th>Status</Th>
                <Th textAlign="center">Action</Th>
              </Tr>
            </Thead>
            <Tbody>
              <X509TableRows />
            </Tbody>
          </Table>
        </FormLayout>
      </Box>
    </Stack>
  );
}

export default MainnetCertificates;
