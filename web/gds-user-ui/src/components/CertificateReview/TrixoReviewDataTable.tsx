import { Stack, Table, Tbody, Tr, Td, Tag, TagLabel, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { COUNTRIES } from 'constants/countries';
import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';
import getColorScheme from 'utils/getColorScheme';
import { TStep, loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { currencyFormatter } from 'utils/utils';
import { t } from '@lingui/macro';
interface TrixoReviewProps {
  data?: any;
}
function TrixoReviewDataTable({ data }: TrixoReviewProps) {
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
                colorScheme={getColorScheme(data.financial_transfers_permitted)}>
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
                Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes
                where it is licensed/approved/registered?
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
                colorScheme={getColorScheme(data?.financial_transfers_permitted)}>
                <TagLabel fontWeight={'bold'}>
                  {data?.financial_transfers_permitted?.toUpperCase()}
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
              {data?.kyc_threshold ? (
                <Text>
                  {currencyFormatter(data?.kyc_threshold, {
                    currency: data?.kyc_threshold_currency
                  }) || 'USD'}{' '}
                  {data?.kyc_threshold_currency}
                </Text>
              ) : (
                'N/A'
              )}
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
            <Td>
              {data?.applicable_regulations?.map((o: any, i: any) => {
                if (o?.length > 0) {
                  return <Text key={i}>{o || 'N/A'}</Text>;
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
              {currencyFormatter(data?.compliance_threshold, {
                currency: data?.compliance_threshold_currency
              }) || 'N/A'}{' '}
              {data?.compliance_threshold_currency || 'N/A'}
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
                colorScheme={getColorScheme(data?.must_safeguard_pii)}>
                <TagLabel fontWeight={'bold'}>{data?.must_safeguard_pii ? 'YES' : 'NO'}</TagLabel>
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
                colorScheme={getColorScheme(data?.safeguards_pii)}>
                <TagLabel fontWeight={'bold'}>{data?.safeguards_pii ? 'YES' : 'NO'}</TagLabel>
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
