import React from 'react';
import { Heading, Stack, Table, Tbody, Tr, Td, Text, Th } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
type TrisaImplementationProps = {
  data: any;
};
const TrisaImplementation: React.FC<TrisaImplementationProps> = ({ data }) => {
  return (
    <Stack
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      fontSize={18}
      p={4}
      mb={7}
      px={{ base: 4, md: 7 }}>
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pb={2} pt={2}>
          <Trans id="TRISA Implementations">TRISA Implementations</Trans>
        </Heading>
        <Text pb={4}>
          <Trans id="You must request a new X.509 Identity Certificate to change your Endpoint and Common Name.">
            You must request a new X.509 Identity Certificate to change your Endpoint and Common
            Name.
          </Trans>
        </Text>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(1)': { fontWeight: 'bold' },
              Tr: {
                borderStyle: 'hidden'
              }
            }}>
            <Tbody
              sx={{
                '*': {
                  'td:first-child': {
                    width: '50%'
                  },
                  td: {
                    borderBottom: 'none',
                    paddingInlineStart: 0,
                    paddingY: 2.5
                  }
                }
              }}>
              <Tr pt={'1rem !important'}>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" fontSize={'md'} pl={'1rem !important'}>
                  <Trans id="MainNet">MainNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Th pt={'1rem !important'} pl={'1rem !important'}>
                  <Trans id="Endpoint">Endpoint</Trans>
                </Th>
                <Td>{data?.mainnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Th pl={'1rem !important'}>
                  <Trans id="Common Name">Common Name</Trans>
                </Th>
                <Td>{data?.mainnet?.common_name || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td colSpan={2}></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" fontSize={'md'} pl={'1rem !important'}>
                  <Trans id="TestNet">TestNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Th pt={'1rem !important'} pl={'1rem !important'}>
                  <Trans id="Endpoint">Endpoint</Trans>
                </Th>
                <Td>{data?.testnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Th pl={'1rem !important'}>
                  <Trans id="Common Name">Common Name</Trans>
                </Th>
                <Td>{data?.testnet?.common_name || 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default TrisaImplementation;
