import {
  Grid,
  GridItem,
  Heading,
  HStack,
  Stack,
  StackProps,
  Table,
  TableCaption,
  TableContainer,
  Tag,
  TagLabel,
  Tbody,
  Td,
  Text,
  Tfoot,
  Th,
  Thead,
  Tr,
  VStack
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import { COUNTRIES } from 'constants/countries';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import React from 'react';
import { renderAddress } from 'utils/address-utils';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { currencyFormatter, getColorScheme, hasValue } from 'utils/utils';

type SimpleLayoutProps = { children: React.ReactNode } & StackProps;
const SimpleLayout: React.FC<SimpleLayoutProps> = ({ children, ...rest }) => {
  return (
    <Stack border="3px solid #E5EDF1" borderRadius="10px" p={6} align="start" {...rest}>
      {children}
    </Stack>
  );
};

export default function MemberDetails() {
  const certificate = loadDefaultValueFromLocalStorage();

  return (
    <Stack>
      <Heading size="lg" fontWeight="bold">
        Member Details
      </Heading>

      <Stack spacing={6}>
        <SimpleLayout>
          <Heading size="md" mb={2}>
            VaspBot Holdings
          </Heading>

          <Stack>
            <Text>TRISA Member ID: 0ffe693f-9752-45cb-830f-7dbae12d6baf</Text>
            <Text>TRISA Verification: VERIFIED on 2022-02-01T17:13:02Z</Text>
          </Stack>
        </SimpleLayout>

        <Stack direction="row" spacing={6}>
          <SimpleLayout height="100%" maxHeight="350px" width="100%" maxW="50%" overflowY="scroll">
            <Heading size="md" mb={2}>
              TestNet: Active
            </Heading>
            <Stack wordBreak="break-word">
              <Text>TRISA TestNet Identity Signature: </Text>
              <Text>
                LRJrjvxOZ77z7tZz6XGKoaTHxbl9ojKu2Ouf1tCZ7mauabqWHLJhXufdvI9gwAD6HH5P0hBLeBY7Si5mxmfItFDN2B+ejiSozMUDtBxa0GHipeTMJ1Xo2DtD1+NJoTSVt2ycrYyQlZtqjdvpA6weCp3SS4brTzVeVWnrIU1IEtXFb2HwHzvIAjCjXFBAgWZ3wLIC+1VOTbJyZBc/7Rtuild9BdqY26fnrgJ2Rq622McY9jAp4h3I5XWL5VQwi30YcM3lG8nSqfk5DRyjAsqifsk9xC68t5QBu7TJgbzw2fkg0P0yqxsk8tC2huDt1Vv4sH8s+iXRFhjB5cz7cgEFtrzuISRppx1a99Hw1dO4Mc7IeylhEvAo8CnWyy2NBuKEU1QubwYEbWCOhSLSxS+7ye4T2cCCtm4S0o026+usG2IubuVW2DE/6lV1njgHwWSPQLvkA+c+0aMSdTKWzDVyaLvqL17jBDQVA+LFyNrWCW8urBXiZbMsyHwxT3uniX42Jp9V68skYf7mFSMtejNLBNeo9CyBf9zrPd80IfdCOp+fJfY8sidfAS3sLlxRUZLM0Lri+/jqSdMepO3ZFPQUuJ++ptl5PyMpPCSuaNh4Wje0OXGHW4tIdgfnx+4u157ksnnVuuashQK6sV/ETueXh6M8zMrn2lv0SB4=
              </Text>
            </Stack>
          </SimpleLayout>
          <SimpleLayout height="100%" maxH="350px" width="100%" maxW="50%" overflowY="scroll">
            <Heading size="md" mb={2}>
              TestNet: Active
            </Heading>
            <Stack>
              <Text>TRISA TestNet Identity Signature: </Text>
              <Text wordBreak="break-word">
                LRJrjvxOZ77z7tZz6XGKoaTHxbl9ojKu2Ouf1tCZ7mauabqWHLJhXufdvI9gwAD6HH5P0hBLeBY7Si5mxmfItFDN2B+ejiSozMUDtBxa0GHipeTMJ1Xo2DtD1+NJoTSVt2ycrYyQlZtqjdvpA6weCp3SS4brTzVeVWnrIU1IEtXFb2HwHzvIAjCjXFBAgWZ3wLIC+1VOTbJyZBc/7Rtuild9BdqY26fnrgJ2Rq622McY9jAp4h3I5XWL5VQwi30YcM3lG8nSqfk5DRyjAsqifsk9xC68t5QBu7TJgbzw2fkg0P0yqxsk8tC2huDt1Vv4sH8s+iXRFhjB5cz7cgEFtrzuISRppx1a99Hw1dO4Mc7IeylhEvAo8CnWyy2NBuKEU1QubwYEbWCOhSLSxS+7ye4T2cCCtm4S0o026+usG2IubuVW2DE/6lV1njgHwWSPQLvkA+c+0aMSdTKWzDVyaLvqL17jBDQVA+LFyNrWCW8urBXiZbMsyHwxT3uniX42Jp9V68skYf7mFSMtejNLBNeo9CyBf9zrPd80IfdCOp+fJfY8sidfAS3sLlxRUZLM0Lri+/jqSdMepO3ZFPQUuJ++ptl5PyMpPCSuaNh4Wje0OXGHW4tIdgfnx+4u157ksnnVuuashQK6sV/ETueXh6M8zMrn2lv0SB4=
                LRJrjvxOZ77z7tZz6XGKoaTHxbl9ojKu2Ouf1tCZ7mauabqWHLJhXufdvI9gwAD6HH5P0hBLeBY7Si5mxmfItFDN2B+ejiSozMUDtBxa0GHipeTMJ1Xo2DtD1+NJoTSVt2ycrYyQlZtqjdvpA6weCp3SS4brTzVeVWnrIU1IEtXFb2HwHzvIAjCjXFBAgWZ3wLIC+1VOTbJyZBc/7Rtuild9BdqY26fnrgJ2Rq622McY9jAp4h3I5XWL5VQwi30YcM3lG8nSqfk5DRyjAsqifsk9xC68t5QBu7TJgbzw2fkg0P0yqxsk8tC2huDt1Vv4sH8s+iXRFhjB5cz7cgEFtrzuISRppx1a99Hw1dO4Mc7IeylhEvAo8CnWyy2NBuKEU1QubwYEbWCOhSLSxS+7ye4T2cCCtm4S0o026+usG2IubuVW2DE/6lV1njgHwWSPQLvkA+c+0aMSdTKWzDVyaLvqL17jBDQVA+LFyNrWCW8urBXiZbMsyHwxT3uniX42Jp9V68skYf7mFSMtejNLBNeo9CyBf9zrPd80IfdCOp+fJfY8sidfAS3sLlxRUZLM0Lri+/jqSdMepO3ZFPQUuJ++ptl5PyMpPCSuaNh4Wje0OXGHW4tIdgfnx+4u157ksnnVuuashQK6sV/ETueXh6M8zMrn2lv0SB4=
              </Text>
            </Stack>
          </SimpleLayout>
        </Stack>

        <SimpleLayout>
          <Heading size="md" mb={2}>
            <Trans id="Basic Details">Basic Details</Trans>
          </Heading>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <Tbody>
                <Tr>
                  <Td>
                    <Trans id="Website">Website</Trans>
                  </Td>
                  <Td>{certificate?.website || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Date of Incorporation / Establishment">
                      Date of Incorporation / Establishment
                    </Trans>
                  </Td>
                  <Td>{certificate?.established_on || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Business Category">Business Category</Trans>
                  </Td>
                  <Td>{(BUSINESS_CATEGORY as any)[certificate.business_category] || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="VASP Category">VASP Category</Trans>
                  </Td>
                  <Td>
                    {certificate?.vasp_categories && certificate?.vasp_categories.length
                      ? certificate?.vasp_categories?.map((categ: any) => {
                          return (
                            <Tag key={categ} color={'white'} bg={'blue.400'} mr={2} mb={1}>
                              {getBusinessCategiryLabel(categ)}
                            </Tag>
                          );
                        })
                      : 'N/A'}
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>
        </SimpleLayout>

        <SimpleLayout>
          <Heading size="md" mb={2}>
            <Trans id="Legal Person">Legal Person</Trans>
          </Heading>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5,
                  verticalAlign: 'baseline'
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <TableCaption placement="top" color="#000" textAlign="start" paddingInlineStart={0}>
                <Trans id="Name Identifiers">Name Identifiers</Trans>
              </TableCaption>
              <Tbody>
                <Tr>
                  <Td fontStyle="italic">
                    <Trans id="The name and type of name by which the legal person is known.">
                      The name and type of name by which the legal person is known.
                    </Trans>
                  </Td>
                  <Td>
                    {certificate.entity.name?.name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <Stack direction="row" key={index}>
                            <Stack width="50%">
                              <Text as="div">{nameIdentifier.legal_person_name || 'N/A'}</Text>
                            </Stack>
                            <Stack>
                              <Text as="div">
                                {getNameIdentiferTypeLabel(
                                  nameIdentifier.legal_person_name_identifier_type
                                )}
                              </Text>
                            </Stack>
                          </Stack>
                        );
                      }
                    )}

                    {certificate.entity.name?.local_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <Stack direction="row" key={index}>
                            <Stack width="50%">
                              <Text as="div">{nameIdentifier.legal_person_name || 'N/A'}</Text>
                            </Stack>
                            <Stack>
                              <Text as="div">
                                {getNameIdentiferTypeLabel(
                                  nameIdentifier.legal_person_name_identifier_type
                                )}
                              </Text>
                            </Stack>
                          </Stack>
                        );
                      }
                    )}

                    <>
                      {certificate.entity.name?.phonetic_name_identifiers?.map(
                        (nameIdentifier: any, index: number) => {
                          return (
                            <Stack direction="row" key={index}>
                              <Stack width="50%">
                                <Text as="div">{nameIdentifier.legal_person_name || 'N/A'}</Text>
                              </Stack>
                              <Stack>
                                <Text as="div">
                                  {getNameIdentiferTypeLabel(
                                    nameIdentifier.legal_person_name_identifier_type
                                  )}
                                </Text>
                              </Stack>
                            </Stack>
                          );
                        }
                      )}
                    </>
                  </Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Addresses">Addresses</Trans>
                  </Td>
                  <Td>
                    <Stack direction="row">
                      <Stack width="50%">
                        <Text whiteSpace="break-spaces" lineHeight={1.5}>
                          {certificate.entity?.geographic_addresses?.map(
                            (address: any, index: number) => (
                              <React.Fragment key={index}>{renderAddress(address)}</React.Fragment>
                            )
                          )}
                        </Text>
                      </Stack>
                      <Stack>
                        <Text as="div">
                          {certificate.entity?.geographic_addresses?.[0] && 'Legal Person'}
                        </Text>
                      </Stack>
                    </Stack>
                  </Td>
                </Tr>
                {/* <Tr>
                  <Td>Customer Number</Td>
                  <Td>N/A</Td>
                </Tr> */}
                <Tr>
                  <Td>
                    <Trans id="Country of Registration">Country of Registration</Trans>
                  </Td>
                  <Td>
                    {(COUNTRIES as any)[certificate.entity?.country_of_registration] || 'N/A'}
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5,
                  verticalAlign: 'baseline'
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <TableCaption placement="top" color="#000" textAlign="start" paddingInlineStart={0}>
                <Trans id="National Identification">National Identification</Trans>
              </TableCaption>
              <Tbody>
                <Tr>
                  <Td>
                    <Trans id="Identification Number">Identification Number</Trans>
                  </Td>
                  <Td>
                    {certificate.entity?.national_identification?.national_identifier || 'N/A'}
                  </Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Identification Type">Identification Type</Trans>
                  </Td>
                  <Td>
                    {getNationalIdentificationLabel(
                      certificate.entity?.national_identification?.national_identifier_type
                    )}
                  </Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Country of Issue">Country of Issue</Trans>
                  </Td>
                  <Td>
                    {(COUNTRIES as any)[
                      certificate.entity?.national_identification?.country_of_issue
                    ] || 'N/A'}
                  </Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Country of Registration">Country of Registration</Trans>
                  </Td>
                  <Td>USA</Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="Reg Authority">Reg Authority</Trans>
                  </Td>
                  <Td>
                    {certificate.entity?.national_identification?.registration_authority || 'N/A'}
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>
        </SimpleLayout>

        <SimpleLayout>
          <Heading size="md" mb={2}>
            Contacts
          </Heading>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <Tbody>
                {['technical', 'legal', 'administrative', 'billing'].map((contact, index) => (
                  <Tr key={index}>
                    <Td textTransform="capitalize">
                      {contact === 'legal' ? `Compliance / ${contact}` : contact} Contact
                    </Td>
                    <Td>
                      {hasValue(certificate.contacts?.[contact]) ? (
                        <>
                          {certificate.contacts?.[contact]?.name && (
                            <>
                              {certificate.contacts?.[contact]?.name} <br />
                            </>
                          )}
                          {certificate.contacts?.[contact]?.email && (
                            <>
                              {certificate.contacts?.[contact]?.email} <br />
                            </>
                          )}
                          {certificate.contacts?.[contact]?.phone && (
                            <>
                              {certificate.contacts?.[contact]?.phone} <br />
                            </>
                          )}
                        </>
                      ) : (
                        'N/A'
                      )}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </TableContainer>
        </SimpleLayout>

        <SimpleLayout>
          <Heading size="md" mb={2}>
            TRISA Implementation
          </Heading>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <Tbody>
                <Tr>
                  <Td>TestNet TRISA Endpoint</Td>
                  <Td>{certificate?.trisa_endpoint_testnet?.trisa_endpoint || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>TestNet Certificate Common Name</Td>
                  <Td>{certificate?.trisa_endpoint_testnet?.common_name || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>MainNet TRISA Endpoint</Td>
                  <Td>{certificate?.trisa_endpoint_mainnet?.trisa_endpoint || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>MainNet Certificate Common Name</Td>
                  <Td>{certificate?.trisa_endpoint_mainnet?.common_name || 'N/A'}</Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>
        </SimpleLayout>

        <SimpleLayout>
          <Heading size="md" mb={2}>
            TRIXO Questionnaire
          </Heading>
          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5,
                  lineHeight: 1.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <Tbody>
                <Tr>
                  <Td>Primary National Jurisdiction</Td>
                  <Td>
                    {(COUNTRIES as any)[certificate.trixo?.primary_national_jurisdiction] || 'N/A'}
                  </Td>
                </Tr>
                <Tr>
                  <Td>Name of Primary Regulator</Td>
                  <Td>{certificate.trixo?.primary_regulator || 'N/A'}</Td>
                </Tr>
                <Tr>
                  <Td>Other Jurisdictions</Td>
                  <Td>
                    {certificate.trixo?.other_jurisdictions?.length > 0
                      ? certificate.trixo?.other_jurisdictions?.map((o: any, i: any) => {
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
                  <Td whiteSpace="break-spaces" lineHeight={1.5}>
                    Is your organization permitted to send and/or receive transfers of virtual
                    assets in the jurisdictions in which it operates?
                  </Td>
                  <Td>
                    <Tag
                      size={'sm'}
                      key={'sm'}
                      variant="subtle"
                      colorScheme={getColorScheme(certificate.trixo.financial_transfers_permitted)}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate.trixo.financial_transfers_permitted === 'yes' ? 'YES' : 'NO'}
                      </TagLabel>
                    </Tag>
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>

          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5,
                  verticalAlign: 'baseline',
                  whiteSpace: 'break-spaces',
                  lineHeight: 1.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <TableCaption placement="top" color="#000" textAlign="start" paddingInlineStart={0}>
                CDD & Travel Rule Policies
              </TableCaption>
              <Tbody>
                <Tr>
                  <Td>
                    Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and
                    Sanctions standards per the requirements of the jurisdiction(s) regulatory
                    regimes where it is licensed/approved/registered?
                  </Td>
                  <Td>
                    <Tag
                      size={'sm'}
                      key={'sm'}
                      variant="subtle"
                      colorScheme={getColorScheme(
                        certificate.trixo.has_required_regulatory_program
                      )}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate.trixo?.has_required_regulatory_program?.toUpperCase() || 'N/A'}
                      </TagLabel>
                    </Tag>
                  </Td>
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
                      colorScheme={getColorScheme(certificate.trixo.financial_transfers_permitted)}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate.trixo?.financial_transfers_permitted.toUpperCase()}
                      </TagLabel>
                    </Tag>
                  </Td>
                </Tr>
                <Tr>
                  <Td>
                    <Trans id="At what threshold and currency does your organization conduct KYC?">
                      At what threshold and currency does your organization conduct KYC?
                    </Trans>
                  </Td>
                  <Td>
                    {currencyFormatter(certificate?.trixo?.kyc_threshold || 0, {
                      currency: certificate.trixo.kyc_threshold_currency
                    })}
                    {'  '}
                    {certificate.trixo.kyc_threshold_currency}
                  </Td>
                </Tr>
                <Tr>
                  <Td>Country of Registration</Td>
                  <Td>USA</Td>
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
                      colorScheme={getColorScheme(certificate?.trixo.must_comply_travel_rule)}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate?.trixo?.must_comply_travel_rule === 'yes' ? 'YES' : 'NO'}
                      </TagLabel>
                    </Tag>
                  </Td>
                </Tr>
                <Tr>
                  <Td>Applicable Regulations</Td>
                  <Td>
                    {certificate.trixo?.applicable_regulations?.map((reg: any) => {
                      if (reg?.name.length > 0) {
                        return <React.Fragment>{reg.name || 'N/A'}</React.Fragment>;
                      }
                    })}
                  </Td>
                </Tr>

                <Tr>
                  <Td>What is the minimum threshold for Travel Rule compliance?</Td>
                  <Td>
                    {currencyFormatter(certificate?.trixo?.compliance_threshold || 0, {
                      currency: certificate.trixo.compliance_threshold_currency
                    })}
                    {'  '}
                    {certificate.trixo.compliance_threshold_currency}
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>

          <TableContainer width="100%">
            <Table
              variant="simple"
              sx={{
                td: {
                  borderBottom: 'none',
                  paddingInlineStart: 0,
                  paddingY: 2.5,
                  verticalAlign: 'baseline',
                  whiteSpace: 'break-spaces',
                  lineHeight: 1.5
                },
                'td:first-child': {
                  width: '50%'
                },
                'td:last-child': {
                  fontWeight: 'bold'
                }
              }}>
              <TableCaption placement="top" color="#000" textAlign="start" paddingInlineStart={0}>
                Data Protection Policies
              </TableCaption>
              <Tbody>
                <Tr>
                  <Td>Is your organization required by law to safeguard PII?</Td>
                  <Td>
                    <Tag
                      size={'sm'}
                      key={'sm'}
                      variant="subtle"
                      colorScheme={getColorScheme(certificate.trixo.must_safeguard_pii)}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate.trixo.must_safeguard_pii ? 'YES' : 'NO'}
                      </TagLabel>
                    </Tag>
                  </Td>
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
                      colorScheme={getColorScheme(certificate.trixo.safeguards_pii)}>
                      <TagLabel fontWeight={'bold'}>
                        {certificate.trixo.safeguards_pii ? 'YES' : 'NO'}
                      </TagLabel>
                    </Tag>
                  </Td>
                </Tr>
              </Tbody>
            </Table>
          </TableContainer>
        </SimpleLayout>
      </Stack>
    </Stack>
  );
}
