import {
  Table,
  TableCaption,
  Stack,
  Heading,
  Button,
  Thead,
  Tr,
  Th,
  Tbody
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import X509TableRows from './X509TableRows';

function X509IdentityCertificateInventoryDataTable() {
  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="unstyled" css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
        <TableCaption placement="top" textAlign="start" p={0} m={0}>
          <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
            <Heading fontSize={'1.2rem'}>X.509 Identity Certificates</Heading>
            <Button borderRadius={5}>Request New Identity Certificate</Button>
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
  );
}

export default X509IdentityCertificateInventoryDataTable;
