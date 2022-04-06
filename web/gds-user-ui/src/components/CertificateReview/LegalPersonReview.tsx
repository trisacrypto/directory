import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
interface LegalReviewProps {}

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
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Addressess</Td>
                <Td>
                  <Tr>
                    <Td>
                      {legalPerson?.geographic_addresses?.[0]?.address_line.map(
                        (line: any, i: any) => {
                          return <Text key={i}>{line}</Text>;
                        }
                      )}
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
                <Td>{legalPerson?.national_identification?.country_of_issue}</Td>
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
                <Td>{legalPerson?.national_identification?.national_identifier_type}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Issue</Td>
                <Td>{legalPerson?.national_identification?.country_of_issue}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Reg Authority</Td>
                <Td>{legalPerson?.national_identification?.registration_authority}</Td>
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
