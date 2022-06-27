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
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
interface TrixoReviewProps {}

const TrixoReview: React.FC<TrixoReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [trixo, setTrixo] = React.useState<any>({});
  const getColorScheme = (status: string | boolean) => {
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
            <Trans id="Section 5: TRIXO Questionnaire">Section 5: TRIXO Questionnaire</Trans>
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
            {t`Edit`}{' '}
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              'td:nth-child(2)': { maxWidth: '75%' },
              Tr: { borderStyle: 'hidden' },
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
            <Tbody>
              <Tr>
                <Td>
                  <Trans id="Primary National Jurisdiction">Primary National Jurisdiction</Trans>
                </Td>
                <Td>{(COUNTRIES as any)[trixo?.primary_national_jurisdiction] || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Name of Primary Regulator">Name of Primary Regulator</Trans>
                </Td>
                <Td>{trixo?.primary_regulator || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Other Jurisdictions">Other Jurisdictions</Trans>
                </Td>
                <Td>
                  {trixo?.other_jurisdictions?.length > 0
                    ? trixo?.other_jurisdictions?.map((o: any, i: any) => {
                        if (o?.regulator_name?.length > 0) {
                          return (
                            <Text key={i}>
                              {o.country} : {o.regulator_name}{' '}
                            </Text>
                          );
                        }
                      })
                    : 'N/A'}
                </Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?">
                    Is your organization permitted to send and/or receive transfers of virtual
                    assets in the jurisdictions in which it operates?
                  </Trans>
                </Td>
                <Td>
                  <Tag
                    size={'sm'}
                    key={'sm'}
                    variant="subtle"
                    colorScheme={getColorScheme(trixo.financial_transfers_permitted)}>
                    <TagLabel fontWeight={'bold'}>
                      {trixo?.financial_transfers_permitted?.toString().toUpperCase()}
                    </TagLabel>
                  </Tag>
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  <Trans id="CDD & Travel Rule Policies">CDD & Travel Rule Policies</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?">
                    Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and
                    Sanctions standards per the requirements of the jurisdiction(s) regulatory
                    regimes where it is licensed/approved/registered?
                  </Trans>
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
                  <Trans id="Does your organization conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers?">
                    Does your organization conduct KYC/CDD before permitting its customers to
                    send/receive virtual asset transfers?
                  </Trans>
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
                <Td>
                  <Trans id="At what threshold and currency does your organization conduct KYC?">
                    At what threshold and currency does your organization conduct KYC?
                  </Trans>
                </Td>
                <Td pl={0}>
                  <Text>
                    {currencyFormatter(trixo.kyc_threshold, {
                      currency: trixo.kyc_threshold_currency
                    }) || 'N/A'}{' '}
                    {trixo.kyc_threshold_currency || 'N/A'}
                  </Text>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Is your organization required to comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered?">
                    Is your organization required to comply with the application of the Travel Rule
                    standards in the jurisdiction(s) where it is licensed/approved/registered?
                  </Trans>
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
                <Td>
                  <Trans id="Applicable Regulations">Applicable Regulations</Trans>
                </Td>
                <Td>
                  {trixo?.applicable_regulations?.map((reg: any) => {
                    if (reg?.name.length > 0) {
                      return <Text key={reg.name}>{reg.name || 'N/A'}</Text>;
                    }
                  })}
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="What is the minimum threshold for Travel Rule compliance?">
                    What is the minimum threshold for Travel Rule compliance?
                  </Trans>
                </Td>
                <Td pl={0}>
                  {currencyFormatter(trixo.compliance_threshold, {
                    currency: trixo.compliance_threshold_currency
                  }) || 'N/A'}{' '}
                  {trixo.compliance_threshold_currency || 'N/A'}
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td colSpan={2} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
                  <Trans id="Data Protection Policies">Data Protection Policies</Trans>
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Is your organization required by law to safeguard PII?">
                    Is your organization required by law to safeguard PII?
                  </Trans>
                </Td>
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
                  <Trans id="Does your organization secure and protect PII, including PII received from other VASPs under the Travel Rule?">
                    Does your organization secure and protect PII, including PII received from other
                    VASPs under the Travel Rule?
                  </Trans>
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
