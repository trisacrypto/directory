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
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import X509TableRows from './X509TableRows';

function X509IdentityCertificateInventoryDataTable() {
  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="simple">
        <TableCaption placement="top" textAlign="start" p={0} m={0}>
          <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
            <Heading fontSize={'1.2rem'}>
              <Trans id="X.509 Identity Certificates">X.509 Identity Certificates</Trans>
            </Heading>
            <Button borderRadius={5}>
              <Trans id="Request New Identity Certificate">Request New Identity Certificate</Trans>
            </Button>
          </Stack>
        </TableCaption>
        <Thead>
          <Tr>
            <Th>
              <Trans id="No">No</Trans>
            </Th>
            <Th>
              <Trans id="Signature ID">Signature ID</Trans>
            </Th>
            <Th>
              <Trans id="Issue Date">Issue Date</Trans>
            </Th>
            <Th>
              <Trans id="Expiration Date">Expiration Date</Trans>
            </Th>
            <Th>
              <Trans id="Status">Status</Trans>
            </Th>
            <Th textAlign="center">
              <Trans id="Action">Action</Trans>
            </Th>
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
