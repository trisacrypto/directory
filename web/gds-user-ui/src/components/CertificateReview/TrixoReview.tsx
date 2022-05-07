import React, { useEffect } from 'react';
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
  Tag,
  TagLabel
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { COUNTRIES } from 'constants/countries';
import { currencyFormatter } from 'utils/utils';
interface TrixoReviewProps {}

const TrixoReview: React.FC<TrixoReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trixo, setTrixo] = React.useState<any>({});
  const getColorScheme = (status: string | boolean) => {
    console.log('[getColorScheme] status', status);
    if (status === 'yes' || status === true) {
      return 'green';
    } else {
      return 'orange';
    }
  };
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.trixo
    };
    setTrixo(stepData);
  }, [steps]);

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
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            Section 5: TRIXO Questionnaire
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(5)}
            _hover={{
              bg: '#10aaed'
            }}>
            {' '}
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' },
              'td:nth-child(2)': { maxWidth: '75%' },
              Tr: { borderStyle: 'hidden' },
              'td:first-child': {
                width: '50%'
              },
              td: {
                borderBottom: 'none',
                paddingInlineStart: 0,
                paddingY: 2.5
              }
            }}>
            <Tbody>
              <Tr>
                <Td>Primary National Jurisdiction</Td>
                <Td>{(COUNTRIES as any)[trixo?.primary_national_jurisdiction] || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>Name of Primary Regulator</Td>
                <Td>{trixo?.primary_regulator || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>Other Jurisdictions</Td>
                <Td>
                  {trixo?.other_jurisdictions?.length > 0
                    ? trixo?.other_jurisdictions?.map((o: any, i: any) => {
                        if (o?.regulator_name?.length > 0) {
                          return (
                            <>
                              <Tr>
                                <Td>{o.country}</Td>
                                <Td>{o.regulator_name}</Td>
                              </Tr>
                            </>
                          );
                        }
                      })
                    : 'N/A'}
                </Td>
              </Tr>
              <Tr>
                <Td>
                  Is your organization permitted to send and/or receive transfers of virtual assets
                  in the jurisdictions in which it operates?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.financial_transfers_permitted)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo?.financial_transfers_permitted?.toUpperCase()}
                    </TagLabel>
                  </Tag>
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  CDD & Travel Rule Policies
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and
                  Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes
                  where it is licensed/approved/registered?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.has_required_regulatory_program)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo?.has_required_regulatory_program?.toUpperCase() || 'N/A'}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization conduct KYC/CDD before permitting its customers to
                  send/receive virtual asset transfers?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo?.financial_transfers_permitted)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo?.financial_transfers_permitted?.toUpperCase()}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>At what threshold and currency does your organization conduct KYC?</Td>
                <Td pl={0}>
                  <Tr>
                    <Td>
                      {currencyFormatter(trixo.kyc_threshold, {
                        currency: trixo.kyc_threshold_currency
                      }) || 'N/A'}
                    </Td>
                    <Td pl={0}>{trixo.kyc_threshold_currency || 'N/A'}</Td>
                  </Tr>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Is your organization required to comply with the application of the Travel Rule
                  standards in the jurisdiction(s) where it is licensed/approved/registered?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.must_comply_travel_rule)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo?.must_comply_travel_rule ? 'YES' : 'NO'}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Applicable Regulations</Td>
                <Td>
                  <Tr>
                    <Td pl={0}>
                      {trixo?.applicable_regulations?.map((reg: any) => {
                        if (reg?.name.length > 0) {
                          return (
                            <React.Fragment key={reg.name}>{reg.name || 'N/A'}</React.Fragment>
                          );
                        }
                      })}
                    </Td>
                    <Td></Td>
                  </Tr>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>What is the minimum threshold for Travel Rule compliance?</Td>
                <Td pl={0}>
                  <Tr>
                    <Td>
                      {currencyFormatter(trixo.compliance_threshold, {
                        currency: trixo.compliance_threshold_currency
                      }) || 'N/A'}
                    </Td>
                    <Td pl={0}>{trixo.compliance_threshold_currency || 'N/A'}</Td>
                  </Tr>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  Data Protection Policies
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Is your organization required by law to safeguard PII?</Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.must_safeguard_pii)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo.must_safeguard_pii ? 'YES' : 'NO'}
                    </TagLabel>
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  Does your organization secure and protect PII, including PII received from other
                  VASPs under the Travel Rule?
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.safeguards_pii)}>
                    <TagLabel fontWeight={'bold'}>{trixo.safeguards_pii ? 'YES' : 'NO'}</TagLabel>
                  </Tag>
                </Td>
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
