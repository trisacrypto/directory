import { Stack, Box, Table, TableCaption, Heading, Thead, Tr, Th, Tbody } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import ConfirmIdentityCertificate from 'components/ConfirmIdentityCertificateRequest';
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
                <Heading fontSize={'1.2rem'}>
                  <Trans id="X.509 Identity Certificates">X.509 Identity Certificates</Trans>
                </Heading>
                <ConfirmIdentityCertificate marginLeft="auto !important">
                  <Trans id="Request New Identity Certificate">
                    Request New Identity Certificate
                  </Trans>
                </ConfirmIdentityCertificate>
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
      </Box>
    </Stack>
  );
}

export default MainnetCertificates;
