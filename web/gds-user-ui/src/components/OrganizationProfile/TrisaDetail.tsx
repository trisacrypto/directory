import React, { useState } from 'react';
import {
  Box,
  Heading,
  VStack,
  Flex,
  Input,
  Stack,
  Text,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  SimpleGrid,
  List,
  ListItem,
  Table,
  Tbody,
  Tr,
  Td,
  HStack,
  Tag
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
type TrisaDetailProps = {
  data: any;
};
const TrisaDetail: React.FC<TrisaDetailProps> = ({ data }) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      fontSize={18}
      p={4}>
      <Stack width={'100%'}>
        <Heading as={'h1'} fontSize={19} pb={10}>
          {' '}
          TRISA Details{' '}
        </Heading>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody
              sx={{
                '*': {
                  fontSize: '1rem',
                  'td:first-child': {
                    width: '50%',
                    paddingLeft: '1rem'
                  },
                  td: {
                    borderBottom: 'none',
                    paddingInlineStart: 0,
                    paddingY: 2.5
                  }
                }
              }}>
              <Tr>
                <Td borderBottom={'none'} pl={'1rem !important'}>
                  ID
                </Td>
                <Td borderBottom={'none'} whiteSpace="break-spaces" lineHeight={1.5}>
                  {data?.organization?.vasp_id || 'N/A'}
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td pl={'1rem !important'}>
                  <Trans id="Fist Listed">Fist Listed</Trans>
                </Td>
                <Td>{data?.organization?.first_listed || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="Verified On">Verified On</Trans>
                </Td>
                <Td>{data?.organization?.verified_on || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="Last Updated">Last Updated</Trans>
                </Td>
                <Td>{data?.organization?.last_updated || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  <Trans id="TestNet">TestNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>
                  <Trans id="TestNet TRISA Endpoint">TestNet TRISA Endpoint</Trans>
                </Td>
                <Td pl={0}>{data?.testnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="TestNet Certificate Common Name">
                    TestNet Certificate Common Name
                  </Trans>
                </Td>
                <Td pl={0}>{data?.testnet?.common_name || 'N/A'}</Td>
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
                <Td pt={'1rem !important'}>MainNet TRISA Endpoint</Td>
                <Td pl={0}>{data?.mainnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="MainNet Certificate Common Name">
                    MainNet Certificate Common Name
                  </Trans>
                </Td>
                <Td pl={0}>{data?.mainnet?.common_name || 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
export default TrisaDetail;
