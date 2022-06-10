import React, { useState } from 'react';
import { Heading, Stack, Table, Tbody, Tr, Td } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
type TrisaDetailProps = {
  data: any;
};
const TrisaDetail: React.FC<TrisaDetailProps> = ({ data }) => {
  return (
    <Stack
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      fontSize={18}
      p={4}
      px={7}>
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pb={7} pt={4}>
          {' '}
          TRISA Details{' '}
        </Heading>
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
                  fontSize: '1rem',
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
              <Tr>
                <Td pt={'1rem !important'}>ID</Td>
                <Td borderBottom={'none'} whiteSpace="break-spaces" lineHeight={1.5}>
                  {data?.organization?.vasp_id || 'N/A'}
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>
                  <Trans id="First Listed">First Listed</Trans>
                </Td>
                <Td>{data?.organization?.first_listed || 'N/A'}</Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pt={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="Verified On">Verified On</Trans>
                </Td>
                <Td>{data?.organization?.verified_on || 'N/A'}</Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pt={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="Last Updated">Last Updated</Trans>
                </Td>
                <Td pb={'2rem !important'}>{data?.organization?.last_updated || 'N/A'}</Td>
              </Tr>
              <Tr pt={'1rem !important'}>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  <Trans id="TestNet">TestNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'} pl={'1rem !important'}>
                  <Trans id="Endpoint">Endpoint</Trans>
                </Td>
                <Td>{data?.testnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td pl={'1rem !important'}>
                  <Trans id="Common Name">Common Name</Trans>
                </Td>
                <Td>{data?.testnet?.common_name || 'N/A'}</Td>
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
                <Td pt={'1rem !important'} pl={'1rem !important'}>
                  <Trans id="Endpoint">Endpoint</Trans>
                  Endpoint
                </Td>
                <Td>{data?.mainnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td pl={'1rem !important'}>
                  <Trans id="Common Name">Common Name</Trans>
                </Td>
                <Td>{data?.mainnet?.common_name || 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default TrisaDetail;
