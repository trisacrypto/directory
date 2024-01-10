import React from 'react';
import { Heading, Stack, Table, Tbody, Tr, Td, Thead, Tag, Th } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { splitAndDisplay, format2ShortDate } from 'utils/utils';
import { t } from '@lingui/macro';
type TrisaDetailProps = {
  data: any;
  type?: string;
};
const TrisaDetail: React.FC<TrisaDetailProps> = ({ data, type }) => {
  const statusCheck = () => {
    switch (data?.status) {
      case 'NO_VERIFICATION':
        return (
          <Tag bg={'orange'} color={'white'}>
            <Trans id="Not Verified">Not Verified</Trans>
          </Tag>
        );
      case 'VERIFIED':
        return (
          <Tag colorScheme="green">
            <Trans id="Verified">Verified</Trans>
          </Tag>
        );
      case 'REJECTED' || 'ERRORED':
        return <Tag colorScheme="red">{t`${splitAndDisplay(data?.status, '_')}`}</Tag>;
      default:
        return (
          <Tag colorScheme="yellow">{t`${splitAndDisplay(data?.status, '_')}`}</Tag>
        );
    }
  };
  return (
    <Stack
      border="1px solid #DFE0EB"
      bg={'white'}
      color={'#252733'}
      fontSize={18}
      p={4}
      my={8}
      >
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pt={4} mb={2} pl={5}>
          {t`Your TRISA ${type} Details`}
        </Heading>
        <Stack fontSize={18} pb={4} overflow={'auto'}>
          <Table
            sx={{
              thead: { fontWeight: 'bold' },
            }}>
            <Thead>
              <Tr>
                <Th>
                  <Trans id="ID">ID</Trans>
                </Th>
                <Th>
                  <Trans id="First Listed">First Listed</Trans>
                </Th>
                <Th>
                  <Trans id="Verified On">Verified On</Trans>
                </Th>
                <Th>
                  <Trans id="Last Updated">Last Updated</Trans>
                </Th>
                <Th>
                  <Trans id="Status">Status</Trans>
                </Th>
              </Tr>
            </Thead>
            <Tbody>
              <Tr>
                <Td>{data?.id || 'N/A'}</Td>
                <Td>{data?.first_listed ? format2ShortDate(data?.first_listed) : 'N/A'}</Td>
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
