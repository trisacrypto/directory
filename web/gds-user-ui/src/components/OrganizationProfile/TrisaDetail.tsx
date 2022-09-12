import React, { useState } from 'react';
import { Heading, Stack, Table, Tbody, Tr, Td, Thead, Tag } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { splitAndDisplay, format2ShortDate } from 'utils/utils';
import { t } from '@lingui/macro';
type TrisaDetailProps = {
  data: any;
};
const TrisaDetail: React.FC<TrisaDetailProps> = ({ data }) => {
  const statusCheck = () => {
    switch (data?.status) {
      case 'NO_VERIFICATION':
        return (
          <Tag bg={'orange'} color={'white'} size={'sm'}>
            {t`${splitAndDisplay(data?.status, '_')}`}
          </Tag>
        );
      case 'VERIFIED':
        return (
          <Tag colorScheme="green" size={'sm'}>
            <Trans id="Verified">VERIFIED</Trans>
          </Tag>
        );
      case 'REJECTED' || 'ERRORED':
        return <Tag colorScheme="red" size={'sm'}>{t`${splitAndDisplay(data?.status, '_')}`}</Tag>;
      default:
        return (
          <Tag colorScheme="yellow" size={'sm'}>{t`${splitAndDisplay(data?.status, '_')}`}</Tag>
        );
    }
  };
  return (
    <Stack
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      fontSize={18}
      p={4}
      mb={10}
      mt={10}
      px={7}>
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pb={7} pt={4}>
          <Trans id="TRISA Details">TRISA Details</Trans>
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
                <Td pt={'1rem !important'}>
                  <Trans id="ID">ID</Trans>
                </Td>
                <Td>
                  <Trans id="First Listed">First Listed</Trans>
                </Td>
                <Td>
                  <Trans id="Verified On">Verified On</Trans>
                </Td>
                <Td>
                  <Trans id="Last Updated">Last Updated</Trans>
                </Td>
                <Td>
                  <Trans id="Status">Status</Trans>
                </Td>
              </Tr>
            </Thead>
            <Tbody
              sx={{
                '*': {
                  fontSize: '1rem',

                  'td:first-child': {},
                  td: {
                    paddingInlineStart: 0.5
                  }
                }
              }}>
              <Tr>
                <Td>{data?.id || 'N/A'}</Td>
                <Td>{data?.first_listed || 'N/A'}</Td>
                <Td>{data?.verified_on ? format2ShortDate(data?.verified_on) : 'N/A'}</Td>
                <Td>{data?.last_updated ? format2ShortDate(data?.last_updated) : 'N/A'}</Td>
                <Td>{data?.status ? statusCheck() : 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Stack>
  );
};
export default TrisaDetail;
