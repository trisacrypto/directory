import { Stack, Table, Tbody, Tr, Td } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

interface TrisaImplementationReviewProps {
  mainnet?: any;
  testnet?: any;
}
function TrisaImplementationReviewDataTable({ mainnet, testnet }: TrisaImplementationReviewProps) {
  return (
    <Stack fontSize={'1rem'}>
      <Table
        sx={{
          'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
          Tr: { borderStyle: 'hidden' }
        }}>
        <Tbody
          sx={{
            'td:first-child': {
              width: '50%',
              paddingLeft: '1rem'
            },
            td: {
              borderBottom: 'none',
              paddingInlineStart: 0,
              paddingY: 2.5
            }
          }}>
          <Tr>
            <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
              <Trans id="TestNet">TestNet</Trans>
            </Td>
          </Tr>
          <Tr>
            <Td pt={'1rem !important'}>
              <Trans id="TestNet TRISA Endpoint">TestNet TRISA Endpoint</Trans>
            </Td>
            <Td pl={0}>{testnet?.endpoint || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="TestNet Certificate Common Name">TestNet Certificate Common Name</Trans>
            </Td>
            <Td pl={0}>{testnet?.common_name || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td colSpan={2}></Td>
          </Tr>
          <Tr>
            <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
              <Trans id="MainNet">MainNet</Trans>
            </Td>
          </Tr>
          <Tr>
            <Td pt={'1rem !important'}>
              <Trans id="MainNet TRISA Endpoint">MainNet TRISA Endpoint</Trans>
            </Td>
            <Td pl={0}>{mainnet?.endpoint || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="MainNet Certificate Common Name">MainNet Certificate Common Name</Trans>
            </Td>
            <Td pl={0}>{mainnet?.common_name || 'N/A'}</Td>
          </Tr>
        </Tbody>
      </Table>
    </Stack>
  );
}

export default TrisaImplementationReviewDataTable;
