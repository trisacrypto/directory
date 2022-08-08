import { Stack, Table, Tbody, Tr, Td, Tag, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { addressType } from 'constants/address';
import { COUNTRIES } from 'constants/countries';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import React, { useEffect } from 'react';
import { useSelector, RootStateOrAny } from 'react-redux';
import { renderAddress } from 'utils/address-utils';
import { TStep, loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';

interface LegalReviewProps {
  data?: any;
}
function LegalPersonReviewDataTable({ data }: LegalReviewProps) {
  return (
    <Stack fontSize={18}>
      <Table
        sx={{
          'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold', paddingLeft: 0 },
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
        <Tbody
          sx={{
            '*': {
              fontSize: '1rem'
            }
          }}>
          <Tr>
            <Td fontSize={'1rem'} fontWeight="bold" colSpan={3} background="#E5EDF1">
              <Trans id="Name Identifiers">Name Identifiers</Trans>
            </Td>
          </Tr>
          <Tr>
            <Td pt={0}>
              <Trans id="The name and type of name by which the legal person is known.">
                The name and type of name by which the legal person is known.
              </Trans>
            </Td>
            <Td>
              <Tbody>
                <Tr>
                  {data.name?.name_identifiers?.map((nameIdentifier: any, index: number) => {
                    return (
                      <React.Fragment key={index}>
                        <Td paddingLeft={'0 !important'} border="none">
                          {nameIdentifier.legal_person_name || 'N/A'}
                        </Td>
                        {nameIdentifier.legal_person_name_identifier_type ? (
                          <Td paddingLeft={0} border="none">
                            (
                            {getNameIdentiferTypeLabel(
                              nameIdentifier.legal_person_name_identifier_type
                            )}
                            )
                          </Td>
                        ) : null}
                      </React.Fragment>
                    );
                  })}
                </Tr>
              </Tbody>
              <>
                {data.name?.local_name_identifiers?.map((nameIdentifier: any, index: number) => {
                  return (
                    <React.Fragment key={index}>
                      <Td paddingLeft={0} pt={0} border="none">
                        {nameIdentifier.legal_person_name}
                      </Td>
                      <Td paddingLeft={0} pt={0} border="none">
                        (
                        {getNameIdentiferTypeLabel(
                          nameIdentifier.legal_person_name_identifier_type
                        )}
                        )
                      </Td>
                    </React.Fragment>
                  );
                })}
              </>
              <>
                {data.name?.phonetic_name_identifiers?.map((nameIdentifier: any, index: number) => {
                  return (
                    <React.Fragment key={index}>
                      <Td paddingLeft={0} pt={0} border="none">
                        {nameIdentifier.legal_person_name}
                      </Td>
                      <Td paddingLeft={0} pt={0} border="none">
                        (
                        {getNameIdentiferTypeLabel(
                          nameIdentifier.legal_person_name_identifier_type
                        )}
                        )
                      </Td>
                    </React.Fragment>
                  );
                })}
              </>
            </Td>
          </Tr>

          <Tr>
            <Td pt={0} pl={'1rem !important'}>
              <Trans id="Addresses">Addresses</Trans>
            </Td>
            <Td pl={0}>
              <Tr>
                <Td paddingLeft={'0px !important'} pt={0}>
                  {data?.geographic_addresses?.map((address: any, index: number) => (
                    <React.Fragment key={index}>
                      {data?.geographic_addresses?.length > 1 && (
                        <Text py={1} fontWeight={'bold'}>
                          Address {index + 1} : {(addressType as any)[address.address_type]}
                        </Text>
                      )}
                      {renderAddress(address)}
                    </React.Fragment>
                  ))}
                </Td>
                {data?.geographic_addresses?.length === 1 && (
                  <Td pt={0}>
                    ({(addressType as any)[data?.geographic_addresses?.[0].address_type]})
                  </Td>
                )}
              </Tr>
            </Td>
          </Tr>

          <Tr>
            <Td pt={0} pl={'1rem !important'}>
              <Trans id="Country of Registration">Country of Registration</Trans>
            </Td>
            <Td paddingLeft={0} pt={0}>
              <Tbody>
                <Tr>
                  <Td pl={'0 !important'}>
                    {(COUNTRIES as any)[data?.country_of_registration] || 'N/A'}
                  </Td>
                </Tr>
              </Tbody>
            </Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td fontSize={'1rem'} pt={'2rem'} fontWeight="bold" colSpan={3} background="#E5EDF1">
              <Text mb={1}>
                <Trans id="National Identification">National Identification</Trans>
              </Text>
            </Td>
          </Tr>
          <Tr>
            <Td></Td>
          </Tr>
          <Tr>
            <Td pt={0}>
              <Trans id="Identification Number">Identification Number</Trans>
            </Td>
            <Td paddingLeft={0}>{data?.national_identification?.national_identifier}</Td>
          </Tr>
          <Tr>
            <Td pt={0}>
              <Trans id="Identification Type">Identification Type</Trans>
            </Td>
            <Td pt={0}>
              <Tag color={'white'} bg={'blue'} size={'lg'}>
                {getNationalIdentificationLabel(
                  data?.national_identification?.national_identifier_type
                )}
              </Tag>
            </Td>
          </Tr>
          <Tr>
            <Td>
              <Trans id="Country of Registration">Country of Registration</Trans>
            </Td>
            <Td>{(COUNTRIES as any)[data?.country_of_registration] || 'N/A'}</Td>
          </Tr>
          <Tr>
            <Td pt={0}>
              <Trans id="Reg Authority">Reg Authority</Trans>
            </Td>
            <Td pt={0}>{data?.national_identification?.registration_authority || 'N/A'}</Td>
          </Tr>
        </Tbody>
      </Table>
    </Stack>
  );
}

export default LegalPersonReviewDataTable;
