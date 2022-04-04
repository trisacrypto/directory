import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';

interface LegalSectionProps {}

const LegalPersonReview: React.FC<LegalSectionProps> = (props) => {
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [legalPerson, setLegalPerson] = React.useState<any>({});

  useEffect(() => {
    const stepData = getStepData(steps, 2);
    if (stepData) {
      setLegalPerson(stepData);
    }
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
                <Td>{legalPerson['']}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Addressess</Td>
                <Td>
                  {legalPerson['entity.geographic_addresses']?.[0]?.address_line[0]} <br />
                  {legalPerson['entity.geographic_addresses']?.[0]?.address_line[1]} <br />
                  {legalPerson['entity.geographic_addresses']?.[0]?.address_line[2]}
                </Td>
                <Td>Legal Person</Td>
              </Tr>
              <Tr>
                <Td>Customer Number</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Registration</Td>
                <Td>{legalPerson['entity.national_identification.country_of_issue']}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>National Identification</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Number</Td>
                <Td>{legalPerson['entity.national_identification.national_identifier']}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Type</Td>
                <Td>{legalPerson['entity.national_identification.national_identifier_type']}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Issue</Td>
                <Td>{legalPerson['entity.national_identification.country_of_issue']}</Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Reg Authority</Td>
                <Td>{legalPerson['entity.national_identification.registration_authority']}</Td>
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
