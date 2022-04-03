import React, { FC } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
interface LegalSectionProps {
  data: any;
}

const TrixoReview = (props: LegalSectionProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      bg={'white'}
      fontSize={18}
      p={5}
      px={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Section 5: TRIXO Questionnaire</Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            {' '}
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody>
              <Tr>
                <Td>Primary National Jurisdiction</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Name of Primary Regulator</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Other Jurisdictions</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>CDD & Travel Rule Policies</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and
                  Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes
                  where it is licensed/approved/registered?
                </Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization conduct KYC/CDD before permitting its customers to
                  send/receive virtual asset transfers?
                </Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Is your organization required to comply with the application of the Travel Rule
                  standards in the jurisdiction(s) where it is licensed/approved/registered?
                </Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>What is the minimum threshold for Travel Rule compliance?</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Data Protection Policies</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Is your organization required by law to safeguard PII?</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization secure and protect PII, including PII received from other
                  VASPs under the Travel Rule?
                </Td>
                <Td></Td>
                <Td></Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
TrixoReview.defaultProps = {
  data: {}
};
export default TrixoReview;
