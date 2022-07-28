import { Stack, Table, Tbody, Tr, Td, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { TStep, loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';

function TrisaImplementationReviewDataTable() {
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
            <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
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
              <Trans id="TestNet Certificate Common Name">TestNet Certificate Common Name</Trans>
            </Td>
            <Td pl={0}>{trisa?.testnet?.common_name || 'N/A'}</Td>
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
            <Td pl={0}>{trisa?.mainnet?.endpoint || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="MainNet Certificate Common Name">MainNet Certificate Common Name</Trans>
            </Td>
            <Td pl={0}>{trisa?.mainnet?.common_name || 'N/A'}</Td>
          </Tr>
        </Tbody>
      </Table>
    </Stack>
  );
}

export default TrisaImplementationReviewDataTable;
