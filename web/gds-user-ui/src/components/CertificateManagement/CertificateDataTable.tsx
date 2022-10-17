import { Table, TableCaption, Stack, Heading, Thead, Tr, Th, Tbody } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import { BsInfoCircle } from 'react-icons/bs';
import SealingCertificateLearnMore from './SealingCertificateLearnMore';
import SealingCertificateTableRows from './SealingCertificateTableRows';

function CertificateDataTable() {
  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="unstyled" css={{ borderCollapse: 'separate', borderSpacing: '0 9px' }}>
        <TableCaption placement="top" textAlign="start" p={0} m={0}>
          <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
            <Heading fontSize={'1.2rem'}>
              <Trans id="Sealing Certificates">Sealing Certificates</Trans>
            </Heading>
            <SealingCertificateLearnMore>
              <BsInfoCircle style={{ marginRight: 4 }} />
              <Trans id="Learn More">Learn More</Trans>
            </SealingCertificateLearnMore>
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
              <Trans id="Details">Details</Trans>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          <SealingCertificateTableRows />
        </Tbody>
      </Table>
    </FormLayout>
  );
}

export default CertificateDataTable;
