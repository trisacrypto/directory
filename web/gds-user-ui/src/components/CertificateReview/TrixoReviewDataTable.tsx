import { Stack, Table, Tbody, Tr, Td, Tag, TagLabel, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { COUNTRIES } from 'constants/countries';
import getColorScheme from 'utils/getColorScheme';
import { t } from '@lingui/macro';
import ComplianceThresholdRow from './ComplianceThresholdRow';
import KycThresholdRow from './KycThresholdRow';
interface TrixoReviewProps {
  data?: any;
}
function TrixoReviewDataTable({ data }: TrixoReviewProps) {
  const getConductsCustomerKYC = (conductsCustomerKYC: boolean) => {
    return conductsCustomerKYC ? t`Yes` : t`No`;
  };

  return (
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
            <Td>{(COUNTRIES as any)[data?.primary_national_jurisdiction] || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="Name of Primary Regulator">Name of Primary Regulator</Trans>
            </Td>
            <Td>{data?.primary_regulator || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="Other Jurisdictions">Other Jurisdictions</Trans>
            </Td>
            <Td>
              {data?.other_jurisdictions?.length > 0
                ? data?.other_jurisdictions?.map((o: any, i: any) => {
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
                Is your organization permitted to send and/or receive transfers of virtual assets in
                the jurisdictions in which it operates?
              </Trans>
            </Td>
            <Td>
              <Tag
                size={'sm'}
                key={'sm'}
                variant="subtle"
                colorScheme={getColorScheme(data?.financial_transfers_permitted)}>
                <TagLabel fontWeight={'bold'}>
                  {data?.financial_transfers_permitted?.toString().toUpperCase()}
                </TagLabel>
              </Tag>
            </Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td colSpan={3} fontWeight="bold" background="#E5EDF1" pl={'1rem !important'}>
              <Trans id="Customer Due Diligence (CDD) & Travel Rule Policies">
                Customer Due Diligence (CDD) & Travel Rule Policies
              </Trans>
            </Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans
                id="Does your organization have a programme that sets minimum Anti-Money Laundering (AML), Countering the Financing of Terrorism (CFT), Know your Counterparty/Customer Due
                    Diligence (KYC/CDD) and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?">
                Does your organization have a programme that sets minimum Anti-Money Laundering
                (AML), Countering the Financing of Terrorism (CFT), Know your Counterparty/Customer
                Due Diligence (KYC/CDD) and Sanctions standards per the requirements of the
                jurisdiction(s) regulatory regimes where it is licensed/approved/registered?
              </Trans>
            </Td>
            <Td>
              <Tag
                size={'sm'}
                key={'sm'}
                variant="subtle"
                colorScheme={getColorScheme(data?.has_required_regulatory_program)}>
                <TagLabel fontWeight={'bold'}>
                  {data?.has_required_regulatory_program?.toUpperCase() || 'N/A'}
                </TagLabel>
              </Tag>
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="Does your organization conduct Know your KYC/CDD before permitting its customers to send/receive virtual asset transfers?">
                Does your organization conduct Know your KYC/CDD before permitting its customers to
                send/receive virtual asset transfers?
              </Trans>
            </Td>
            <Td>
              <Tag
                size={'sm'}
                key={'sm'}
                variant="subtle"
                colorScheme={getColorScheme(data?.conducts_customer_kyc || 'no')}>
                <TagLabel fontWeight={'bold'}>
                  {getConductsCustomerKYC(data?.conducts_customer_kyc || false)}
                </TagLabel>
              </Tag>
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="At what threshold and currency does your organization conduct KYC checks?">
                At what threshold and currency does your organization conduct KYC checks?
              </Trans>
            </Td>
            <Td pl={0}>
              <KycThresholdRow data={data} />
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
                colorScheme={getColorScheme(data?.must_comply_travel_rule)}>
                <TagLabel fontWeight={'bold'}>
                  {data?.must_comply_travel_rule ? 'YES' : 'NO'}
                </TagLabel>
              </Tag>
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="Applicable Regulations">Applicable Regulations</Trans>
            </Td>
            <Td display="flex" flexWrap="wrap" gap={1}>
              {data?.applicable_regulations?.map((o: any, i: any) => {
                if (o?.length > 0) {
                  return (
                    <Tag bg={'blue'} color={'white'} key={i}>
                      {o}
                    </Tag>
                  );
                } else {
                  return <Text key={i}>{'N/A'}</Text>;
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
              <ComplianceThresholdRow data={data} />
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td colSpan={3} background="#E5EDF1" fontWeight="bold" pl={'1rem !important'}>
              <Trans id="Data Protection Policies">Data Protection Policies</Trans>
            </Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans
                id="Is your organization required by law to safeguard Personally Identifiable
                Information (PII)?">
                Is your organization required by law to safeguard Personally Identifiable
                Information (PII)?
              </Trans>
            </Td>
            <Td>
              <Tag
                size={'sm'}
                key={'sm'}
                variant="subtle"
                colorScheme={getColorScheme(data?.must_safeguard_pii)}>
                <TagLabel fontWeight={'bold'}>{data?.must_safeguard_pii ? t`YES` : t`NO`}</TagLabel>
              </Tag>
            </Td>
            <Td></Td>
          </Tr>
          <Tr>
            <Td>
              <Trans
                id="Does your organization secure and protect Personally Identifiable
                Information (PII), including Personally Identifiable
                Information (PII) received from other VASPs under the Travel Rule?">
                Does your organization secure and protect Personally Identifiable Information (PII),
                including Personally Identifiable Information (PII) received from other VASPs under
                the Travel Rule?
              </Trans>
            </Td>
            <Td>
              <Tag
                size={'sm'}
                key={'sm'}
                variant="subtle"
                colorScheme={getColorScheme(data?.safeguards_pii)}>
                <TagLabel fontWeight={'bold'}>{data?.safeguards_pii ? t`YES` : t`NO`}</TagLabel>
              </Tag>
            </Td>
            <Td></Td>
          </Tr>
        </Tbody>
      </Table>
    </Stack>
  );
}

export default TrixoReviewDataTable;
