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
  Tag,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import { COUNTRIES } from 'constants/countries';
import { renderAddress } from 'utils/address-utils';
import { addressType } from 'constants/address';
import { Trans } from '@lingui/react';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

interface LegalReviewProps {
  data?: any;
}
// NOTE: need some clean up.

const LegalPersonReview: React.FC<LegalReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [legalPerson, setLegalPerson] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  useEffect(() => {
    const fetchData = async () => {
      const getStepperData = await getRegistrationDefaultValue();
      const stepData = {
        ...getStepperData.entity
      };
      setLegalPerson(stepData);
    };
    fetchData();
  }, [steps]);
  return (
    <Box
      border="2px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={textColor}
      fontSize={18}
      bg={useColorModeValue('white', '#171923')}
      p={5}
      px={5}>
      <Stack>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            <Trans id="Section 2: Legal Person">Section 2: Legal Person</Trans>
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(2)}
            _hover={{
              bg: '#10aaed'
            }}>
            <Trans id="Edit">Edit</Trans>
          </Button>
        </Box>
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
                <Td
                  fontSize={'1rem'}
                  fontWeight="bold"
                  colSpan={3}
                  background={useColorModeValue('#E5EDF1', 'gray.900')}>
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
                      {legalPerson.name?.name_identifiers?.map(
                        (nameIdentifier: any, index: number) => {
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
                        }
                      )}
                    </Tr>
                  </Tbody>
                  <>
                    {legalPerson.name?.local_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
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
                      }
                    )}
                  </>
                  <>
                    {legalPerson.name?.phonetic_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
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
                      }
                    )}
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
                      {legalPerson?.geographic_addresses?.map((address: any, index: number) => (
                        <React.Fragment key={index}>
                          {legalPerson?.geographic_addresses?.length > 1 && (
                            <Text py={1} fontWeight={'bold'}>
                              Address {index + 1} : {(addressType as any)[address.address_type]}
                            </Text>
                          )}
                          {renderAddress(address)}
                        </React.Fragment>
                      ))}
                    </Td>
                    {legalPerson?.geographic_addresses?.length === 1 && (
                      <Td pt={0}>
                        ({(addressType as any)[legalPerson?.geographic_addresses?.[0].address_type]}
                        )
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
                        {(COUNTRIES as any)[legalPerson?.country_of_registration] || 'N/A'}
                      </Td>
                    </Tr>
                  </Tbody>
                </Td>
              </Tr>
              <Tr>
                <Td></Td>
              </Tr>
              <Tr>
                <Td
                  fontSize={'1rem'}
                  pt={'2rem'}
                  fontWeight="bold"
                  colSpan={3}
                  background={useColorModeValue('#E5EDF1', 'gray.900')}>
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
                <Td paddingLeft={0}>{legalPerson?.national_identification?.national_identifier}</Td>
              </Tr>
              <Tr>
                <Td pt={0}>
                  <Trans id="Identification Type">Identification Type</Trans>
                </Td>
                <Td pt={0}>
                  <Tag bg={'blue.400'} size={'lg'}>
                    {getNationalIdentificationLabel(
                      legalPerson?.national_identification?.national_identifier_type
                    )}
                  </Tag>
                </Td>
              </Tr>
              <Tr>
                <Td>
                  <Trans id="Country of Registration">Country of Registration</Trans>
                </Td>
                <Td>{(COUNTRIES as any)[legalPerson?.country_of_registration] || 'N/A'}</Td>
              </Tr>
              <Tr>
                <Td pt={0}>
                  <Trans id="Reg Authority">Reg Authority</Trans>
                </Td>
                <Td pt={0}>
                  {legalPerson?.national_identification?.registration_authority || 'N/A'}
                </Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
LegalPersonReview.defaultProps = {
  data: {}
};
export default LegalPersonReview;
