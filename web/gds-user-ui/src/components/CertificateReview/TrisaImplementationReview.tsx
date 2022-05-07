import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button, Divider } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
interface TrisaImplementationReviewProps {}

const TrisaImplementationReview = (props: TrisaImplementationReviewProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trisa, setTrisa] = React.useState<any>({});
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      mainnet: getStepperData.trisa_endpoint_mainnet,
      testnet: getStepperData.trisa_endpoint_testnet
    };

    setTrisa(stepData);
  }, [steps]);
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      bg={'white'}
      fontSize={'1rem'}
      p={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            Section 4: TRISA Implementation
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(4)}
            _hover={{
              bg: '#10aaed'
            }}>
            Edit
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody
              sx={{
                ' td': {
                  fontSize: '1rem'
                },
                'td:first-child': {
                  width: '50%'
                },
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5
                }
              }}>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  TestNet
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>TestNet TRISA Endpoint</Td>
                <Td pl={0}>{trisa?.testnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>TestNet Certificate Common Name</Td>
                <Td pl={0}>{trisa?.testnet?.common_name || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td colSpan={2}></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  MainNet
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>MainNet TRISA Endpoint</Td>
                <Td pl={0}>{trisa?.mainnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>MainNet Certificate Common Name</Td>
                <Td pl={0}>{trisa?.mainnet?.common_name || 'N/A'}</Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
TrisaImplementationReview.defaultProps = {
  data: {}
};
export default TrisaImplementationReview;
