import React, { FC, useEffect } from 'react';
import { Stack, Box, Text, Heading, Table, Tbody, Tr, Td, Button } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData, hasValue } from 'utils/utils';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
interface ContactsProps {
  data: any;
}

const ContactsReview = (props: ContactsProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [contacts, setContacts] = React.useState<any>({});
  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      ...getStepperData.contacts
    };
    setContacts(stepData);
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
            Section 3: Contacts
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            onClick={() => jumpToStep(3)}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            Edit{' '}
          </Button>
        </Box>
        <Stack fontSize={'1rem'}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              Tr: { borderStyle: 'hidden' },
              'td:first-child': {
                width: '50%'
              },
              td: {
                borderBottom: 'none',
                paddingInlineStart: 0,
                paddingY: 2.5
              }
            }}>
            <Tbody>
              {['technical', 'legal', 'administrative', 'billing'].map((contact, index) => (
                <Tr key={index}>
                  <Td textTransform="capitalize">
                    {contact === 'legal' ? `Compliance / ${contact}` : contact} Contact
                  </Td>
                  <Td>
                    {hasValue(contacts?.[contact]) ? (
                      <>
                        {contacts?.[contact]?.name && <Text>{contacts?.[contact]?.name}</Text>}
                        {contacts?.[contact]?.email && <Text>{contacts?.[contact]?.email}</Text>}
                        {contacts?.[contact]?.phone && <Text>{contacts?.[contact]?.phone}</Text>}
                      </>
                    ) : (
                      'N/A'
                    )}
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
ContactsReview.defaultProps = {
  data: {}
};
export default ContactsReview;
