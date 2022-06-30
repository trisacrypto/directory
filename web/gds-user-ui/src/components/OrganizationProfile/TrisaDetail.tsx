import React, { useState } from 'react';
import { Heading, Stack, Table, Tbody, Tr, Td, Thead } from '@chakra-ui/react';
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
      mb={10}
      px={7}>
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pb={7} pt={4}>
          {' '}
          TRISA Details{' '}
        </Heading>
        <Stack fontSize={18}>
          <Table
            sx={{
              thead: { fontWeight: 'bold' },
              Tr: {
                borderStyle: 'hidden'
              }
            }}>
            <Thead
              sx={{
                td: {
                  paddingInlineStart: 0.5
                }
              }}>
              <Tr>
                <Td pt={'1rem !important'}>ID</Td>
                <Td>First Listed</Td>
                <Td>Verified On</Td>
                <Td>Last Updated</Td>
              </Tr>
            </Thead>
            <Tbody
              sx={{
                '*': {
                  fontSize: '1rem',

                  'td:first-child': {},
                  td: {
                    paddingInlineStart: 0.5,
                    width: '20%'
                  }
                }
              }}>
              <Tr>
                <Td>{data?.organization?.vasp_id || 'N/A'}</Td>
                <Td>{data?.organization?.first_listed || 'N/A'}</Td>
                <Td>{data?.organization?.verified_on || 'N/A'}</Td>
                <Td>{data?.organization?.last_updated || 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default TrisaDetail;
