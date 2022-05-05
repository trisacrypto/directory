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
  Divider
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { getNameIdentiferTypeLabel } from 'constants/name-identifiers';
import { getNationalIdentificationLabel } from 'constants/national-identification';
import { COUNTRIES } from 'constants/countries';
interface LegalReviewProps {}
// NOTE: need some clean up.
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
          <Heading fontSize={20}>Section 2: Legal Person</Heading>
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
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'bold', paddingLeft: 0 },
              Tr: { borderStyle: 'hidden' }
            }}>
            <Tbody
              sx={{
                '*': {
                  fontSize: '1rem'
                }
              }}>
              <Tr>
                <Td fontSize={'1rem'} fontWeight="bold" colSpan={3}>
                  <Text mb={1}>Name Identifiers</Text>
                  <Divider />
                </Td>
              </Tr>
              <Tr>
                <Td fontStyle={'italic'} pt={0}>
                  The name and type of name by which the legal person is known.
                </Td>
                <Td>
                  <Tr>
                    {legalPerson.name?.name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <React.Fragment key={index}>
                            <Td paddingLeft={0} border="none">
                              {nameIdentifier.legal_person_name || 'N/A'}
                            </Td>
                            <Td paddingLeft={0} border="none">
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
                            </Td>
                          </React.Fragment>
                        );
                      }
                    )}
                  </Tr>
                  <>
                    {legalPerson.name?.local_name_identifiers?.map(
                      (nameIdentifier: any, index: number) => {
                        return (
                          <React.Fragment key={index}>
                            <Td paddingLeft={0} pt={0} border="none">
                              {nameIdentifier.legal_person_name}
                            </Td>
                            <Td paddingLeft={0} pt={0} border="none">
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
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
                              {getNameIdentiferTypeLabel(
                                nameIdentifier.legal_person_name_identifier_type
                              )}
                            </Td>
                          </React.Fragment>
                        );
                      }
                    )}
                  </>
                </Td>
              </Tr>
              <Tr>
                <Td pt={0}>Addresses</Td>
                <Td>
                  <Tr>
                    <Td paddingLeft={0} pt={0}>
                      {legalPerson?.geographic_addresses?.map((address: any, index: number) => (
                        <React.Fragment key={index}>{renderAddress(address)}</React.Fragment>
                      ))}
                      {/* {legalPerson?.geographic_addresses?.[0]?.address_line.map(
                        (line: any, i: any) => {
                          return <Text key={i}>{line}</Text>;
                        }
                      )} */}
                    </Td>
                    <Td pt={0}>{legalPerson?.geographic_addresses?.[0] && 'Legal Person'}</Td>
                  </Tr>
                </Td>
              </Tr>
              <Tr>
                <Td pt={0}>Country of Registration</Td>
                <Td paddingLeft={0} pt={0}>
                  {(COUNTRIES as any)[legalPerson?.country_of_registration] || 'N/A'}
                </Td>
              </Tr>
              <Tr>
                <Td fontSize={'1rem'} fontWeight="bold" colSpan={3}>
                  <Text mb={1}>National Identification</Text>
                  <Divider />
                </Td>
              </Tr>
              <Tr>
                <Td pt={0}>Identification Number</Td>
                <Td paddingLeft={0}>{legalPerson?.national_identification?.national_identifier}</Td>
              </Tr>
              <Tr>
                <Td pt={0}>Identification Type</Td>
                <Td pt={0}>
                  <Tag color={'white'} bg={'blue.400'} size={'lg'}>
                    {getNationalIdentificationLabel(
                      legalPerson?.national_identification?.national_identifier_type
                    )}
                  </Tag>
                </Td>
              </Tr>
              <Tr>
                <Td>Country of Issue</Td>
                <Td>
                  {(COUNTRIES as any)[legalPerson?.national_identification?.country_of_issue] ||
                    'N/A'}
                </Td>
              </Tr>
              <Tr>
                <Td pt={0}>Reg Authority</Td>
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
