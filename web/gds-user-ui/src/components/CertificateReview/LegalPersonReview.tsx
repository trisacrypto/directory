import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button, Tag } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
interface LegalReviewProps {}

const isValidIvmsAddress = (address: any) => {
  if (address) {
    return !!(address.country && address.address_type);
  }
  return false;
};

const hasAddressLine = (address: any) => {
  if (isValidIvmsAddress(address)) {
    return Array.isArray(address.address_line) && address.address_line.length > 0;
  }
  return false;
};

const hasAddressField = (address: any) => {
  if (isValidIvmsAddress(address) && !hasAddressLine(address)) {
    return !!(address.street_name && (address.building_number || address.building_name));
  }
  return false;
};

const hasAddressFieldAndLine = (address: any) => {
  if (hasAddressField(address) && hasAddressLine(address)) {
    console.warn('[ERROR]', 'cannot render address');
    return true;
  }
  return false;
};

export const renderLines = (address: any) => (
  <address data-testid="addressLine">
    {address.address_line.map(
      (addressLine: any, index: number) => addressLine && <div key={index}>{addressLine} </div>
    )}
    <div>{address?.country}</div>
  </address>
);

export const renderField = (address: any) => (
  <address data-testid="addressField">
    {address.sub_department ? (
      <>
        {address?.sub_department} <br />
      </>
    ) : null}
    {address.department ? (
      <>
        {address?.department} <br />
      </>
    ) : null}
    {address.building_number} {address?.street_name}
    <br />
    {address.post_box ? (
      <>
        P.O. Box: {address?.post_box} <br />
      </>
    ) : null}
    {address.floor || address.room || address.building_name ? (
      <>
        {address?.floor} {address?.room} {address?.building_name} <br />
      </>
    ) : null}
    {address.district_name ? (
      <>
        {address?.district_name} <br />
      </>
    ) : null}
    {address.town_name || address.town_location_name || address.country_sub_division ? (
      <>
        {address?.town_name} {address?.town_location_name} {address?.country_sub_division}{' '}
        {address?.post_code} <br />
      </>
    ) : null}
    {address?.country}
  </address>
);

const renderAddress = (address: any) => {
  if (hasAddressFieldAndLine(address)) {
    console.warn('[ERROR]', 'invalid address with both fields and lines');
    return <div>Invalid Address</div>;
  }

  if (hasAddressLine(address)) {
    return renderLines(address);
  }

  if (hasAddressField(address)) {
    return renderField(address);
  }

  console.warn('[ERROR]', 'could not render address');
  return <div>Unparseable Address</div>;
};

const LegalPersonReview: React.FC<LegalReviewProps> = (props) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [legalPerson, setLegalPerson] = React.useState<any>({});

  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.entity
    };
    console.log('legal step data', stepData);
    setLegalPerson(stepData);
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
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Section 2: Legal Person</Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'34px'}
            onClick={() => jumpToStep(2)}
            _hover={{
              bg: '#10aaed'
            }}>
            {' '}
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold' },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody>
              <Tr>
                <Td>Name Identifiers</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td fontStyle={'italic'}>
                  The name and type of name by which the legal person is known.
                </Td>
                <Td>
                  <Tr>
                    {legalPerson.name?.name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <>
                            <Td>{nameIdentifier.legal_person_name}</Td>
                            <Td>
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
                            </Td>
                          </>
                        );
                      }
                    )}
                  </Tr>
                  <Tr>
                    {legalPerson.name?.local_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <>
                            <Td>{nameIdentifier.legal_person_name}</Td>
                            <Td>
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
                            </Td>
                          </>
                        );
                      }
                    )}
                  </Tr>
                  <Tr>
                    {legalPerson.name?.phonetic_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <>
                            <Td>{nameIdentifier.legal_person_name}</Td>
                            <Td>
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
                            </Td>
                          </>
                        );
                      }
                    )}
                  </Tr>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Addresses</Td>
                <Td>
                  <Tr>
                    <Td>
                      {legalPerson?.geographic_addresses?.map((address: any, index: number) => (
                        <React.Fragment key={index}>{renderAddress(address)}</React.Fragment>
                      ))}
                      {/* {legalPerson?.geographic_addresses?.[0]?.address_line.map(
                        (line: any, i: any) => {
                          return <Text key={i}>{line}</Text>;
                        }
                      )} */}
                    </Td>
                    <Td>{legalPerson?.geographic_addresses?.[0] && 'Legal Person'}</Td>
                  </Tr>
                </Td>
                <Td></Td>
              </Tr>
              {/* <Tr>
                <Td>Customer Number</Td>
                <Td></Td>
                <Td></Td>
              </Tr> */}
              <Tr>
                <Td>Country of Registration</Td>
                <Td>{legalPerson?.country_of_registration || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>National Identification</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Number</Td>
                <Td>{legalPerson?.national_identification?.national_identifier}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Type</Td>
                <Td>
                  <Tag color={'white'} bg={'blue.400'} size={'lg'}>
                    {getNationalIdentificationLabel(
                      legalPerson?.national_identification?.national_identifier_type
                    )}
                  </Tag>
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Issue</Td>
                <Td>{legalPerson?.national_identification?.country_of_issue || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Reg Authority</Td>
                <Td>{legalPerson?.national_identification?.registration_authority || 'N/A'}</Td>
                <Td></Td>
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
