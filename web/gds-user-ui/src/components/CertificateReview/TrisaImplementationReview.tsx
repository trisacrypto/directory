import React, { FC, useEffect } from 'react';
import {
  Stack,
  Box,
  Text,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
  Divider,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
interface TrisaImplementationReviewProps {}

const TrisaImplementationReview = (props: TrisaImplementationReviewProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trisa, setTrisa] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

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
      color={textColor}
      fontSize={'1rem'}
      p={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            <Trans id="Section 4: TRISA Implementation">Section 4: TRISA Implementation</Trans>
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(4)}
            _hover={{
              bg: '#10aaed'
            }}>
            <Trans id="Edit">Edit</Trans>
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody
              sx={{
                ' td': {
                  fontSize: '1rem'
                },
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
                <Td
                  colSpan={2}
                  background={useColorModeValue('#E5EDF1', 'gray.900')}
                  fontWeight="bold"
                  pl={'1rem !important'}>
                  <Trans id="TestNet">TestNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>
                  <Trans id="TestNet TRISA Endpoint">TestNet TRISA Endpoint</Trans>
                </Td>
                <Td pl={0}>{trisa?.testnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="TestNet Certificate Common Name">
                    TestNet Certificate Common Name
                  </Trans>
                </Td>
                <Td pl={0}>{trisa?.testnet?.common_name || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td colSpan={2}></Td>
              </Tr>
              <Tr>
                <Td
                  colSpan={2}
                  background={useColorModeValue('#E5EDF1', 'gray.900')}
                  fontWeight="bold"
                  pl={'1rem !important'}>
                  <Trans id="MainNet">MainNet</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td pt={'1rem !important'}>
                  <Trans id="MainNet TRISA Endpoint">MainNet TRISA Endpoint</Trans>
                </Td>
                <Td pl={0}>{trisa?.mainnet?.endpoint || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="MainNet Certificate Common Name">
                    MainNet Certificate Common Name
                  </Trans>
                </Td>
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
